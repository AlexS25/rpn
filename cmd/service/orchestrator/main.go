package main

import (
	"github.com/AlexS25/rpn/internal/application"
	// "../../../internal/application"
)

func main() {
	app := application.New()
	app.RunOrchestrator()
}
