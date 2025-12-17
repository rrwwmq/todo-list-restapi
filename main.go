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
	sqlDB, err := database.DB.DB()
	if err != nil {
		panic(err)
	}

	defer sqlDB.Close()

	repo := repository.NewGormRepository()
	todoList := todo.NewList(repo)
	httpHandlers := restHTTP.NewHttpHandlers(todoList)
	httpServer := restHTTP.NewHttpServer(httpHandlers)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("Error starting http server", err)
	}
}
