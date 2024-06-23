package web

import (
	"log"
	"net/http"
)

type Server struct {
	PORT     string
	Handlers map[string]http.HandlerFunc
}

func NewServer(port string) *Server {
	return &Server{
		PORT: port,
	}
}

func (s *Server) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers = make(map[string]http.HandlerFunc)
	s.Handlers[path] = handler
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	for path, handler := range s.Handlers {
		mux.HandleFunc(path, handler)
	}

	log.Println("Server running on port", s.PORT)
	err := http.ListenAndServe(":"+s.PORT, mux)
	if err != nil {
		panic(err)
	}
}
