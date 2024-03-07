package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"shop/internal/config"
	"shop/internal/domains"
	"shop/internal/repository/postgres"
)

func main() {
	cfg, err := config.New("configs/config.yml")
	exitOnError(err)

	repository, err := postgres.New(cfg.DataBase)
	exitOnError(err)
	err = repository.AddCategory(domains.ProductCategory{Name: "flowers"})
	exitOnError(err)
	err = repository.AddCategory(domains.ProductCategory{Name: "phone"})
	exitOnError(err)
	err = repository.AddCategory(domains.ProductCategory{Name: "aboba"})
	exitOnError(err)
	c, err := repository.GetAllCategories()
	exitOnError(err)
	fmt.Printf("c: %v\n", c)
	c1, err := repository.GetCategoryByID(5)
	exitOnError(err)
	fmt.Println(c1)
	err = repository.AddProduct(domains.Product{Name: "aboba", Description: "a", Price: 156.45, Quantity: 256, CategoryID: 1, ImagePath: "aboba.txt"})
	exitOnError(err)

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

func exitOnError(err error) {
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
}
