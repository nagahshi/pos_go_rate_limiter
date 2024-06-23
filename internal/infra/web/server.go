package web

import (
	"log"
	"net/http"
)

type Server struct {
	PORT            string
	Handlers        map[string]http.HandlerFunc
	MiddlewareChain []func(http.HandlerFunc) http.HandlerFunc
}

func NewServer(port string) *Server {
	return &Server{
		PORT: port,
	}
}

func (s *Server) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers = make(map[string]http.HandlerFunc)
	s.Handlers[path] = buildChain(handler, s.MiddlewareChain...)
}

func (s *Server) AddMiddleware(middleware func(http.HandlerFunc) http.HandlerFunc) {
	s.MiddlewareChain = append(s.MiddlewareChain, middleware)
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

func buildChain(f http.HandlerFunc, m ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	if len(m) == 0 {
		return f
	}

	return m[0](buildChain(f, m[1:cap(m)]...))
}
