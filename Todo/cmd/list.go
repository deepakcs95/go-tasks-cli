package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// Task struct to represent a task
type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks in your todo list",
	Long: `The list command allows you to view all tasks in your todo list.
You can see the task descriptions and other details.
For example:

todo-cli list

This will list all tasks in your todo list.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Example tasks slice
		tasks := []Task{
			{ID: 1, Description: "Tidy up my desk", CreatedAt: time.Now().Add(-10 * time.Minute), IsComplete: false},
			{ID: 2, Description: "Write documentation", CreatedAt: time.Now().Add(-1 * time.Hour), IsComplete: true},
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTask\tCreated\tDone")
		for _, task := range tasks {
			fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", task.ID, task.Description, customTimeAgo(task.CreatedAt), task.IsComplete)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// customTimeAgo formats the time difference in a more granular way
func customTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration.Seconds() < 60 {
		return fmt.Sprintf("%.0f seconds ago", duration.Seconds())
	} else if duration.Minutes() < 60 {
		return fmt.Sprintf("%.0f minutes ago", duration.Minutes())
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	} else if duration.Hours() < 48 {
		return "yesterday"
	} else {
		return fmt.Sprintf("%.0f days ago", duration.Hours()/24)
	}
}