package main

import (
	"log"
	"os"
)

type Command int

const (
	Add Command = iota
	Update
	Delete
	MarkInProgress
	MarkDone
	List
)

var commandName = map[Command]string{
	Add:            "add",
	Update:         "update",
	Delete:         "delete",
	MarkInProgress: "mark-in-progress",
	MarkDone:       "mark-done",
	List:           "list",
}

func (cn Command) String() string {
	return commandName[cn]
}

func main() {

	command := os.Args[1]
	switch command {
	case Add.String():
		log.Println("Add command received")
	}
}

type Task struct {
	Id int
}
