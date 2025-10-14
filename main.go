package main

import (
	"firstRestAPI/restHTTP"
	"firstRestAPI/todo"
	"fmt"
)

func main() {
	todoList := todo.NewList()
	httpHandlers := restHTTP.NewHttpHandlers(todoList)
	httpServer := restHTTP.NewHttpServer(httpHandlers)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("Error starting http server", err)
	}
}
