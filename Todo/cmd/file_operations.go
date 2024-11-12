package cmd

import (
	"encoding/csv"
	"os"
	"sync"
	"time"
)

var fileMutex sync.Mutex

// Task struct to represent a task
type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}

// LoadTasks loads tasks from a CSV file
func LoadTasks(filename string) ([]Task, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for _, record := range records {
		// Parse the record into a Task struct
		// Assume record[0] is ID, record[1] is Description, record[2] is CreatedAt, record[3] is IsComplete
		// Add parsing logic here
	}

	return tasks, nil
}

// SaveTasks saves tasks to a CSV file
func SaveTasks(filename string, tasks []Task) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, task := range tasks {
		// Convert the Task struct to a record
		// Add conversion logic here
	}

	return nil
} 