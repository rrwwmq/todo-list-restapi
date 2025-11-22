package restHTTP

import (
	"encoding/json"
	"errors"
	"firstRestAPI/todo"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type httpHandlers struct {
	todoList *todo.List
}

func NewHttpHandlers(todoList *todo.List) *httpHandlers {
	return &httpHandlers{
		todoList: todoList,
	}
}

/*
	pattern: /tasks
	method: post
	info: json in restHTTP request body

	succeed:
		-- status code: 201, created
		-- response body: JSON represent created task
	failed:
		-- status code: 400, 404, 409, 500, ...
		-- response body: JSON with error + time
*/

func (h *httpHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)

		}

		return
	}

	b, err := json.MarshalIndent(todoTask, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
	pattern: /tasks/{title}
	method: get
	info: pattern

	succeed:
		-- status code: 200, ok
		-- response body; JSON represented found task
	failed:
		-- status code: 400, 404, 409, 500, ...
		-- response body: JSON with error + time
*/

func (h *httpHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	task, err := h.todoList.GetTask(title)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(task, "", "   ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response", err)
		return
	}
}

/*
	pattern: /tasks
	method: get
	info: -

	succeed:
		-- status code: 200, ok
		-- response body; JSON represented found tasks
	failed:
		-- status code: 400, 409, 500, ...
		-- response body: JSON with error + time
*/

func (h *httpHandlers) HandlerGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todoList.ListTasks()
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}
	b, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response", err)
		return
	}
}

/*
	pattern: /tasks?completed=false
	method: get
	info: query params

	succeed:
		-- status code: 200, ok
		-- response body; JSON represented found task
	failed:
		-- status code: 400, 404, 409, 500, ...
		-- response body: JSON with error + time
*/

func (h *httpHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks, err := h.todoList.ListUncompletedTasks()
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}
	b, err := json.MarshalIndent(uncompletedTasks, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response", err)
		return
	}
}

/*
	pattern: /tasks/{title}
	method: patch
	info: pattern + json in request body

	succeed:
		-- status code: 200, ok
		-- response body; JSON represented changed task
	failed:
		-- status code: 400, 404, 409, 500, ...
		-- response body: JSON with error + time
*/

func (h *httpHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeTaskDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeTaskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]

	var (
		changedTask todo.Task
		err         error
	)

	if completeTaskDTO.Complete {
		changedTask, err = h.todoList.CompleteTask(title)
	} else {
		changedTask, err = h.todoList.UncompleteTask(title)
	}

	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(changedTask, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response", err)
		return
	}
}

/*
	pattern: /tasks/{title}
	method: delete
	info: pattern

	succeed:
		-- status code: 204, no content
		-- response body; -
	failed:
		-- status code: 400, 404, 409, 500, ...
		-- response body: JSON with error + time
*/

func (h *httpHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	if err := h.todoList.DeleteTask(title); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
