package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


type apiFunc func(http.ResponseWriter, *http.Request) error

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
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))
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