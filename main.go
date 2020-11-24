package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mfcekirdek/background-task-poc/dispatcher"
)

func main() {
	task1 := func() error {
		fmt.Printf("job1: performing work.\n")
		time.Sleep(2 * time.Second)
		fmt.Printf("Work done.\n")
		return nil
	}

	task2 := func() error {
		fmt.Printf("job2: performing work.\n")
		time.Sleep(4 * time.Second)
		fmt.Printf("Work done.\n")
		return nil
	}

	dispatcher := dispatcher.CreateNewDispatcher()
	dispatcher.AddTask(task1)
	dispatcher.AddTask(task2)
	dispatcher.Start(2)

	for {
		if dispatcher.Finished() {
			log.Println("All jobs finished")
			break
		}
	}

}
