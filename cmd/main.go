package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nagahshi/pos_go_rate_limiter/configs"
	server "github.com/nagahshi/pos_go_rate_limiter/internal/infra/web"
)

func main() {
	// LoadConfig is a function that reads the configuration file and returns a Conf struct
	// with the configuration values
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	server := server.NewServer(cfg.PORT)
	server.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
}
