package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	ork "github.com/Reit437/Calculator-2.0/internal/app"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/calculate", ork.CalculateHandler).Methods("POST")
	router.HandleFunc("/api/v1/expressions", ork.ExpressionsHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", ork.ExpressionByIDHandler).Methods("GET") // добавляем новый маршрут
	router.HandleFunc("/internal/task", ork.TaskHandler).Methods("GET")                     // добавляем новый маршрут для задачи
	router.HandleFunc("/internal/task", ork.ResultHandler).Methods("POST")

	fmt.Println("Сервер запущен на порту 80...")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
