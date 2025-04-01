package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"service/dilivery"
	"service/infrastructure"
	"service/usecases"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	nlpService, err := infrastructure.NewLocalNLPService()
	if err != nil {
		log.Fatal("Ошибка загрузки NLP:", err)
	}

	log.Println("Connecting to DB with:", os.Getenv("POSTGRES_CONN"))

	db, err := infrastructure.NewPostgresDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	queryProcessor := usecases.NewQueryProcessor(nlpService, db)
	handler := delivery.NewQueryHandler(queryProcessor)

	http.HandleFunc("/query", handler.HandleQuery)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
