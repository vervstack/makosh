package main

import (
	//_transport_imports

	"github.com/sirupsen/logrus"

	"github.com/godverv/makosh/cmd/service/makosh"
)

func main() {
	app, err := makosh.New()
	if err != nil {
		logrus.Fatal(err)
	}

	err = app.Start()
	if err != nil {
		logrus.Fatal(err)
	}
}
