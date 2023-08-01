package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)


type apiFunc func(http.ResponseWriter, *http.Request) error

func permissionDenied(w http.ResponseWriter){
	writeJSON(w, http.StatusForbidden, ApiError{Error: "Permission denied"})
	
}


//jsonwebtoken shanenigans
func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){

		tokenString := r.Header.Get("x-jwt-token")
		token, err := ValidateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		userID, err := GetID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetAccountByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		if err != nil {
			writeJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func ValidateJWT(tokenString string)(*jwt.Token, error){
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _,ok := t.Method.(*jwt.SigningMethodHMAC); ! ok{
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte("jsonSecret"), nil
	})
}

func createJWT(account *Account) (string, error){
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("jsonSecret"))
}

type ApiError struct{
	Error string 	`json:"error"`
}


type APIServer struct {
	listenAdrr string
	store Storage
}


func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w,r); err != nil{
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewApiServer(listenAdrr string, store Storage) *APIServer {
	return &APIServer{
		listenAdrr: listenAdrr,
		store: store,
	}
}


//handle functions
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error{

	if r.Method == "GET" {
		return s.handleGetAccount(w,r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w,r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w,r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}


//get all accounts
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error{
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, accounts)
}

//get all account by id
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error{
	if r.Method == "GET" {
		id, err:= GetID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)

	if err != nil{
		return err
	}

	return writeJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE"{
		return s.handleDeleteAccount(w,r)
	} 

	return fmt.Errorf("invalid method %s", r.Method)
}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjozNTc5NSwiZXhwaXJlc0F0IjoxNTAwMH0.rFWAXnTum-oWAcUDO1N0cDo-W7UUrgezZ7GLBz3IrrQ

//create account
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error{
	createAcc := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(createAcc);err != nil{
		return err
	}

	account := NewAccount(createAcc.FirstName, createAcc.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Println("json token: ", tokenString)

	return writeJSON(w, http.StatusOK, account)
}

//delete account
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error{
	id, err:= GetID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id) ; err != nil {
		return err
	}

	return writeJSON(w , http.StatusOK, map[string]int{"deleted": id})
}

//transer funds
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error{
	transferReq := new(TranferRequest)

	if err := json.NewDecoder(r.Body).Decode(transferReq) ; err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, transferReq)
}

//write json function
func writeJSON(w http.ResponseWriter, status int, v any) error{
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

//run function
func (s *APIServer) Run(){
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth( makeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("Server is running on port", s.listenAdrr)

	http.ListenAndServe(s.listenAdrr, router)
}

//getId function

func GetID (r *http.Request) (int, error){

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil{
		return id, fmt.Errorf("invalid id %d given", id)
	}

	return id, nil
}