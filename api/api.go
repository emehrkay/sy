package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/emehrkay/sy/service"
)

func New(port string, monitorService service.Monitor, router *http.ServeMux) *server {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return &server{
		port:           port,
		monitorService: monitorService,
		router:         router,
	}
}

type JsonError struct {
	Message string `json:"message"`
}

type server struct {
	monitorService service.Monitor
	router         *http.ServeMux
	port           string
}

func (s *server) Run() error {
	fmt.Printf("SERVER STARTED: %s\n", s.port)
	s.Routes()

	return http.ListenAndServe(s.port, s.router)
}

func (s *server) Routes() {
	s.router.HandleFunc("POST /api/v1/devices/{deviceID}/heartbeat", s.createHeartbeat)
	s.router.HandleFunc("POST /api/v1/devices/{deviceID}/stats", s.createStats)
	s.router.HandleFunc("GET /api/v1/devices/{deviceID}/stats", s.getStats)
}

func (s *server) respondJson(w http.ResponseWriter, response any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (s *server) respondError(w http.ResponseWriter, message string, statusCode int) {
	//TODO: convert errors to http errors -- not found to 404 etc
	s.respondJson(w, JsonError{
		Message: message,
	}, statusCode)
}

func requestBody[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()
	req := new(T)
	err := json.NewDecoder(r.Body).Decode(req)
	return req, err
}
