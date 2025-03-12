# Calculator-2.0

Этот проект представляет собой простой API калькулятора, написанный на Go. API позволяет выполнять базовые арифметические операции над математическими выражениями.
    
## Как это всё работает?

1. Пользователь вводит выражение
2. Выражение попадает в Оркестратор, а он дает его CalculateHandler
3. CalculateHandler вызывает функцию Calc и дает ей введенное выражение
4. CalculateHandler получает обратно мапу, где ключ - id подвыражения, а значение - само подвыражение
5. CalculateHandler отправляет пользователю id последнего подвыражения, формирует задачи для Агента и запускает его
6. Агент запускается сразу в нескольких горутинах, каждая горутина берет задание и решает его, если в операнде есть id, то заменить его значением из мапы(подробнее ниже)
7. Агент отправляет результаты в Оркестратор в ResultHandler
8. ResulHandler уведомляет пользователя, что выражение решено
Пользователь может в процессе вычисления (или после него) запрашивать как весь список подвыражений, так и одно подвыражение по его id(как запускать проект рассказано ниже
)
## Подробнее о работе API

Здесь будет рассказано о некоторых моментах работы API
### Разбиение на подвыражения функцией Calc
Функция Calc запрашивается CalculateHandler, проводит главные проверки на ошибки в записи выражения и самое главное разбивает все выражение на подвыражения.
![Работа Calc](/images/Calc.jpg)
Выражение разбивается на подзадачи. Подзадаче "операнд1 операция операнд2" присваивается id("id"+порядковый номер подзадачи), на который и заменяется подзадача в выражении. Подзадача и ее id добавляются в мапу mapid, в CalculateHandler возвращается эта самая мапа.
### Паралельное решение подзадач Агентом
Так как у нас есть ограничение на количество запущенных горутин, то нам надо в случае не хватки их изначального количества, запустить новые, но опять не больше установленного значения.
![Agent func main()](/images/Agent_func_main().jpg)
Но как решать подзадачи где какой-либо из операндов с "id"? Для этого результат каждой подзадачи помимо отправки в Оркестратор ResultHandler сохраняется в мапе valmap.
![Agent func Agent()](/images/Agent_func_Agent().jpg)
В valmap по ключу, в виде id всех подзадач, присваивается значение "no". Если один из операндов(или оба) содержат "id", то идет проверка, не заменились ли значения по ключу(id подзадачи) на число, если значение до сих пор "no", то ждем немного времени и проверяем заново. Если значение не "no", а число, то меняем операнд-ы на новое значение. Дальше, когда операнды только числа, вычисляем результат подзадачи и меняем в valmap[id подзадачи] значение на результат подзадачи. Такие действия выполняются со всеми горутинами
## Установка и запуск

### Установка
1. Перейдите в директорию, в которую хотите уставновить проект:
откройте терминал и прописывайте `cd ..`, пока не окажетесь в директории диска
нажмите правой кнопкой мыши по папке, в которую хотите установить проект, скопируйте строчку расположение
введите `cd`+скопированное расположение+\+имя папки(пробел только после cd)
2. Введите `git clone https://github.com/Reit437/Calculator-2.0.git`
3. Введите `go get github.com/Reit437/Calculator-2.0`
4. Введите `go mod tidy`
Готово! Проект установлен
### Запуск
Запуск через терминал
1. Если вы закрыи терминал, то откройте его и повторите 1 пункт из Установки, если не не закрывали, то введите `cd Calculator-2.0`
2. Введите в терминал `go run ./cmd/app/main.go`
3. Откройте GitBash и введите(поле expression можно менять как угодно):
```
curl --location 'http://localhost:5000/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1.2 + ( -8 * 9 / 7 + 56 - 7 ) * 8 - 35 + 74 / 41 - 8"
}'
```
Вы получите обратно id, это id по которому вы сможете получить ответ на всё выражение
ВАЖНО! Пробелы в вашем выражении должны быть строго как в образце сверху иначе будет ошибка.

4. Пока будет решаться выражение вы можете ввести(в GitBash):

1.`curl --location 'localhost:5000/api/v1/expressions'`, чтобы посмотреть все сформированные подзадачи

2.`curl --location 'localhost:5000/api/v1/expressions/id1'`, чтобы посмотреть определенную подзадачу(можете менять id1 на любой id, но строго в таком формате)

6. После надписи в терминале "Выражение решено", можете ввести команду 4.2 с id, который вам дали при вводе выражения и увидеть ответ
7. Для запуска тестов:

   1. Для тестов Calc введите: `go test -v ./pkg/calc`
   2. Для тестов CalculateHandler и ExpressionsHandler введите: `go test -v ./internal/app`
8. Для изменения переменных среды откройте файл в `/internal/config/variables.env` и измените их
## Примеры работы
#### Правильная работа программы
1. Обычное выражение:
   ```
    curl --location 'http://localhost:5000/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
    "expression": "2 + 2 * 2"
    }'
   ```
   Ответ по `curl --location 'localhost:5000/api/v1/expressions/id2'`:
   ```{
   {
    "expression": {
        "Id": "id2",
        "status": "solved",
        "result": "6.000"
        }
   }
   ```
2. Выражение со скобками:
   ```
    curl --location 'http://localhost:5000/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
    "expression": "2 + 2 * 2 / ( 8 * 5 )"
    }'
   ```
   Ответ по `curl --location 'localhost:5000/api/v1/expressions/id4'`:
   ```
    {
    "expression": {
        "Id": "id4",
        "status": "solved",
        "result": "2.100"
        }
    }
   ```
3. Сложное выражение со скобками отрицательными числами и дробными числами:
    ```
    curl --location 'http://localhost:5000/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
    "expression": "1.2 + ( -8 * 9 / 7 + 56 - 7 ) * 8 - 35 + 74 / 41 - 8"
    }'
   ```
    Ответ по `curl --location 'localhost:5000/api/v1/expressions/id10'`:
   ```
   {
    "expression": {
        "Id": "id10",
        "status": "solved",
        "result": "269.710"
        }
    }

   ```
4. Ошибочное выражение:
    ```
    curl --location 'http://localhost:5000/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
    "expression": "2 + 2 * * 2 / ( 8 * 5 )"
    }'
    ```
    Ответ:
   `Невалидные данные`
