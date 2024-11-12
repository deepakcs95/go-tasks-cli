/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var mutex sync.Mutex
type Task struct {
	ID int
	Task string
	Completed bool
	CreatedAt time.Time
}

func main() {
	fmt.Println("Starting Todo List Application")
	err := initCSV("tasks.csv")
	if err != nil {
		fmt.Println("Error initializing CSV file:", err)
	}
	addTask("tasks.csv", "First Task")
	addTask("tasks.csv", "Second Task")
	addTask("tasks.csv", "Third Task")
	addTask("tasks.csv", "Fourth Task")
	listTasks("tasks.csv")
	updateTask("tasks.csv", "3", "true")
	listTasks("tasks.csv")
	// deleteTask("tasks.csv", "2")
	listTasks("tasks.csv")
}


func initCSV(filename string) error {
	mutex.Lock()
	defer mutex.Unlock()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		} 
		defer file.Close()

	
	writer := csv.NewWriter(file)
		defer writer.Flush()

		record := []string{
			"ID",
			"Task",
			"Completed",
			"CreatedAt",
	}
	if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write task: %w", err)
		}
		fmt.Println("CSV file created successfully")
	}
	return nil
}
func openCSVFile(filename string, flag int) (*os.File, error) {
	file, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

func writeCSVFile(file *os.File, task Task) error {
	mutex.Lock()
	defer mutex.Unlock()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		strconv.Itoa(task.ID),
		task.Task,
		strconv.FormatBool(task.Completed),
		task.CreatedAt.Format(time.RFC3339),
	}
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write task: %w", err)
	}
	return nil
}

func readCSVFile(file *os.File) ([]Task, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	var tasks []Task
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}
		id, _ := strconv.Atoi(record[0])
		completed, _ := strconv.ParseBool(record[2])
		createdAt, _ := time.Parse(time.RFC3339, record[3])
		
		tasks = append(tasks, Task{
			ID:        id,
			Task:      record[1],
			Completed: completed,
			CreatedAt: createdAt,
		})
	}
	return tasks, nil
}

func getNextID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func addTask(filename string, taskDescription string) error {
	// First read existing tasks to get the next ID
	file, err := openCSVFile(filename, os.O_RDONLY)
	if err != nil {
		return err
	}
	records, err := readCSVFile(file)
	file.Close()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	// Get next ID and create new task
	nextID := getNextID(records)
	newTask := Task{
		ID: nextID,
		Task: taskDescription,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// Append the new task
	file, err = openCSVFile(filename, os.O_APPEND|os.O_WRONLY)
	if err != nil {
		return err
	}
	defer file.Close()

	err = writeCSVFile(file, newTask)
	if err != nil {
		return fmt.Errorf("failed to write task: %w", err)
	}
	fmt.Printf("Task added successfully with ID: %d\n", nextID)
	return nil
}

func listTasks(filename string) error {
	file, err := openCSVFile(filename, os.O_RDONLY)
	if err != nil {
		return err
	}
	defer file.Close()

	records, err := readCSVFile(file)
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	for _, record := range records {
		fmt.Println(record)
	}
	return nil
}

func getTaskByID(filename string, taskID string) (Task, error) {
	file, err := openCSVFile(filename, os.O_RDONLY)
	if err != nil {
		return Task{}, err
	}
	defer file.Close()

	records, err := readCSVFile(file)
	if err != nil {
		return Task{}, fmt.Errorf("failed to read tasks: %w", err)
	}

	var task Task
	for _, record := range records {
		if strconv.Itoa(record.ID) == taskID {
			task = record
			break
		}
	}

	if task.ID == 0 {
		return Task{}, fmt.Errorf("task with ID %s not found", taskID)
	}

	return task, nil
}

 // Start of Selection
func deleteTask(filename string, taskID string) error {
	 

	// First, read all records with read-only access
	readFile, err := openCSVFile(filename, os.O_RDONLY)
	if err != nil {
		return err
	}
	records, err := readCSVFile(readFile)
	readFile.Close() // Close immediately after reading
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	// Filter out the task to be deleted
	var remainingRecords []Task
	for _, record := range records {
		if strconv.Itoa(record.ID) != taskID {
			remainingRecords = append(remainingRecords, record)
		}
	}

	// Truncate and write to the file
	writeFile, err := os.Create(filename) // This truncates the file
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer writeFile.Close()

	writer := csv.NewWriter(writeFile)
	defer writer.Flush()

	

	// Write all remaining records
	for _, record := range remainingRecords {
		if err := writeCSVFile(writeFile, record); err != nil {
			return fmt.Errorf("failed to write task: %w", err)
		}
	}

	fmt.Println("Task deleted successfully")
	return nil
}

func updateTask(filename string, taskID string, completed string) error {
	// Get existing task first
	task, err := getTaskByID(filename, taskID)
	if err != nil {
		return err
	}

	// Create updated task
	completedBool, err := strconv.ParseBool(completed)
	if err != nil {
		return fmt.Errorf("failed to parse completed value: %w", err)
	}
	updatedTask := Task{task.ID, task.Task, completedBool, time.Now()}

	// Read all tasks
	file, err := openCSVFile(filename, os.O_RDONLY)
	if err != nil {
		return err
	}
	tasks, err := readCSVFile(file)
	file.Close()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	// Create new file
	file, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header at the top
	if err := writer.Write([]string{"ID", "Task", "Completed", "CreatedAt"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write all tasks, replacing the updated one
	for _, t := range tasks {
		if t.ID == task.ID {
			record := []string{
				strconv.Itoa(updatedTask.ID),
				updatedTask.Task,
				strconv.FormatBool(updatedTask.Completed),
				updatedTask.CreatedAt.Format(time.RFC3339),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write task: %w", err)
			}
		} else {
			record := []string{
				strconv.Itoa(t.ID),
				t.Task,
				strconv.FormatBool(t.Completed),
				t.CreatedAt.Format(time.RFC3339),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write task: %w", err)
			}
		}
	}

	fmt.Println("Task updated successfully")
	return nil
}
