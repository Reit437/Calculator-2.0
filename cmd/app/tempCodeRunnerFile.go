router.HandleFunc("/api/v1/calculate", ork.CalculateHandler).Methods("POST")
	router.HandleFunc("/api/v1/expressions/{id}", ork.ExpressionByIDHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions", ork.ExpressionsHandler).Methods("GET") // добавляем новый маршрут
	router.HandleFunc("/internal/task", ork.TaskHandler).Methods("GET")             // добавляем новый маршрут для задачи
	router.HandleFunc("/internal/task", ork.ResultHandler).Methods("POST")