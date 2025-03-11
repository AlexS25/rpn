[Структура проекта и как скачать](#sprint_1)

[Работа приложения в режиме **Web-сервер**](#sprint_2)

[Распределенный вычислитель арифметических выражений (с оркестратором и агентами)](#sprint_3)


# <a id="sprint_1">Go калькулятор</a>
Данный проект реализует базовые возможности калькулятора и позволяет вводить выражение одной строкой, а на выходе получать ответ или описание ошибки, если такая закралась в выражение. Реализовано два механизма работы:
 - [x] Консольный ввод выражений;
 - [x] Режим web-сервера.

В качестве механизма работы с выражением была выбрана так-называемая *Польская нотация* которая позволяет проводить операции над распарсенными числами выражения.
Примеры работы приложения проводились под ОС *Ubuntu*


## Структура проекта
```
.
├── cmd  // Для хранения пакетов `main`
│   └── main.go  // Точка входа в программу
├── internal  // Для хранения приватных пакетов
│   └── application
│       ├── application.go  // Само приложение
│       └── application_test.go
├── pkg  // Для хранения публичных пакетов
│   ├── calculation  // Пакеты самого приложения
│   │   ├── calculation.go  // Пакет с логикой самого приложения
│   │   ├── calculation_test.go
│   │   └── errors.go  // Описание базовых ошибок
│   └── collections  // Дополнительные пакеты коллекций, которые использыются в проете
│       └── stack
│           ├── stack_empty_interface.go  // Пакет для работы со стеком через пустой интерфейс 
│           ├── stack_empty_interface_test.go
│           ├── stackstr.go  // Пакет для работы со стеком через `string`
│           └── stackstr_test.go
└── README.md
```

## Загрузка проекта
Для загрузки проекта достаточно выполнить следующую команду в терминале на вашем ПК:

```bash
git clone https://github.com/AlexS25/rpn.git
```

После чего перейти в директорию с проектом

## Просмотр режимов работы приложения
Запустив приложение с ключем `--help` можно увидеть подсказку, c доступными ключами запуска.

```bash
go run ./cmd/main.go --help
```

В выводе команды можно увидеть следующее:

```bash
Usage: cmd [arguments]

Manage calc\'s list of trusted argumetns

  --help     - Show current help.
  --cmd      - Use command line interface.
  --srv      - Run as web-server. (default mod)

Example of application launch:
./cmd --srv
```

Из чего понятно, что для запуска приложения в режиме **консольного ввода** достаточно запустить его с ключем `--cmd`:

```bash
go run ./cmd/main.go --cmd
```

После чего, получим приглашение для ввода выражения, и результат решения, а для выхода достаточно ввести команду `exit`:

```bash
==> Input expression: 2 + 2 * 2
2024/12/22 21:24:01 2 + 2 * 2 = 6
==> Input expression: (2+2) * 2
2024/12/22 21:24:21 (2+2) * 2 = 8
==> Input expression: 2+2*+2
2024/12/22 21:24:36 2+2*+2 calculate failed with error: invalid expression: extra operator in expression: 2+2*+2
==> Input expression: (2+2*2
2024/12/22 21:24:55 (2+2*2 calculate failed with error: invalid expression: not found close bracket in expression: (2+2*2
==> Input expression: 2+2)*2
2024/12/22 21:25:03 2+2)*2 calculate failed with error: invalid expression: not found open bracket in expression: 2+2)*2
==> Input expression: exit
2024/12/22 21:25:08 application was successfully closed
```

Для запуска калькулятора в режиме работы **Web-сервера** требуется запустить приложение с ключем `--srv` или запустить просто без указания ключей:

```bash
go run ./cmd/main.go
```

После запуска, приложение ожидает `HTTP` запрос (по умолчанию на порту *8080*)

## <a id="sprint_2">Работа приложения в режиме **Web-сервер**</a>
### Запуск сервера
Как было сказано выше, для запуска приложения в режиме **Web-сервер**, его достаточно запустить без ключей или с ключем `--srv`. В таком режиме, приложение будет работать на порту *8080*. 

```bash
go run ./cmd/main.go
2024/12/22 21:42:29 [INFO] ==> Running on port :8080
```

Для изменения порта, достаточно указать номер порта в переменной окружения `PORT`. Это можно сделать двумя способами.
Первый вариант, указать переменную перед запускаемым приложением. В этом случае, переменная не будет записана в переменные окружения.

```bash
PORT=8081 go run ./cmd/main.go 
2024/12/22 21:43:15 [INFO] ==> Running on port :8081
```

Второй вариант подразумевает экспортирование в окружение переменных посредством комманды `export`.

```bash
# Экспортируем переменную в окружение переменных
export PORT=8082
# Можно проверить, что переменная действительно попала в окружение переменных, выполнив команду
env | grep PORT
# В результате увидим, что переменная там есть
PORT=8082
# Теперь осталось только запустить приложение
go run ./cmd/main.go
# Видим, что приложение запустилось на порту 8082
2024/12/22 21:53:14 [INFO] ==> Running on port :8082
```

Теперь, пока мы не закроем терминал, приложение можно запускать не указывая номер порта.
Если по каким-то причинам переменную окружения требуется удалить, а терминал не хочется закрывать, можно это сделать командой `unset`.

```bash
unset PORT
```

После чего, переменная будет удалена из окружения переменных.

### Отправка запросов
Отправка запросов проводится по пути `/api/v1/calculate` на заданный порт хоста, по умолчанию это порт 8080.
Для начала запустим **web-сервер** на порту по умолчанию (8080)

```bash
go run ./cmd/main.go 
```

В результате чего увидим сообщение:

```bash
...
2024/12/22 22:08:07 [INFO] ==> Running on port :8080
...
```

После запуска сервера можно приступить к отправке запросов на сервер. Для этого можно использовать консольную утилиту `curl`. 
Запрос отправляется в формате `json` с одним полем `expression` и будет выглядеть следующим образом:

```json
{
  "expression": "2+2*2"
}
```

Ответ также будет получен в `json` формате с одним полем `result` и будет выглядеть следующим образом:

```json
{
  "result": 6
}
```

Примеры отправки запроса утилитой `curl`:
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

На выходе получим ответ:

```bash
{
  "result":6
}
```

Так-же ведется журналировани, где можно увидеть следующую информацию:

```bash
...
2024/12/22 22:19:44 [INFO] Star HTTP request: "/api/v1/calculate"
2024/12/22 22:19:44 [INFO] request: 2+2*2
2024/12/22 22:19:44 [INFO] result: 6.000000
2024/12/22 22:19:44 [INFO] HTTP request: "/api/v1/calculate", method: "POST", duration: 63937
...
```

Можно попробовать отправить пример со скобками и пробелами в выражении:


```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2 + 2) *   2"
}'
```

На выходе получим ответ:

```bash
{
  "result":8
}
```

Журналирование так-же зафиксирует данную информацию:
```bash
2024/12/22 22:23:47 [INFO] Star HTTP request: "/api/v1/calculate"
2024/12/22 22:23:47 [INFO] request: (2 + 2) *   2
2024/12/22 22:23:47 [INFO] result: 8.000000
2024/12/22 22:23:47 [INFO] HTTP request: "/api/v1/calculate", method: "POST", duration: 51676
```

Помимо прочего обрабатываются возможные ошибки. Ответ так-же приходит в `json` формате с двумя полями: `error` - именование ошибки и `description` - подробное описание ошибки, если оно присутствует.
Пример с лишним оператором:

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*+2"
}'
```

В ответ получем:

```bash
{
  "error": "Expression is not valid",
  "description": "invalid expression: extra operator in expression: (2+2)*+2"
}
```

Ошибка так-же фиксируется в журнале:
```bash
2024/12/22 22:30:57 [INFO] Star HTTP request: "/api/v1/calculate"
2024/12/22 22:30:57 [INFO] request: (2+2)*+2
2024/12/22 22:30:57 [ERROR] invalid expression: extra operator in expression: (2+2)*+2
2024/12/22 22:30:57 [INFO] HTTP request: "/api/v1/calculate", method: "POST", duration: 65990
```

Пример с лишними скобками:

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2*+2"
}'
```

В ответ получим:

```bash
{
  "error":"Expression is not valid",
  "description":"invalid expression: not found close bracket in expression: (2+2*+2"
}
```

Фиксация в журнале

```bash
2024/12/22 22:35:02 [INFO] Star HTTP request: "/api/v1/calculate"
2024/12/22 22:35:02 [INFO] request: (2+2*+2
2024/12/22 22:35:02 [ERROR] invalid expression: not found close bracket in expression: (2+2*+2
2024/12/22 22:35:02 [INFO] HTTP request: "/api/v1/calculate", method: "POST", duration: 39717
```

Если будет отправлено некорректное поле в запросе, то оно будет приравнено к пустому полу:

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression1": "(2+2)*+2"
}'
```

В ответ получим:

```json
{
  "error":"Expression is not valid",
  "description":"expression is required"
}
```

Фиксация в журнале:

```bash
2024/12/22 22:38:05 [INFO] Star HTTP request: "/api/v1/calculate"
2024/12/22 22:38:05 [ERROR] expression is required
2024/12/22 22:38:05 [INFO] HTTP request: "/api/v1/calculate", method: "POST", duration: 105756
```

## Проверка кода возврата
Помимо самого выражения, возвращается и код возврата. Для его вывода можно добавить ключ `-v` для утилиты `curl` перенаправить поток ошибок на поток вывода, и с помощью утилиты `grep` отфильтровать вывод.

Пример возврата выражения с кодом `200`

```bash
curl -v --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*2"
}' 2>&1 | grep -E "expression|error|HTTP"
```

В выводе получим следующий ответ:

```bash
> POST /api/v1/calculate HTTP/1.1
< HTTP/1.1 200 OK
{"result":8}
```

Как видно, код возврата 200, и сам `result` тоже присутствует.

Так-же можно проверить и ошибочное выражение

```bash
curl -v --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*+2"
}' 2>&1 | grep -E "result|error|HTTP"
```

В данном случае код возврата будет 422

```bash
> POST /api/v1/calculate HTTP/1.1
< HTTP/1.1 422 Unprocessable Entity
{
  "error":"Expression is not valid",
  "description":"invalid expression: extra operator in expression: (2+2)*+2"
}
```

## <a id="sprint_3">Распределенный вычислитель арифметических выражений (с оркестратором и агентами)</a>

Для загрузки проекта необходимо выполнить команду в терминале:

```bash
git clone https://github.com/AlexS25/rpn.git
```

После клонирования проекта следует перейти в директорию с проектором и проводить работу из нее
```bash
cd rpn
```

### Запуск оркестратора

Запускаем оркестартор:

```bash
go run -race cmd/service/orchestrator/main.go
```

После запуска увидим, что оркестратор запущен

```bash
...
[INFO] ==> Running `Orchestrator` on port :8080
...
```

При попытке запустить вторую копию увидим сообщение, что порт уже используется

```bash
...
[INFO] ==> Running `Orchestrator` on port :8080
listen tcp :8080: bind: address already in use
...
```

Если требуется запустить вторую копию приложения, то следует указать порт, на котором следует запустить приложение:
```bash
# Можно задать переменну окружения перед самой командой
PORT=8081 go run -race cmd/service/orchestrator/main.go

# Или Экспортируем переменную в окружение переменных
export PORT=8082
```

### Добавление вычисления арифметического выражения

Для добавления вырежния можно воспользоваться утилитой `curl`. Подключаться к серверу для добавления задания будем по порту по умолчанию - *8080*, если сервер был запущен на другом порту, то следует его изменить.

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2"
}'

curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2-2"
}'

curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1+2*3/(4-3)"
}'

curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "5 -4 + 3 + 2 *1"
}'

curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "55/11*2"
}'
```

После добавления задания, мы получим `id` выражения, чтобы в дальнейшем можно было его отследить

```bash
...
{"id":1}
{"id":2}
{"id":3}
{"id":4}
{"id":5}
...
```

Если требуется посмотреть код ответа, то можно попробовать добавить следующим образом

```bash
curl -v --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*2"
}' 2>&1 | grep -E "expression|error|HTTP"
```

В ответ получим информацию с кодом возврата **201**:

```bash
> POST /api/v1/calculate HTTP/1.1
< HTTP/1.1 201 Created
```

При попытке отправить некорректное выражение (например лишняя скобка) 

```bash
curl -v --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*2)"
}' 2>&1 | grep -E "expression|error|HTTP"
```

Увидим следующую картину, где код возврата будет уже **422**, и дополнительно получим информацию по ошибке

```bash
> POST /api/v1/calculate HTTP/1.1
< HTTP/1.1 422 Unprocessable Entity
{"error":"Expression is not valid","description":"invalid expression: not found open bracket in expression: (2+2)*2)"}
```

### Получение списка выражений 

Для получения списка выражений, необходимо отправить `GET` запрос на оркестратор следующим образом

```bash
curl --location 'localhost:8080/api/v1/expressions'
```

В результате чего получим список отправленных выражений

```bash
{"expressions":
  [
   {"id":1,"status":"expression available","result":0}
   ,{"id":2,"status":"expression available","result":0}
   ,{"id":3,"status":"expression available","result":0}
   ,{"id":4,"status":"expression available","result":0}
   ,{"id":5,"status":"expression available","result":0}
   ,{"id":6,"status":"expression available","result":0}
  ]
}
```

Как видим, добавлены только успешные выражения

Если требуется получить информацию по конкретному выражения, то следует выполнить следующий запрос, где после `:` идет `id` выражения

```bash
curl --location 'localhost:8080/api/v1/expressions/:1'
```

Видим состояние нашего выражения

```bash
{"id":1,"status":"expression available","result":0}
```

Можно попробовать получить не существующее выражение

```bash
curl -v --location  'localhost:8080/api/v1/expressions/:111' 2>&1 | grep -E "expression|error|HTTP"
```

На выходе получим сообщение об ошибке и код возврата **404**

```bash
...
> GET /api/v1/expressions/:111 HTTP/1.1
< HTTP/1.1 404 Not Found
Not found expression by id "111".
...
```

### Запуск агента (worker)

Для запуска агента необходимо открыть терминал и перейти в директорию проекта.

Далее выполняем команду 

```bash
go run cmd/service/worker/main.go
```

После чего запустится агент, который будет брать задания с оркестратора для выполнения. Агентов можно запускать несколько штук, для этого придется открыть новый терминал и повторно запустить команду выше.

Посредством переменных окружения, агенту можно задать количество потоков обрабоки выражения `COMPUTING_POWER`, а так-же время выполнения для каждого выражения

```bash
# Количество потоков обработки выражений
export COMPUTING_POWER=2
# время выполнения операции сложения в миллисекундах
export TIME_ADDITION_MS=1000
# время выполнения операции вычитания в миллисекундах
export TIME_SUBTRACTION_MS=2000
# время выполнения операции умножения в миллисекундах
export TIME_MULTIPLICATIONS_MS=3000
# время выполнения операции деления в миллисекундах 
export TIME_DIVISIONS_MS=4000

go run cmd/service/worker/main.go
```

После того, как непрерывно будет идти сообшение

```bash
...
[INFO] ==> Run getTask
[INFO] ==> Run getTask
...
```

Можно проверять состояние наших выражений
```bash
curl --location 'localhost:8080/api/v1/expressions'
...
{"expressions":
  [
    {"id":6,"status":"expression calculated","result":8}
    ,{"id":1,"status":"expression calculated","result":4}
    ,{"id":2,"status":"expression calculated","result":0}
    ,{"id":3,"status":"expression calculated","result":7}
    ,{"id":4,"status":"expression calculated","result":6}
    ,{"id":5,"status":"expression calculated","result":10}
  ]
}
...
```

Так-же по логам можно проверить время выполнения операций, которое мы задали через переменные окружения

```bash
...
[INFO] ==> Run sendTask
Arg1: "1", arg2: "6.000000", operation: "+"
[INFO] Current response: {"id":3,"result":7}
[INFO] ==> Stop goroutine num: 7, duration: 1002314316
...
[INFO] ==> Run sendTask
Arg1: "4.000000", arg2: "2", operation: "*"
[INFO] Current response: {"id":6,"result":8}
[INFO] ==> Stop goroutine num: 7, duration: 3002548978
...
```


### Эмулируем работу агента

Для проверки корректности работы воркера, попробуем проделать операцию вычисления вручную
Предварительно не забываем отключить и оркестратор и агент

Запускам оркестратора

```bash
go run -race cmd/service/orchestrator/main.go
```

Пытаемся получить задания для обработки

```bash
curl --location 'localhost:8080/internal/task'
```

Получаем сообщение об ошибке

```bash
{"error":"Not found tasks","description":"Not found tasks"}
```
G
Заоднем провери код ошибки

```bash
curl -v --location 'localhost:8080/internal/task/' 2>&1 | grep -E "expression|error|HTTP"
```

Получим код **404**

```bash
> GET /internal/task/ HTTP/1.1
< HTTP/1.1 404 Not Found
{"error":"Not found tasks","description":"Not found tasks"}
```

Добавляем выражение

```bash
curl --location 'localhost:8080/api/v1/calculate/' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*2"
}'
```

Пытаемся повторно получить задания для обработки

```bash
curl --location 'localhost:8080/internal/task/'
```

Получили ответ

```bash
{"id":1,"arg1":"2","arg2":"2","operation":"+","operation_time":"2025-03-12 02:08:07.371"}
```

Пытаемся повторно получить выражени для рассчета

```bash
curl -v --location 'localhost:8080/internal/task/' 2>&1 | grep -E "expression|error|HTTP"
## Получили ошибку
> GET /internal/task/ HTTP/1.1
< HTTP/1.1 404 Not Found
{"error":"Not found tasks","description":"Not found tasks"}
```

Проверяем состояние заданий

```bash
curl --location 'localhost:8080/api/v1/expressions/'
```

Получаем ответ, что задание вычисляется

```bash
{"expressions":[{"id":1,"status":"expression evaluated","result":0}]}
```

Решаем выражение от отправляем результат

```bash
curl --location 'localhost:8080/internal/task/' \
--header 'Content-Type: application/json' \
--data '{
  "id": 1,
  "result": 4
}'
```

Смотрим состояние

```bash
curl --location 'localhost:8080/api/v1/expressions/'
```

Видим, что наже выраженин опять доступно

```bash
{"expressions":[{"id":1,"status":"expression available","result":0}]}j
```

Опять получаем задание для обработки

```bash
curl --location 'localhost:8080/internal/task/'
```

Видим ответ

```bash
{"id":1,"arg1":"4.000000","arg2":"2","operation":"*","operation_time":"2025-03-12 02:18:09.119"}
```

Решаем и отправляем решение


```bash
curl --location 'localhost:8080/internal/task/' \
--header 'Content-Type: application/json' \
--data '{
  "id": 1,
  "result": 8
}'
```

Смотрим состояние

```bash
curl --location 'localhost:8080/api/v1/expressions/'
```

Видим, что у нас достсупно только одно выражение и оно решено

```bash
{"expressions":[{"id":1,"status":"expression calculated","result":8}]}
```
