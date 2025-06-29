package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

const FILENAME = "tasks.json"

func main() {
	if err := CreateFileIfNotExists(); err != nil {
		log.Fatalf("create file if not exists: %v", err)
	}

	tasks, err := Load(FILENAME)
	if err != nil {
		log.Fatalf("load file: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments provided")
	}

	switch action := os.Args[1]; action {
	case "add":
		desc := os.Args[2]
		nextId := GetNextId(tasks)
		newTask := Task{
			ID:          nextId,
			Description: desc,
			Status:      TODO,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		tasks = Add(tasks, &newTask)
		if err := WriteJsonToFile(tasks); err != nil {
			log.Fatalf("write to JSON file: %v", err)
		}
		fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)

	case "list":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "	")

		if len(os.Args) == 2 {
			if err := encoder.Encode(tasks); err != nil {
				log.Fatalf("encode tasks: %v", err)
			}
		}
		if len(os.Args) == 3 {
			var tasksToEncode []Task
			switch arg := os.Args[2]; arg {
			case "todo":
				tasksToEncode = GetTodo(tasks)
			case "in-progress":
				tasksToEncode = GetInProgress(tasks)
			case "done":
				tasksToEncode = GetDone(tasks)
			default:
				log.Fatal("Invalid argument")
			}
			if err := encoder.Encode(tasksToEncode); err != nil {
				log.Fatalf("encode tasks: %v", err)
			}
		}

	case "update":
		if len(os.Args) < 4 {
			log.Fatal("Not enough arguments provided")
		}
		_id := os.Args[2]
		desc := os.Args[3]
		id, err := FormatInputId(_id)
		if err != nil {
			log.Fatalf("format id string to int: %v", err)
		}
		tsk, found := Get(id, tasks)
		if found == false {
			fmt.Printf("Task not found with ID: %v\n", id)
			return
		}
		Update(tsk, desc)
		if err := WriteJsonToFile(tasks); err != nil {
			log.Fatalf("write json to file: %v", err)
		}

	case "delete":
		_id := os.Args[2]
		id, err := FormatInputId(_id)
		if err != nil {
			log.Fatalf("format id string to int: %v", err)
		}
		Del(id, &tasks)
		if err := WriteJsonToFile(tasks); err != nil {
			log.Fatalf("write json to file: %v", err)
		}

	case "mark-in-progress":
		_id := os.Args[2]
		id, err := FormatInputId(_id)
		if err != nil {
			log.Fatalf("format id string to int: %v", err)
		}
		tsk, found := Get(id, tasks)
		if found == false {
			fmt.Printf("Task not found with ID: %v\n", id)
			return
		}
		MarkInProgress(tsk)
		if err := WriteJsonToFile(tasks); err != nil {
			log.Fatalf("write json to file: %v", err)
		}

	case "mark-done":
		_id := os.Args[2]
		id, err := FormatInputId(_id)
		if err != nil {
			log.Fatalf("format id string to int: %v", err)
		}
		tsk, found := Get(id, tasks)
		if !found {
			fmt.Printf("Task not found with ID: %v\n", id)
			return
		}
		MarkDone(tsk)
		if err := WriteJsonToFile(tasks); err != nil {
			log.Fatalf("write json to file: %v", err)
		}

	default:
		log.Fatal("Invalid action provided")
	}
}
