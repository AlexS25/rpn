package application

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/AlexS25/rpn/pkg/rpn"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

/*
Функция запуска приложения
читаем stdin после нажатия Enter вывод результата на экран
`exit` - выход из приложения
*/
func (a *Application) Run() error {
	for {
		log.Print("==> Input expression: ")
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

		result, err := rpn.CalcNew(text)
		if err != nil {
			log.Println(text, "calculate failed with error:", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}
