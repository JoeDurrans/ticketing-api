package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"ticketing-api/data"
	"ticketing-api/types"
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

	router.HandleFunc("GET /ping", makeHTTPHandleFunc(s.handlePing))

	router.HandleFunc("POST /account/login", makeHTTPHandleFunc(s.handleLogin))

	router.HandleFunc("POST /account", IsAdmin(makeHTTPHandleFunc(s.handleCreateAccount)))
	router.HandleFunc("GET /account", IsEditor(makeHTTPHandleFunc(s.handleGetAccounts)))
	router.HandleFunc("GET /account/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleGetAccountByID)))
	router.HandleFunc("PUT /account/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleUpdateAccount)))
	router.HandleFunc("DELETE /account/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleDeleteAccount)))

	router.HandleFunc("POST /ticket", IsAuthenticated(makeHTTPHandleFunc(s.handleCreateTicket)))
	router.HandleFunc("GET /ticket", makeHTTPHandleFunc(s.handleGetTickets))
	router.HandleFunc("GET /ticket/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleGetTicketByID)))
	router.HandleFunc("PUT /ticket/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleUpdateTicket)))
	router.HandleFunc("DELETE /ticket/{id}", IsAuthenticated(makeHTTPHandleFunc(s.handleDeleteTicket)))

	router.HandleFunc("GET /ticket/{id}/chat", makeHTTPHandleFunc(s.handleChatGroup))
	router.HandleFunc("GET /ticket/{id}/chat/message", makeHTTPHandleFunc(s.handleGetMessages))

	server := &http.Server{
		Addr:    s.addr,
		Handler: CreateStack(Logging)(router),
	}

	log.Println("Starting server on", s.addr)

	return server.ListenAndServe()
}

func (s *APIServer) handlePing(w http.ResponseWriter, r *http.Request) error {
	encodeResponse(w, http.StatusOK, &APIResponse{Status: 200, Message: "pong"})
	return nil
}

func getID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, &types.BadRequest{}
	}

	return id, nil
}

func decodeRequest(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return &types.BadRequest{}
	}

	return nil
}

func encodeResponse(w http.ResponseWriter, status int, v any) error {
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
			switch err.(type) {
			case *types.Unauthorized:
				encodeResponse(w, http.StatusUnauthorized, &APIResponse{Status: http.StatusUnauthorized, Message: err.Error()})
			case *types.Forbidden:
				encodeResponse(w, http.StatusForbidden, &APIResponse{Status: http.StatusForbidden, Message: err.Error()})
			case *types.NotFound:
				encodeResponse(w, http.StatusNotFound, &APIResponse{Status: http.StatusNotFound, Message: err.Error()})
			default:
				encodeResponse(w, http.StatusInternalServerError, &APIResponse{Status: http.StatusInternalServerError, Message: err.Error()})
			}
		}
	})
}
