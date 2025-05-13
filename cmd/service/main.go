package main

import (
	//_transport_imports

	"github.com/sirupsen/logrus"

	"go.vervstack.ru/makosh/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		logrus.Fatal(err)
	}

	err = app.Start()
	if err != nil {
		logrus.Fatal(err)
	}
}
