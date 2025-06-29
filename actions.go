package main

import (
	"slices"
	"time"
)

func Add(tasks []Task, t *Task) []Task {
	tasks = append(tasks, *t)
	return tasks
}

func Update(t *Task, desc string) {
	t.Description = desc
	t.UpdatedAt = time.Now()
}

func Get(id int, tasks []Task) (*Task, bool) {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], true
		}
	}
	return nil, false
}

func GetDone(tasks []Task) []Task {
	var done []Task
	for _, tsk := range tasks {
		if tsk.Status == DONE {
			done = append(done, tsk)
		}
	}
	return done
}

func GetInProgress(tasks []Task) []Task {
	var inProgress []Task
	for _, tsk := range tasks {
		if tsk.Status == IN_PROGRESS {
			inProgress = append(inProgress, tsk)
		}
	}
	return inProgress
}

func GetTodo(tasks []Task) []Task {
	var todo []Task
	for _, tsk := range tasks {
		if tsk.Status == TODO {
			todo = append(todo, tsk)
		}
	}
	return todo
}

func Del(id int, tasks *[]Task) {
	*tasks = slices.DeleteFunc(*tasks, func(tsk Task) bool {
		return tsk.ID == id
	})
}

func MarkInProgress(tsk *Task) {
	tsk.Status = IN_PROGRESS
	tsk.UpdatedAt = time.Now()
}

func MarkDone(tsk *Task) {
	tsk.Status = DONE
	tsk.UpdatedAt = time.Now()
}
