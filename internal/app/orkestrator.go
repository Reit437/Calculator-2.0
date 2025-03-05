package orkestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	Calc "github.com/Reit437/Calculator-2.0/pkg/calc"
	errors "github.com/Reit437/Calculator-2.0/pkg/errors"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type ExpressionRequest struct {
	Expression string `json:"expression"` // Прием первого запроса от пользователя с выражением
}

type SubExp struct { //подвыражения, запрашиваемые пользователем
	Id     string `json:"Id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

type Response struct {
	ID string `json:"Id"` //главный id, по которому пользователь получает конечный ответ на все выражение(ответ для CalculateHandler)
}

type Task struct { //задания, отправляемые агенту
	Id             string `json:"Id"`
	Arg1           string `json:"Arg1"`
	Arg2           string `json:"Arg2"`
	Operation      string `json:"Operation"`
	Operation_time string `json:"Operation_time"`
}

type AllExpressionsResponse struct {
	Expressions []SubExp `json:"expressions"` // ответ для ExpressionsHandler
}

type ExpressionResponse struct {
	Expression SubExp `json:"expression"` // ответ для ExpressionByIdHandler
}

type TaskResponse struct {
	Tasks Task `json:"Tasks"` //ответ для TaskHandler
}

type ResultResp struct { //прием результатов от агента
	Id     string `json:"Id"`
	Result string `json:"result"`
}

var (
	subExpressions = make(map[string]string)
	mu             sync.Mutex
	Id             []SubExp
	Maxid          int
	Tasks          []Task
	res            float64
	v              int
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	/*Прием запроса с выражением от пользователя
	Разбиение его на подвыражения,
	Формирование заданий для агента,
	Запуск агента*/
	if r.Method != http.MethodPost {
		http.Error(w, errors.ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	var req ExpressionRequest //прием запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, errors.ErrUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// вызов функции Calc для разбора выражения
	subExpr, expErr := Calc.Calc(req.Expression)
	//проверка на ошибки при разбиении
	if expErr == 422 {
		http.Error(w, errors.ErrUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}

	Id = []SubExp{}
	Maxid = 0

	//проходимся по мапе из Calc и добавляем в соответствующем формате в Id
	for expid, exp := range subExpr {
		Maxid++
		resp := SubExp{Id: expid, Status: "not solved", Result: exp}
		Id = append(Id, resp)
	}

	//сортировка Id по id
	sort.Slice(Id, func(i, j int) bool {
		id1, _ := strconv.Atoi(Id[i].Id[2:])
		id2, _ := strconv.Atoi(Id[j].Id[2:])
		return id1 < id2
	})

	//формирование ответа
	resp := Response{ID: strconv.Itoa(Maxid)}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	//Формирование заданий
	dir, err := os.Getwd() //установка пути до файла с переменными среды
	if err != nil {
		log.Fatal(err)
	}

	dir = dir[:strings.Index(dir, "Calculator-2.0")+14]
	envPath := filepath.Join(dir, "internal", "config", "variables.env")
	//Загрузка переменных среды
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Ошибка загрузки .env в оркестраторе из %s: %v", envPath, err)
	}

	var (
		addTime  = os.Getenv("TIME_ADDITION_MS")
		subTime  = os.Getenv("TIME_SUBTRACTION_MS")
		multTime = os.Getenv("TIME_MULTIPLICATIONS_MS")
		divTime  = os.Getenv("TIME_DIVISIONS_MS")
	)

	//Формирование массива с заданиями
	for _, i := range Id {
		result := i.Result
		//Ищем знаки операций
		add := strings.Index(result, "+")
		sub := strings.Index(result, " - ")
		mult := strings.Index(result, "*")
		div := strings.Index(result, "/")
		var time, ind = "", 0

		//Если находим операцию, устанавливаем соответствующее время и запоминаем индекс операции
		switch {
		case add != -1:
			time = addTime
			ind = add
		case sub != -1:
			time = subTime
			ind = sub + 1
		case mult != -1:
			time = multTime
			ind = mult
		case div != -1:
			time = divTime
			ind = div
		}

		//Формируем задание
		task := Task{
			Id:             i.Id,
			Arg1:           result[:ind-1],
			Arg2:           result[ind+2:],
			Operation:      string(result[ind]),
			Operation_time: time,
		}
		Tasks = append(Tasks, task) //добавляем задание

		//сортировка заданий по id
		sort.Slice(Tasks, func(i, j int) bool {
			id1, _ := strconv.Atoi(Tasks[i].Id[2:])
			id2, _ := strconv.Atoi(Tasks[j].Id[2:])
			return id1 < id2
		})
	}

	//Создаем последнее задание для остановки агента
	Tasks = append(Tasks, Task{
		Id:             "last",
		Arg1:           "g",
		Arg2:           "g",
		Operation:      "no",
		Operation_time: "",
	})
	//Запускаем агента
	go func() {
		cmd := exec.Command("go", "run", "./internal/services/agent.go")
		err := cmd.Run()
		if err != nil {
			fmt.Println(errors.ErrInternalServerError, err)
		}
	}()
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	//Отправка массива Id с подвыражениями
	mu.Lock()
	defer mu.Unlock()

	//формирование ответа
	response := AllExpressionsResponse{Expressions: Id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(response); err != nil {
		http.Error(w, errors.ErrInternalServerError, http.StatusInternalServerError)
		return
	}
}

func ExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	//Вывод подвыражения по его id
	mu.Lock()
	defer mu.Unlock()

	//определяем запрашиваемый id
	vars := mux.Vars(r)
	expressionID := vars["id"]
	expressId, err := strconv.Atoi(expressionID)
	//проверяем валидность id
	if expressId > Maxid || expressId < 1 || err != nil {
		http.Error(w, errors.ErrNotFound, http.StatusNotFound)
	}

	// поиск выражения
	for _, exp := range Id {
		if exp.Id == expressionID {
			//формирование ответа
			response := ExpressionResponse{Expression: exp}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "    ")
			if err := encoder.Encode(response); err != nil {
				http.Error(w, errors.ErrInternalServerError, http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Expression not found", http.StatusNotFound)
}

// Новый обработчик для /internal/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	// отправка подвыражений агенту
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	//взятие первого элемента и его удаление
	response := TaskResponse{Tasks: Tasks[0]}
	Tasks = Tasks[1:]

	//формирование ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	return
}
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	// прием результатов от агента
	if r.Method != http.MethodPost {
		http.Error(w, errors.ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	//декодирование ответа
	var result ResultResp
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, errors.ErrUnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}

	//проверка на валидность подвыражений
	if result.Id[len(result.Id)-1] == byte(Maxid+1) {
		http.Error(w, errors.ErrNotFound, http.StatusNotFound)
	}

	//замена статуса и результата в Id
	d, err := strconv.ParseFloat(result.Result, 64)
	if err != nil {
		http.Error(w, errors.ErrUnprocessableEntity, http.StatusUnprocessableEntity)
	}
	for i := 0; i < len(Id); i++ {
		if Id[i].Id == result.Id {
			Id[i].Status = "solved"
			Id[i].Result = result.Result
			break
		}
	}
	//Подсчет результата
	res = res + d
	v++
	if v == Maxid {
		fmt.Println("Выражение решено")
	}
}

/*curl --location 'http://localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1.2 + ( -8 * 9 / 7 + 56 - 7 ) * 8 - 35 + 74 / 41 - 8"
}'*/
//curl --location 'localhost/api/v1/expressions'
//curl --location 'localhost/api/v1/expressions/:Id'
