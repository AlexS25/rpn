# Go калькулятор
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

Manage calc's list of trusted argumetns

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

## Работа приложения в режиме **Web-сервер**
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

