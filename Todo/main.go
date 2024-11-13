/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofrs/flock"
	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

var tasksFile ="tasks.csv"

type Task struct {
	ID int
	Task string
	Completed bool
	CreatedAt time.Time
}

func main() {

	var rootCmd = &cobra.Command{Use: "tasks"}

	rootCmd.AddCommand(addCmd, listCmd, completeCmd, deleteCmd, getCmd)

	rootCmd.Execute()
}


var addCmd = &cobra.Command{
	Use: "add",
	Short: "Add a new task to the list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task to add")
			return
		}
		addTask(args[0])
	},
}

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		listTasks()
	},
}
var completeCmd = &cobra.Command{
	Use: "complete",
	Short: "Complete a task",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task ID to complete")
			return
		}
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		completeTask(taskID)
	},
}

var deleteCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete a task",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task ID to delete")
			return
		}
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		deleteTask(taskID)
	},
}
var getCmd = &cobra.Command{
	Use: "get",
	Short: "Get a task",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task ID to get")
			return
		}
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		getTask(taskID)
	},
}

func addTask(task string)  {
	file, err := loadFile(tasksFile)
	if err != nil {
		fmt.Printf("failed to load file: %v\n", err)
		return
	}

	defer closeFile(file)

	tasks, err := readTasks(file)
	if err != nil {
		fmt.Printf("failed to read tasks: %v\n", err)
		return
	}

	
	newTask := Task{ID: len(tasks) + 1, Task: task, Completed: false, CreatedAt: time.Now().UTC()}
	tasks = append(tasks, newTask)
	 writeTasks(file, tasks)
	 fmt.Println("Task added successfully")
}

func listTasks() {
	file, err := loadFile(tasksFile)
	if err != nil {
		fmt.Printf("failed to load file: %v\n", err)
		return
	}

	defer closeFile(file)

	tasks, err := readTasks(file)
	if err != nil {
		fmt.Printf("failed to read tasks: %v\n", err)
		return
	}

	fmt.Println("ID. Task Status Created At")
	for _, task := range tasks {
		fmt.Printf("%d. %s %v %s\n", task.ID, task.Task, task.Completed, timediff.TimeDiff(task.CreatedAt))
	}
}

func completeTask(taskID int) {
	file, err := loadFile(tasksFile)
	if err != nil {
		fmt.Printf("failed to load file: %v\n", err)
		return
	}

	defer closeFile(file)

	tasks, err := readTasks(file)
	if err != nil {
		fmt.Printf("failed to read tasks: %v\n", err)
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks[i].Completed = true
			break
		}
	}

	writeTasks(file, tasks)
	fmt.Println("Task completed successfully")
}

func deleteTask(taskID int) {
	file, err := loadFile(tasksFile)
	if err != nil {
		fmt.Printf("failed to load file: %v\n", err)
		return
	}

	defer closeFile(file)

	tasks, err := readTasks(file)
	if err != nil {
		fmt.Printf("failed to read tasks: %v\n", err)
		return
	}

	var updatedTasks []Task

	for _, task := range tasks {
		if task.ID != taskID {
			updatedTasks = append(updatedTasks, task)
		}
	}

	writeTasks(file, updatedTasks)
	fmt.Println("Task deleted successfully")
}

func getTask(taskID int)  {
	file, err := loadFile(tasksFile)
	if err != nil {
		fmt.Printf("failed to load file: %v\n", err)
		return
	}

	defer closeFile(file)

	tasks, err := readTasks(file)
	if err != nil {
		fmt.Printf("failed to read tasks: %v\n", err)
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			fmt.Println("ID. Task Status Created At")
			fmt.Printf("%d. %s %v %s\n", task.ID, task.Task, task.Completed, timediff.TimeDiff(task.CreatedAt))
		}
	}
}

func loadFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	fileLock := flock.New(filepath + ".lock")

	locked, err := fileLock.TryLock()
	if err != nil {
		return nil, fmt.Errorf("failed to lock file: %w", err)
	}

	if !locked {
		return nil, fmt.Errorf("failed to lock file")
	}


	return f, nil
}

func closeFile(f *os.File) error {
	fileLock := flock.New(f.Name() + ".lock")
	return fileLock.Unlock()
}

func readTasks(file *os.File) ([]Task, error) {
	file.Seek(0,0)
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	var tasks []Task

	for _, record := range records {
		id, _ := strconv.Atoi(record[0])
		isCompleted, _ := strconv.ParseBool(record[3])
		createdAt, _ := time.Parse(time.RFC3339, record[2])
		tasks = append(tasks, Task{ID: id, Task: record[1], Completed: isCompleted, CreatedAt: createdAt})
	}

	return tasks, nil
}

func writeTasks(file *os.File, tasks []Task) error {
	file.Truncate(0)
	file.Seek(0,0)

	writer := csv.NewWriter(file)

	for _, task := range tasks {
		record := []string{strconv.Itoa(task.ID), task.Task, task.CreatedAt.Format(time.RFC3339), strconv.FormatBool(task.Completed)}
		writer.Write(record)
	}

	writer.Flush()
	return nil
}