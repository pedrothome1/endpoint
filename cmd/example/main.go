package main

import (
	"fmt"
	"github.com/pedrothome1/endpoint"
)

var (
	todosClient = endpoint.New("https://jsonplaceholder.typicode.com/todos")
)

type Todo struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userID"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func (t *Todo) String() string {
	return fmt.Sprintf("(id=%d, userID=%d, title=%q, completed=%v)",
		t.ID, t.UserID, t.Title, t.Completed)
}

func main() {
	var todos []Todo
	_, err := todosClient.Get(endpoint.WithJSONReceivers(&todos, nil))
	if err != nil {
		panic(err)
	}

	for _, todo := range todos {
		fmt.Println(todo.String())
	}
}
