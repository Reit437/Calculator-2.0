package orkestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux" // Импортируйте пакет Gorilla Mux для маршрутизации
	"github.com/joho/godotenv"
)

type ExpressionRequest struct {
	Expression string `json:"expression"` // фиксируем это место
}

type SubExp struct {
	Id     string `json:"id"`     // идентификатор
	Status string `json:"status"` // статус
	Result string `json:"result"` // результат
}

type Response struct {
	ID string `json:"id"` // идентификатор
}

type Task struct {
	Id             string `json:"id"`
	Arg1           string `json:"Arg1"`
	Arg2           string `json:"Arg2"`
	Operation      string `json:"Operation"`
	Operation_time string `json:"Operation_time"`
}

type AllExpressionsResponse struct {
	Expressions []SubExp `json:"expressions"` // массив выражений
}

type ExpressionResponse struct {
	Expression SubExp `json:"expression"` // конкретное выражение
}

type TaskResponse struct {
	Tasks Task `json:"tasks"` // Сообщение
}

type ResultResp struct {
	Id     string `json:"Id"`
	Result string `json:"result"`
}

var (
	subExpressions = make(map[string]string) // новая карта для результатов вычислений
	mu             sync.Mutex
	id             []SubExp
	maxid          int
	tasks          []Task
	res            float64
	v              int
)

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExpressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Вызов функции Calc для разбора выражения
	subExpr := Calc(req.Expression)
	id = []SubExp{}
	maxid = 0
	for expid, exp := range subExpr {
		maxid++
		resp := SubExp{Id: expid, Status: "not solved", Result: exp}
		id = append(id, resp) // добавляем новый результат
	}

	resp := Response{ID: strconv.Itoa(maxid)}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	if err := godotenv.Load("variables.env"); err != nil {
		http.Error(w, "Ошибка при загрузке переменных среды", http.StatusInternalServerError)
	}
	var (
		addTime  = os.Getenv("TIME_ADDITION_MS")
		subTime  = os.Getenv("TIME_SUBTRACTION_MS")
		multTime = os.Getenv("TIME_MULTIPLICATIONS_MS")
		divTime  = os.Getenv("TIME_DIVISIONS_MS")
	)
	for _, i := range id {
		result := i.Result
		add := strings.Index(result, "+")
		sub := strings.Index(result, " - ")
		mult := strings.Index(result, "*")
		div := strings.Index(result, "/")
		var time, ind = "", 0
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
		task := Task{
			Id:             i.Id,
			Arg1:           result[:ind-1],
			Arg2:           result[ind+2:],
			Operation:      string(result[ind]),
			Operation_time: time,
		}

		tasks = append(tasks, task)
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Id < tasks[j].Id
		})
	}
	tasks = append(tasks, Task{
		Id:             "last",
		Arg1:           "0",
		Arg2:           "0",
		Operation:      "no",
		Operation_time: "",
	})
	cmd := exec.Command("go", "run", ".\agent.go")
	err := cmd.Run()
	if err != nil {
		http.Error(w, "Что-то пошло не так", 500)
	}
}

func expressionsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	response := AllExpressionsResponse{Expressions: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Новый обработчик для получения конкретного выражения по ID
func expressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)
	expressionID := vars["id"] // Получаем ID из маршрута

	// Поиск выражения
	for _, exp := range id {
		if exp.Id == expressionID {
			response := ExpressionResponse{Expression: exp}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "    ")
			if err := encoder.Encode(response); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Expression not found", http.StatusNotFound)
}

// Новый обработчик для /internal/task
func taskHandler(w http.ResponseWriter, r *http.Request) {

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	response := TaskResponse{Tasks: tasks[0]}
	tasks = tasks[1:]

	// Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем ответ в JSON и отправляем
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	return
}
func resultHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}
	var result ResultResp
	// Декодируем JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}
	if result.Id[len(result.Id)-1] == byte(maxid+1) {
		http.Error(w, "Нет такой задачи", http.StatusNotFound)
	}
	d, err := strconv.ParseFloat(result.Result, 64)
	if err != nil {
		http.Error(w, "Невалидные данные", http.StatusPaymentRequired)
	}
	for i := 0; i < len(id); i++ {
		if id[i].Id == result.Id {
			id[i].Status = "solved"
			id[i].Result = result.Result
			break
		}
	}
	res = res + d
	v++
	if v == maxid {
		fmt.Println("Выражение решено")
	}
}

/*curl --location 'http://localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1 + 2 + (1 - 5) / 5 * ( 5 / 8)"
}'*/
//curl --location 'localhost/api/v1/expressions'
//curl --location 'localhost/api/v1/expressions/:id'
//curl --location 'localhost/api/v1/cleaner'
