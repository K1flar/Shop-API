package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"shop/internal/config"
)

func main() {
	cfg, err := config.New("../configs/config.yml")
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	done := make(chan struct{})
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port), nil); err != nil {
			io.WriteString(os.Stderr, err.Error())
			os.Exit(1)
		}
		done <- struct{}{}
	}()
	fmt.Printf("Starting server %s:%s\n", cfg.Server.Host, cfg.Server.Port)
	<-done
}
