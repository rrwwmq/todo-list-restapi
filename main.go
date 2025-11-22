package main

import (
	"firstRestAPI/database"
	"firstRestAPI/repository"
	"firstRestAPI/restHTTP"
	"firstRestAPI/todo"
	"fmt"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	repo := repository.NewPostgresRepository()
	todoList := todo.NewList(repo)
	httpHandlers := restHTTP.NewHttpHandlers(todoList)
	httpServer := restHTTP.NewHttpServer(httpHandlers)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("Error starting http server", err)
	}
}
