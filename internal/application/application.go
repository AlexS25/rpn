package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AlexS25/rpn/pkg/calculation"
)

var DisableLogging bool = false

type Config struct {
	PortNum string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.PortNum = os.Getenv("PORT")
	if config.PortNum == "" {
		config.PortNum = "8080"
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
			log.Println("failed to read expression from console")
		}

		text = strings.TrimSpace(text)

		if text == "exit" {
			log.Println("application was successfully closed")
			return nil
		}

		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, "calculate failed with error:", err)
		} else {
			log.Println(text, "=", result)
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
	//http.HandleFunc("/", CalcHandler)
	//log.Fatal(http.ListenAndServe(port, nil))

	logging("==> Running on port "+port, os.Stdout)

	mux := http.NewServeMux()
	calc := http.HandlerFunc(CalcHandler)
	mux.HandleFunc("/api/v1/calculate", panicCalc(loggingCalc(calc)))
	log.Fatal(http.ListenAndServe(port, mux))
}
