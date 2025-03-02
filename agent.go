package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Структура для задачи
type Task struct {
	Id             string `json:"id"`
	Arg1           string `json:"Arg1"`
	Arg2           string `json:"Arg2"`
	Operation      string `json:"Operation"`
	Operation_time string `json:"Operation_time"`
}
type SolvExp struct {
	Id     string `json:"id"`
	Result string `json:"result"`
}

// Структура для всего ответа, который содержит объект "tasks"
type APIResponse struct {
	Tasks Task `json:"tasks"`
}

var (
	mu         sync.Mutex
	result     float64
	ID         string
	valmap     = make(map[string]string)
	stopch     = make(chan bool)
	dig        int
	comp_power int
	n          int
)

// Функция Agent будет выполняться в отдельной горутине
func Agent(wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик в WaitGroup
	var (
		result float64
	)
	// URL вашего внутреннего API
	url := "http://localhost/internal/task"

	// Выполняем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при запросе:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела ответа:", err)
		return
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}

	fmt.Printf("%+v\n", apiResp.Tasks)
	task := apiResp.Tasks
	dig++
	fmt.Println(dig, "dig", comp_power)
	if dig == comp_power && task.Id != "last" {
		n++
	}
	if task.Id == "last" {
		close(stopch)
	}
	mu.Lock()
	ID = task.Id
	valmap[ID] = "no"
	mu.Unlock()
	fmt.Println(valmap)
	fmt.Println(task.Arg1, valmap[task.Arg1], task.Arg2, valmap[task.Arg2])
	for strings.Contains(task.Arg1, "id") || strings.Contains(task.Arg2, "id") {
		fmt.Println(strings.Index(task.Arg1, "id"), task.Arg1)
		if strings.Contains(task.Arg1, "id") {
			if valmap[task.Arg1] != "no" {
				task.Arg1 = strings.Replace(task.Arg1, task.Arg1, valmap[task.Arg1], 1)
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
		fmt.Println(strings.Contains(task.Arg2, "id"), valmap[task.Arg2], ">", task.Arg2, "<")
		if strings.Contains(task.Arg2, "id") {
			task.Arg2 = task.Arg2[:len(task.Arg2)-1]
			fmt.Println(">", task.Arg2, "<")
			if valmap[task.Arg2] != "no" {
				task.Arg2 = strings.Replace(task.Arg2, task.Arg2, valmap[task.Arg2], 1)
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}

	t, _ := strconv.Atoi(task.Operation_time)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(t))
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Таймаут вышел")
			return
		case <-stopch:
			return
		default:
			fmt.Println(task.Operation, task.Id)
			task.Arg2 = task.Arg2[:len(task.Arg2)-1]
			fmt.Println(">", task.Arg1, "<", ">", task.Arg2, "<")
			a, erra := strconv.ParseFloat(task.Arg1, 64)
			b, errb := strconv.ParseFloat(task.Arg2, 64)
			if erra != nil || errb != nil {
				fmt.Println("Ошибка при преобразовании значений.")
				fmt.Println(a, erra, b, errb)
				close(stopch)
				return
			} else {
				switch task.Operation {
				case "+":
					result = a + b
				case "-":
					result = a - b
				case "*":
					result = a * b
				case "/":
					if b == 0 {
						fmt.Println("Ошибка: Деление на нуль")
						return
					} else {
						result = a / b
					}
				}
			}
			valmap[ID] = strconv.FormatFloat(result, 'f', 3, 64)
			res := SolvExp{Id: task.Id, Result: strconv.FormatFloat(result, 'f', 3, 64)}
			fmt.Println(res, valmap)
			body, err := json.Marshal(res)
			if err != nil {
				fmt.Println("Ошибка при сериализации:", err)
				return
			}

			resp, err = http.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println("Ошибка при отправке запроса:", err)
				return
			}
			defer resp.Body.Close()
			return
		}
	}
	return
}

func main() {
	if err := godotenv.Load("variables.env"); err != nil {
		fmt.Println("Ошибка при загрузке переменных среды")
	}
	comp_power, _ = strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	fmt.Println(comp_power)
	n = 1
	var wg sync.WaitGroup
	fmt.Println("Запускаем Agent в горутине...") // Логируем запуск
	for i := 0; i < n; i++ {
		for u := 0; u < comp_power; u++ {
			wg.Add(1)
			go Agent(&wg)
			time.Sleep(1 * time.Second)
		}
	}
	wg.Wait()
	valmap = make(map[string]string)
	fmt.Println("Все горутины завершили работу.")
}
