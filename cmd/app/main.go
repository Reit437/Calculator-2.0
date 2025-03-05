package main

import (
	"fmt"
	"log"
	"net/http"

	ork "github.com/Reit437/Calculator-2.0/internal/app"
	"github.com/gorilla/mux"
)

func main() {
	// Создаем новый роутер и запросы
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/api/v1/calculate", ork.CalculateHandler).Methods("POST")
	router.HandleFunc("/api/v1/expressions/{id}", ork.ExpressionByIDHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions", ork.ExpressionsHandler).Methods("GET")
	router.HandleFunc("/internal/task", ork.TaskHandler).Methods("GET")
	router.HandleFunc("/internal/task", ork.ResultHandler).Methods("POST")

	// Запускаем сервер
	fmt.Println("Сервер запущен на порту 80...")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
