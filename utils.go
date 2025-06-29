package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
)

func WriteJsonToFile(tasks []Task) error {
	file, _ := os.OpenFile(FILENAME, os.O_WRONLY, 0644)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")
	err := encoder.Encode(tasks)
	if err != nil {
		return err
	}
	return nil
}

func Load(fileName string) ([]Task, error) {
	var tasks []Task
	file, openErr := os.Open(fileName)
	defer file.Close()
	if openErr != nil {
		return nil, openErr
	}
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		if errors.Is(io.EOF, err) {
			return []Task{}, nil
		} else {
			return nil, err
		}
	}
	return tasks, nil
}

func CreateFileIfNotExists() error {
	if _, err := os.Stat(FILENAME); errors.Is(err, os.ErrNotExist) {
		_, createErr := os.Create(FILENAME)
		if createErr != nil {
			return err
		}
	}
	return nil
}

func FormatInputId(id string) (int, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return intId, nil
}

func GetNextId(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	} else {
		return tasks[len(tasks)-1].ID + 1
	}
}
