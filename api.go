package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct{
	Error string
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

	account := NewAccount("Abel", "Wanyonyi")

	return writeJSON(w, http.StatusOK, account)
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
	return nil
}

//transer funds
func (s *APIServer) handleTranser(w http.ResponseWriter, r *http.Request) error{
	return nil
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

	log.Println("Server is running on port", s.listenAdrr)

	http.ListenAndServe(s.listenAdrr, router)
}