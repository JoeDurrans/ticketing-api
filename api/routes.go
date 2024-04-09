package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"ticketing-api/data"
)

type APIServer struct {
	addr       string
	db         *data.DataAdapter
	chatGroups *sync.Map
}

func CreateAPIServer(addr string, db *data.DataAdapter) *APIServer {
	return &APIServer{
		addr:       addr,
		db:         db,
		chatGroups: &sync.Map{},
	}
}

func (s *APIServer) Start() error {
	router := http.NewServeMux()
	log.Println("Starting server on", s.addr)

	//ping
	router.HandleFunc("GET /ping", makeHTTPHandleFunc(s.handlePing))

	router.HandleFunc("POST /account/login", makeHTTPHandleFunc(s.handleLogin))

	router.HandleFunc("POST /account", IsAdmin(makeHTTPHandleFunc(s.handleCreateAccount)))

	router.HandleFunc("GET /account", IsEditor(makeHTTPHandleFunc(s.handleGetAccounts)))

	router.HandleFunc("GET /ticket", makeHTTPHandleFunc(s.handleGetTickets))
	router.HandleFunc("GET /account/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleGetAccountByID)))
	router.HandleFunc("PUT /account/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleUpdateAccount)))
	router.HandleFunc("DELETE /account/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleDeleteAccount)))
	router.HandleFunc("POST /ticket", IsAuthenticated(makeHTTPHandleFunc(s.handleCreateTicket)))
	router.HandleFunc("GET /ticket/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleGetTicketByID)))
	router.HandleFunc("PUT /ticket/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleUpdateTicket)))
	router.HandleFunc("DELETE /ticket/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleDeleteTicket)))
	router.HandleFunc("GET /ticket/{id}/chat", makeHTTPHandleFunc(s.handleChatGroup))

	stack := CreateStack(Logging)

	server := &http.Server{
		Addr:    s.addr,
		Handler: stack(router),
	}

	return server.ListenAndServe()
}

func (s *APIServer) handlePing(w http.ResponseWriter, r *http.Request) error {
	EncodeResponse(w, http.StatusOK, &APIResponse{Status: 200, Message: "pong"})
	return nil
}

func getID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("error converting id: %w", err)
	}

	return id, nil
}

func DecodeRequest(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding request: %w", err)
	}

	return nil
}

func EncodeResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			EncodeResponse(w, http.StatusInternalServerError, &APIResponse{Status: http.StatusInternalServerError, Message: err.Error()})
		}
	})
}
