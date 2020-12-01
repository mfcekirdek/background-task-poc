package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mfcekirdek/background-task-poc/dispatcher"
)

var Dispatcher dispatcher.Dispatcher

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

	Dispatcher := dispatcher.CreateNewDispatcher()
	Dispatcher.AddTask(task1)
	Dispatcher.AddTask(task2)
	Dispatcher.Start(2)

	for {
		if Dispatcher.Finished() {
			log.Println("All jobs finished")
			break
		}
	}

}
