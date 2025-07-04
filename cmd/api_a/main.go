package main

import (
	"github.com/DeboraIK/lab2-OTEL/internal/webserver"
)

func main() {
	webserver := webserver.NewWebServer("a")
	webserver.Start("a")
}
