package application

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlexS25/rpn/pkg/calculation"
)

var DisableLogging bool = false

// const apiVersion = "v1"

type Config struct {
	PortNum         string
	CompPower       string
	TimeAddition    int
	TimeSubtraction int
	TimeMultiply    int
	TimeDivision    int
}

/*
TIME_ADDITION_MS - время выполнения операции сложения в миллисекундах
TIME_SUBTRACTION_MS - время выполнения операции вычитания в миллисекундах
TIME_MULTIPLICATIONS_MS - время выполнения операции умножения в миллисекундах
TIME_DIVISIONS_MS - время выполнения операции деления в миллисекундах
*/

func ConfigFromEnv() *Config {
	config := new(Config)

	config.PortNum = os.Getenv("PORT")
	if config.PortNum == "" {
		config.PortNum = "8080"
	}

	config.CompPower = os.Getenv("COMPUTING_POWER")
	if config.CompPower == "" {
		config.CompPower = "4"
	}

	val := os.Getenv("TIME_ADDITION_MS")
	config.TimeAddition, _ = strconv.Atoi(val)
	if config.TimeAddition == 0 {
		config.TimeAddition = int(time.Millisecond) * 100
	}

	val = os.Getenv("TIME_SUBTRACTION_MS")
	config.TimeSubtraction, _ = strconv.Atoi(val)
	if config.TimeSubtraction == 0 {
		config.TimeSubtraction = int(time.Millisecond) * 100
	}

	val = os.Getenv("TIME_MULTIPLICATIONS_MS")
	config.TimeMultiply, _ = strconv.Atoi(val)
	if config.TimeMultiply == 0 {
		config.TimeMultiply = int(time.Millisecond) * 100
	}

	val = os.Getenv("TIME_DIVISIONS_MS")
	config.TimeDivision, _ = strconv.Atoi(val)
	if config.TimeDivision == 0 {
		config.TimeDivision = int(time.Millisecond) * 100
	}

	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

/*
Функция запуска приложения
читаем stdin после нажатия Enter вывод результата на экран
`exit` - выход из приложения
*/
func (a *Application) Run() error {
	for {
		fmt.Print("==> Input expression: ")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')

		if err != nil {
			logging("failed to read expression from console", os.Stderr)
		}

		text = strings.TrimSpace(text)

		if text == "exit" {
			logging("application was successfully closed", os.Stdout)
			return nil
		}

		result, err := calculation.Calc(text)
		if err != nil {
			logging(fmt.Sprintf("%q, calculate failed with error: ", text, err), os.Stderr)
		} else {
			logging(fmt.Sprintf("%q = %v", text, result), os.Stdout)
		}
	}
}

func logging(mess string, typeMess io.Writer) {
	if DisableLogging {
		return
	}
	var prefix string = "[ERROR]"
	if typeMess == nil || typeMess == os.Stdout {
		typeMess = os.Stdout
		prefix = "[INFO]"
	}

	log.SetOutput(typeMess)
	log.Println(prefix, mess)
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}

type ResponseError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

// Structures for orchestrator and worker
type ResponseId struct {
	Id int `json:"id"`
}

type Expression struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

type ResponseExprs struct {
	Expressions []Expression `json:"expressions"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	respError := new(ResponseError)
	respError.Error = "Expression is not valid"

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respError.Description = err.Error()
		jsonError, _ := json.Marshal(respError)

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonError)
		logging(err.Error(), os.Stderr)
		return
	}

	if len(request.Expression) == 0 {
		mess := "expression is required"
		respError.Description = mess
		jsonError, _ := json.Marshal(respError)

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonError)
		logging(mess, os.Stderr)
		return
	}
	logging("request: "+request.Expression, os.Stdout)

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		respError.Description = err.Error()
		jsonError, _ := json.Marshal(respError)

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonError)

		logging(err.Error(), os.Stderr)
	} else {
		response := new(Response)
		response.Result = result
		json.NewEncoder(w).Encode(response)

		mess := fmt.Sprintf("result: %f", result)
		//fmt.Fprint(w, mess)
		logging(mess, os.Stdout)
	}
}

func AddExprHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	respError := new(ResponseError)
	respError.Error = "Expression is not valid"

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respError.Description = err.Error()
		jsonError, _ := json.Marshal(respError)

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonError)
		logging(err.Error(), os.Stderr)
		return
	}

	if len(request.Expression) == 0 {
		mess := "expression is required"
		respError.Description = mess
		jsonError, _ := json.Marshal(respError)

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonError)
		logging(mess, os.Stderr)
		return
	}
	logging("request: "+request.Expression, os.Stdout)

	// result, err := calculation.Calc(request.Expression)
	id, err := SafeExpr.AddExpr(request.Expression)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		respError.Description = err.Error()
		jsonError, _ := json.Marshal(respError)

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonError)

		logging(err.Error(), os.Stderr)
	} else {
		response := new(ResponseId)
		response.Id = id
		json.NewEncoder(w).Encode(response)

		mess := fmt.Sprintf("id: %v", id)
		//fmt.Fprint(w, mess)
		logging(mess, os.Stdout)
	}
}

func getExprs() ResponseExprs {
	SafeExpr.mtx.RLock()
	defer SafeExpr.mtx.RUnlock()

	expr := Expression{}
	respExprs := ResponseExprs{}
	for id, de := range SafeExpr.de {
		expr.Id = id

		// 0-free, 1-busy, 2-solved, 5-selected
		var state string
		switch de.state {
		case 0:
			state = "expression available"
		case 1, 5:
			state = "expression evaluated"
		case 2:
			state = "expression calculated"
		default:
			state = "unknown state"
		}
		expr.Status = state
		res, _ := SafeExpr.getSolution(id)
		expr.Result = res

		respExprs.Expressions = append(respExprs.Expressions, expr)
	}
	return respExprs
}

func SendExprStateHandler(w http.ResponseWriter, r *http.Request) {
	response := getExprs()
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	// Так-как muxRouter не обрабатывает корректно неявные запросы - костылим
	path := strings.Trim(r.URL.Path, "/")
	segments := strings.Split(path, "/")

	// запрос без id
	if len(segments) == 3 {
		if err := json.NewEncoder(w).Encode(response); err != nil {
			mess := "Error generation response: " + err.Error()

			http.Error(w, mess, http.StatusInternalServerError)
			logging(mess, os.Stderr)
			return
		}
	}

	// обработка запроса с id
	if len(segments) == 4 {
		idS := segments[3]
		idS = idS[1:]
		// fmt.Fprintf(w, "==> ID is: %q\n", idS)
		logging(fmt.Sprintf("==> ID is: %q\n", idS), os.Stdout)

		id, err := strconv.Atoi(idS)
		if err == nil {
			if expr, ok := getExprById(id); ok {
				// Отправляем ответ, если все ок
				if err := json.NewEncoder(w).Encode(expr); err != nil {
					mess := "Error generation response: " + err.Error()
					// http.Error(w, mess, http.StatusInternalServerError)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(mess))
					logging(mess, os.Stderr)
					return
				}
			} else {
				// обработка, если значение не найдено
				mess := fmt.Sprintf("Not found expression by id %q.", idS)
				// http.Error(w, mess, http.StatusNotFound)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(mess))
				logging(mess, os.Stderr)
				// return
			}
		} else {
			// Не является числом
			mess := fmt.Sprintf("Value %q is not identifier.", idS)
			http.Error(w, mess, http.StatusNotFound)
			logging(mess, os.Stderr)
			return
		}
	}
	// else {
	// 	mess := fmt.Sprintf("Not found %q.", r.RequestURI)
	// 	// http.Error(w, mess, http.StatusNotFound)
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte(mess))
	// 	logging(mess, os.Stderr)
	// }
}

func getExprById(id int) (Expression, bool) {
	SafeExpr.mtx.RLock()
	defer SafeExpr.mtx.RUnlock()

	expr := Expression{}
	if de, ok := SafeExpr.Get(id); ok {
		expr.Id = id

		fmt.Println("==> Cur de:", de)

		// 0-free, 1-busy, 2-solved, 5-selected
		var state string
		switch de.state {
		case 0:
			state = "expression available"
		case 1, 5:
			state = "expression evaluated"
		case 2:
			state = "expression calculated"
		default:
			state = "unknown state"
		}
		expr.Status = state
		res, _ := SafeExpr.getSolution(id)
		expr.Result = res
		return expr, true

	}

	return expr, false
}

// func SendExprByIdStateHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	defer r.Body.Close()

// 	vars := mux.Vars(r)
// 	idS := vars["id"]
// 	fmt.Fprintf(w, "==> ID is: %q\n", idS[1:])
// 	id, err := strconv.Atoi(idS[1:])
// 	if err == nil {
// 		if expr, ok := getExprById(id); ok {
// 			// Отправляем ответ, если все ок
// 			if err := json.NewEncoder(w).Encode(expr); err != nil {
// 				mess := "Error generation response: " + err.Error()
// 				// http.Error(w, mess, http.StatusInternalServerError)
// 				w.WriteHeader(http.StatusInternalServerError)
// 				w.Write([]byte(mess))
// 				logging(mess, os.Stderr)
// 				return
// 			}
// 		} else {
// 			// обработка, если значение не найдено
// 			mess := fmt.Sprintf("Not found expression by id %q.", idS[1:])
// 			// http.Error(w, mess, http.StatusNotFound)
// 			w.WriteHeader(http.StatusNotFound)
// 			w.Write([]byte(mess))
// 			logging(mess, os.Stderr)
// 			// return
// 		}
// 	} else {
// 		// Не является числом
// 		mess := fmt.Sprintf("Value %q is not identifier.", idS[1:])
// 		http.Error(w, mess, http.StatusNotFound)
// 		logging(mess, os.Stderr)
// 		return
// 	}
// }

type Task struct {
	Id            int    `json:"id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime string `json:"operation_time"`
}

type TaskResult struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

func doTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	// Для ошибок в json
	respError := new(ResponseError)

	taskResult := new(TaskResult)
	err := json.NewDecoder(r.Body).Decode(&taskResult)

	// if r.Method == http.MethodGet {
	if err != nil {
		// тут происзодит обработка запроса новой таски
		logging("running `doTask` to send the data", os.Stdout)
		task := new(Task)

		// заготовка для ошибки
		mess := "Not found tasks"
		respError.Error = mess
		respError.Description = mess

		// получаем свободны id
		id := SafeExpr.getFreeExprId()

		if id > 0 {
			fmt.Print("==> Cur value: ")
			fmt.Println(SafeExpr.Get(id))
			cr_expr, err := SafeExpr.getExprForTask(id)

			if err != nil {
				respError.Description = err.Error()
				jsonError, _ := json.Marshal(respError)

				w.WriteHeader(http.StatusNotFound)
				w.Write(jsonError)
				logging(mess, os.Stderr)
				return
			}

			// Все гуд, задачу сформировали
			logging(fmt.Sprintf("Cur expression for id %v :%v\n", id, cr_expr), os.Stdout)
			task.Id = id
			task.Arg1 = cr_expr[0]
			task.Arg2 = cr_expr[1]
			task.Operation = cr_expr[2]
			task.OperationTime = time.Now().Format("2006-01-02 15:04:05.000")

			json.NewEncoder(w).Encode(task)
			return

		} else {
			// задач пока нет
			jsonError, _ := json.Marshal(respError)

			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonError)
			logging(mess, os.Stderr)
			return
		}

		// } else if r.Method == http.MethodPost {
	} else {
		// тут получаем обратно обработанные данные
		logging("running `doTask` to get the data", os.Stdout)

		// taskResult := new(TaskResult)
		// err := json.NewDecoder(r.Body).Decode(&taskResult)

		// // если не получилось десериализовать данные
		// if err != nil {
		// 	respError.Error = "Error geting task execution"
		// 	respError.Description = err.Error()
		// 	jsonError, _ := json.Marshal(respError)

		// 	w.WriteHeader(http.StatusUnprocessableEntity)
		// 	w.Write(jsonError)
		// 	logging(respError.Error+" "+respError.Description, os.Stderr)
		// 	return
		// }

		if err := SafeExpr.pushValFromTask(taskResult.Id, fmt.Sprintf("%f", taskResult.Result)); err != nil {
			// Елси не нашли данные для записи
			respError.Error = "Error writing result"
			respError.Description = err.Error()
			jsonError, _ := json.Marshal(respError)

			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonError)
			return
		}

	}

}

func loggingCalc(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		mess := fmt.Sprintf("Star HTTP request: %q", r.RequestURI)
		logging(mess, os.Stdout)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		mess = fmt.Sprintf("HTTP request: %q, method: %q, duration: %d",
			r.RequestURI, r.Method, duration)
		logging(mess, os.Stdout)
	})
}

func panicCalc(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				mess := fmt.Sprintf("panic calculation error: %v", err)
				respError := new(ResponseError)
				respError.Error = "Internal server error"
				respError.Description = mess
				jsonError, _ := json.Marshal(respError)

				//http.Error(w, string(jsonError), http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonError)

				logging(mess, os.Stderr)

			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *Application) RunServer() {
	port := fmt.Sprintf(":%s", a.config.PortNum)

	logging("==> Running on port "+port, os.Stdout)

	mux := http.NewServeMux()
	calc := http.HandlerFunc(CalcHandler)
	mux.HandleFunc("/api/v1/calculate", panicCalc(loggingCalc(calc)))
	mux.HandleFunc("/api/v1/calculate/", panicCalc(loggingCalc(calc)))
	log.Fatal(http.ListenAndServe(port, mux))
}

var SafeExpr safeExpression

func (a *Application) RunOrchestrator() {
	port := fmt.Sprintf(":%s", a.config.PortNum)

	logging("==> Running `Orchestrator` on port "+port, os.Stdout)

	SafeExpr = *NewSafeExpression()

	router := http.NewServeMux()
	// router := mux.NewRouter()

	addExpr := http.HandlerFunc(AddExprHandler)
	// Добавление арифметичкого выражения
	router.HandleFunc("/api/v1/calculate/", panicCalc(loggingCalc(addExpr)))
	router.HandleFunc("/api/v1/calculate", panicCalc(loggingCalc(addExpr)))
	// получение списка выражений
	sendExprState := http.HandlerFunc(SendExprStateHandler)
	router.HandleFunc("/api/v1/expressions/", panicCalc(loggingCalc(sendExprState)))

	// Получение выражения по `id`
	// sendExprByIdState := http.HandlerFunc(SendExprByIdStateHandler)
	// router.HandleFunc("/api/v1/expressions/{id}", panicCalc(loggingCalc(sendExprByIdState)))

	// Получение задачи для выполнения
	// Прием результата обработки данных
	doTask := http.HandlerFunc(doTaskHandler)
	router.HandleFunc("/internal/task/", panicCalc(loggingCalc(doTask)))

	log.Fatal(http.ListenAndServe(port, router))

	expr := "2 + 3*4/(3 - 1)"

	safeExpr := NewSafeExpression()
	safeExpr.AddExpr(expr)

	safeExpr.AddExpr("2 + 2")
	safeExpr.AddExpr("1 + 2 * 3")

	var id int
	// var cr_expr []string = make([]string, 0)
	for i := 0; i < 18; i++ {
		id = safeExpr.getFreeExprId()
		fmt.Print("==> Cur value: ")
		fmt.Println(safeExpr.Get(id))
		cr_expr, err := safeExpr.getExprForTask(id)

		if err == nil {
			fmt.Printf("==> Cur expression for id %v :%v\n", id, cr_expr)

			if len(cr_expr) == 0 {
				if res, ok := safeExpr.getSolution(id); ok {
					fmt.Printf("==> FOR EXPRESSION WITH ID %d, RESULT = %v\n", id, res)
				}
				// fmt.Println("==> FINE")
				// return
				continue
			}

			res, err := calculation.EvalSimpleExpr(cr_expr[0], cr_expr[1], cr_expr[2])
			if err == nil {
				fmt.Printf("For id %v expr %q, res = %q\n", id, strings.Join(cr_expr, " "), res)

				if err := safeExpr.pushValFromTask(id, res); err != nil {
					fmt.Printf("==> New values for ID %v: ", id)
					fmt.Println(safeExpr.Get(id))

				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("\n")
	for val, expr := range safeExpr.de {
		fmt.Println(val, expr)
	}

	fmt.Print("Result: ")
	fmt.Println(safeExpr.getSolution(id))
}

func (task *Task) getTaskToRun(url string) error {
	// Создаем POST запрос для отправки результата
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating GET request.\n%v", err)
	}
	//Задаем заголовок, что хотим получить json тело
	req.Header.Set("Accept", "application/json")
	// Ну тут выполняем его
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error running GET request.\n%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		return fmt.Errorf("Request failed: %s\n Response: %s\n %v", resp.Status, body, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed read data from response.\n %v", err)
	}
	// jsonTask := Task{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		return fmt.Errorf("Failed json unmarshal.\n %v", err)
	}
	logging(fmt.Sprintf("Current json:\n%v", task), os.Stdout)
	return nil
}

func (task *Task) sendTaskResult(url string) error {

	res, err := calculation.EvalSimpleExpr(task.Arg1, task.Arg2, task.Operation)
	if err != nil {
		return err
	}
	val, _ := strconv.ParseFloat(res, 64)

	taskResult := new(TaskResult)
	taskResult.Id = task.Id
	taskResult.Result = val
	logging(fmt.Sprintf("Current response: %v", taskResult), os.Stdout)

	reqTaskRes, _ := json.Marshal(taskResult)
	// Создаем POST запрос для отправки ответа
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqTaskRes))
	if err != nil {
		return fmt.Errorf("Error creating POST request.\n%v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Ну тут выполняем его
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error running GET request.\n%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		return fmt.Errorf("Request failed: %s\n Response: %s\n %v", resp.Status, body, err)
	}
	return nil

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return fmt.Errorf("Failed read data from response.\n %v", err)
	// }
	// jsonTask := new(Task)
	// err = json.Unmarshal(body, &jsonTask)
	// if err != nil {
	// 	return fmt.Errorf("Failed json unmarshal.\n %v", err)
	// }
	// logging(fmt.Sprintf("Current json:\n%s", jsonTask), os.Stdout)
	// return nil
}

func (a *Application) RunWorker() {
	port, err := strconv.Atoi(a.config.PortNum)
	if err != nil {
		logging(fmt.Sprintf("Error getting port number.\n%s", err.Error()), os.Stderr)
		return
	}
	portS := strconv.Itoa(port)
	// compPower := fmt.Sprintf("%s", a.config.CompPower)

	logging("==> Running `Worker` on port number: "+portS, os.Stdout)

	url := fmt.Sprintf("http://localhost:%v/internal/task", port)

	func(string) {
		task := Task{}
		err := task.getTaskToRun(url)
		if err != nil {
			logging(err.Error(), os.Stderr)
			return
		}

		err = task.sendTaskResult(url)
		if err != nil {
			logging(err.Error(), os.Stderr)
			return
		}
	}(url)

}
