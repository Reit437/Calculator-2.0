package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/calculate", calculateHandler).Methods("POST")
	router.HandleFunc("/api/v1/expressions", expressionsHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", expressionByIDHandler).Methods("GET") // добавляем новый маршрут
	router.HandleFunc("/internal/task", taskHandler).Methods("GET")                     // добавляем новый маршрут для задачи
	router.HandleFunc("/internal/task", resultHandler).Methods("POST")

	fmt.Println("Сервер запущен на порту 80...")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
