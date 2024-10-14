package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Command int
type TaskStatus int

const (
	Add Command = iota
	Update
	Delete
	MarkInProgress
	MarkDone
	List
)

const (
	Todo TaskStatus = iota
	InProgress
	Done
)

var commandName = map[Command]string{
	Add:            "add",
	Update:         "update",
	Delete:         "delete",
	MarkInProgress: "mark-in-progress",
	MarkDone:       "mark-done",
	List:           "list",
}

var taskStatusName = map[TaskStatus]string{
	Todo:       "todo",
	InProgress: "in-progress",
	Done:       "done",
}

func (cn Command) String() string {
	return commandName[cn]
}
func (ts TaskStatus) String() string {
	return taskStatusName[ts]
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getCWD() string {
	wd, err := os.Getwd()
	check(err)
	return wd
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please enter a command")
		return
	}
	command := os.Args[1]
	switch command {
	case Add.String():
		desc := os.Args[2]
		log.Println("Add command received")
		taskId := addTask(desc)
		log.Printf("Task added successfully (ID:%d)", taskId)
	case Update.String():
		sid := os.Args[2]
		newDesc := os.Args[3]
		id, err := strconv.Atoi(sid)
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}
		err = updateTask(id, newDesc)
		if err != nil {
			fmt.Println(err)
			return
		}
	case Delete.String():
		sid := os.Args[2]
		id, err := strconv.Atoi(sid)
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}
		err = deleteTask(id)
		if err != nil {
			fmt.Println(err)
			return
		}

	default:
		fmt.Println("Invalid command")
		return
	}
}

type Task struct {
	Id          int       `json:"-"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func addTask(desc string) int {
	id := getTaskId()
	task := Task{
		Id:          id,
		Description: desc,
		Status:      Todo.String(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	task.save()
	return id
}

func (t Task) save() {
	ok, jsonPath := isDataExist()
	if !ok {
		data := make(map[int]Task)
		data[t.Id] = t
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Marshaling data failed %v", err)
		}
		if err = os.WriteFile(jsonPath, bytes, 0644); err != nil {
			log.Fatalf("Saving data failed %v", err)
		}
		return
	}
	tasks, err := loadData(jsonPath)
	if err != nil {
		log.Fatalf("Loading data failed %v", err)
	}
	tasks[t.Id] = t
	err = updateData(jsonPath, tasks)
}

func isDataExist() (bool, string) {
	exPath := getCWD()
	jsonPath := filepath.Join(exPath, "data.json")
	if _, err := os.Stat(jsonPath); err != nil {
		return false, jsonPath
	}
	return true, jsonPath
}
func getTaskId() int {
	ok, jsonPath := isDataExist()
	if !ok {
		return 1
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatal("Loading JSON file failed")
	}
	var tasks map[int]Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatal("Parsing JSON data failed", err)
	}
	return len(tasks) + 1
}

func loadData(jsonPath string) (map[int]Task, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	var tasks map[int]Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func updateData(jsonPath string, tasks map[int]Task) error {
	bytes, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("marshaling data failed %v", err)
	}
	if err = os.WriteFile(jsonPath, bytes, 0644); err != nil {
		return fmt.Errorf("saving data failed %v", err)
	}
	return nil
}

func updateTask(id int, desc string) error {
	ok, jsonPath := isDataExist()
	if !ok {
		return errors.New("no tasks created yet")
	}
	tasks, err := loadData(jsonPath)
	if err != nil {
		return err
	}
	task, ok := tasks[id]
	if !ok {
		return errors.New("task not found")
	}

	task.Description = desc
	tasks[id] = task
	err = updateData(jsonPath, tasks)
	if err != nil {
		return fmt.Errorf("updating data failed %v", err)
	}
	return nil
}

func deleteTask(id int) error {
	ok, jsonPath := isDataExist()
	if !ok {
		return errors.New("no tasks created yet")
	}
	tasks, err := loadData(jsonPath)
	if err != nil {
		return err
	}
	_, ok = tasks[id]
	if !ok {
		return errors.New("task not found")
	}
	delete(tasks, id)
	err = updateData(jsonPath, tasks)
	if err != nil {
		return err
	}
	return nil
}
