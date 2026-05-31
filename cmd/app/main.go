package main

import (
	"koito_proxy/internal/app"
	"koito_proxy/internal/logging"
)

func main() {
	logging.Setup()

	a, err := app.New()
	if err != nil {
		panic(err)
	}
	if err := a.Run(); err != nil {
		panic(err)
	}
}
