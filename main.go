package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"time"
)

type task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const (
	TODO        = "todo"
	IN_PROGRESS = "in progress"
	DONE        = "done"
	FILENAME    = "tasks.json"
)

func main() {
	if _, err := os.Stat(FILENAME); errors.Is(err, os.ErrNotExist) {
		_, createErr := os.Create(FILENAME)
		if createErr != nil {
			log.Fatal(err)
		}
	}

	tasks, err := load(FILENAME)
	if err != nil {
		log.Fatal(err)
	}

	var nextInt int
	if len(tasks) == 0 {
		nextInt = 1
	} else {
		nextInt = tasks[len(tasks)-1].ID + 1
	}

	if len(os.Args) < 2 {
		fmt.Println("Invalid argument")
		return
	}

	action := os.Args[1]

	switch action {
	case "add":
		desc := os.Args[2]
		newTask := task{
			ID:          nextInt,
			Description: desc,
			Status:      TODO,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		tasks = add(tasks, newTask)
		writeToFile(tasks)
		fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)

	case "list":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "	")

		if len(os.Args) == 2 {
			if err := encoder.Encode(tasks); err != nil {
				log.Fatal(err)
				return
			}
		}
		if len(os.Args) == 3 {
			arg := os.Args[2]
			switch arg {
			case "todo":
				encoder.Encode(getTodo(tasks))
			case "in-progress":
				encoder.Encode(getInProgress(tasks))
			case "done":
				encoder.Encode(getDone(tasks))
			default:
				fmt.Println("Invalid")
			}
		}

	case "update":
		id := os.Args[2]
		desc := os.Args[3]
		idString, _ := strconv.Atoi(id)
		tsk, found := get(idString, tasks)
		if found == false {
			fmt.Printf("Task not found with ID: %v\n", id)
			return
		}
		update(tsk, desc)
		writeToFile(tasks)

	case "delete":
		id := os.Args[2]
		idString, _ := strconv.Atoi(id)
		del(idString, &tasks)
		writeToFile(tasks)

	case "mark-in-progress":
		id := os.Args[2]
		idString, _ := strconv.Atoi(id)
		tsk, found := get(idString, tasks)
		if found == false {
			fmt.Printf("Task not found with ID: %v\n", id)
			return
		}
		markInProgress(tsk)
		writeToFile(tasks)

	case "mark-done":
		id := os.Args[2]
		idString, _ := strconv.Atoi(id)
		tsk, found := get(idString, tasks)
		if found == false {
			fmt.Printf("Task not found with ID: %v\n", id)
			return
		}
		markDone(tsk)
		writeToFile(tasks)

	default:
		fmt.Println("invalid arg")
	}
}

func load(fileName string) ([]task, error) {
	var tasks []task
	file, openErr := os.Open(fileName)
	defer file.Close()
	if openErr != nil {
		log.Fatal(openErr)
	}
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		if errors.Is(io.EOF, err) {
			return []task{}, nil
		} else {
			return nil, err
		}
	}
	return tasks, nil
}

func add(tasks []task, t task) []task {
	tasks = append(tasks, t)
	return tasks
}

func update(t *task, desc string) {
	t.Description = desc
	t.UpdatedAt = time.Now()
}

func get(id int, tasks []task) (*task, bool) {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], true
		}
	}
	return nil, false
}

func getDone(tasks []task) []task {
	var done []task
	for _, tsk := range tasks {
		if tsk.Status == DONE {
			done = append(done, tsk)
		}
	}
	return done
}

func getInProgress(tasks []task) []task {
	var inProgress []task
	for _, tsk := range tasks {
		if tsk.Status == IN_PROGRESS {
			inProgress = append(inProgress, tsk)
		}
	}
	return inProgress
}

func getTodo(tasks []task) []task {
	var todo []task
	for _, tsk := range tasks {
		if tsk.Status == TODO {
			todo = append(todo, tsk)
		}
	}
	return todo
}

func del(id int, tasks *[]task) {
	*tasks = slices.DeleteFunc(*tasks, func(tsk task) bool {
		return tsk.ID == id
	})
}

func markInProgress(tsk *task) {
	tsk.Status = IN_PROGRESS
	tsk.UpdatedAt = time.Now()
}

func markDone(tsk *task) {
	tsk.Status = DONE
	tsk.UpdatedAt = time.Now()
}

func writeToFile(tasks []task) {
	file, _ := os.OpenFile(FILENAME, os.O_WRONLY, 0644)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")
	encoder.Encode(tasks)
}
