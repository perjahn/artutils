package main

import (
	"fmt"
	"net/http"
	"slices"
	"time"
)

func ShowTasksOnce(baseurl string, token string, dontShowScheduled bool) error {
	client := &http.Client{}

	tasks, err := GetTasks(client, baseurl, token)
	if err != nil {
		return err
	}

	headers := []string{"Task ID", "Type", "State", "Description", "Started"}
	lengths := make([]int, len(headers)-1)
	for i := range len(headers) - 1 {
		lengths[i] = len(headers[i])
	}

	for _, task := range tasks {
		if dontShowScheduled && task.State != nil && *task.State == "scheduled" {
			continue
		}
		if task.Id != nil && len(*task.Id) > lengths[0] {
			lengths[0] = len(*task.Id)
		}
		if task.Type != nil && len(*task.Type) > lengths[1] {
			lengths[1] = len(*task.Type)
		}
		if task.State != nil && len(*task.State) > lengths[2] {
			lengths[2] = len(*task.State)
		}
		if task.Description != nil && len(*task.Description) > lengths[3] {
			lengths[3] = len(*task.Description)
		}
	}

	fmt.Printf("%-*s  %-*s  %-*s  %-*s  %s\n",
		lengths[0], headers[0],
		lengths[1], headers[1],
		lengths[2], headers[2],
		lengths[3], headers[3],
		headers[4])

	for _, task := range tasks {
		if dontShowScheduled && *task.State == "scheduled" {
			continue
		}
		fmt.Printf("%-*s  %-*s  %-*s  %-*s  %s\n",
			lengths[0], *task.Id,
			lengths[1], *task.Type,
			lengths[2], *task.State,
			lengths[3], *task.Description,
			task.Started)
	}

	return nil
}

type Task struct {
	Id          string
	Type        string
	State       string
	Description string
	Started     time.Time
	Duration    time.Duration
	Ended       *time.Time
}

// For every iteration, print all task history. A task in the history can have the Ended set to null, i.e. it's still running.
// Multiple tasks with the same Id can exists in the task history.

func ShowTasksFollow(baseurl string, token string) error {
	client := &http.Client{}

	var itercount int = 0
	var taskhistory []Task

	for {
		tasks, err := GetTasks(client, baseurl, token)
		if err != nil {
			return err
		}

		// For each task in tasks, that has just Started, add it to taskhistory
		for _, task := range tasks {
			if task.Started != nil && task.Id != nil {
				index := slices.IndexFunc(taskhistory, func(t Task) bool {
					return t.Id == *task.Id && t.Started.Equal(*task.Started) && t.Ended == nil
				})
				if index == -1 {
					taskhistory = append(taskhistory, Task{
						Id:          *task.Id,
						Type:        *task.Type,
						State:       *task.State,
						Description: *task.Description,
						Started:     *task.Started,
						Ended:       nil,
					})
				}
			}
		}

		// For each task in taskhistory, check if any has been terminated (doesn't exists in tasks any longer), and set Ended to current time
		for i := range taskhistory {
			task := &taskhistory[i]

			now := time.Now()

			if task.Ended == nil {
				task.Duration = now.Sub(task.Started)
			}
			exists := slices.ContainsFunc(tasks, func(t ArtifactoryTask) bool {
				return t.Id != nil && *t.Id == task.Id && t.Started != nil && t.Started.Equal(task.Started)
			})
			if !exists && task.Ended == nil {
				task.Ended = &now
				task.State = "-      "
			}
		}

		headers := []string{"Task ID", "State", "Description", "Started", "Duration", "Ended"}
		lengths := make([]int, len(headers)-1)
		for i := range len(headers) - 1 {
			lengths[i] = len(headers[i])
		}

		for _, task := range taskhistory {
			if len(task.Id) > lengths[0] {
				lengths[0] = len(task.Id)
			}
			if len(task.State) > lengths[1] {
				lengths[1] = len(task.State)
			}
			if len(task.Description) > lengths[2] {
				lengths[2] = len(task.Description)
			}
			if len(task.Started.Local().Format("2006-01-02 15:04:05")) > lengths[3] {
				lengths[3] = len(task.Started.Local().Format("2006-01-02 15:04:05"))
			}
			if len(task.Duration.String()) > lengths[4] {
				lengths[4] = len(task.Duration.String())
			}
		}

		fmt.Print("\033[2J\033[3J\033[H")

		fmt.Printf("%-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
			lengths[0], headers[0],
			lengths[1], headers[1],
			lengths[2], headers[2],
			lengths[3], headers[3],
			lengths[4], headers[4],
			headers[5])

		for _, task := range taskhistory {
			var ended string
			if task.Ended == nil {
				ended = "-"
			} else {
				ended = task.Ended.Format("2006-01-02 15:04:05")
			}

			fmt.Printf("%-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
				lengths[0], task.Id,
				lengths[1], task.State,
				lengths[2], task.Description,
				lengths[3], task.Started.Local().Format("2006-01-02 15:04:05"),
				lengths[4], task.Duration,
				ended)
		}

		fmt.Printf("\nIteration: %d\n", itercount+1)
		itercount++

		time.Sleep(100 * time.Millisecond)
	}
}
