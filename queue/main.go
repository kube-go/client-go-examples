package main

import (
	"fmt"
	"k8s.io/client-go/util/workqueue"
	"time"
)

// This function will add 5 words to queue and then exit
func writeToQueue(queue workqueue.Interface) {
	fmt.Println("Started Writer goroutine")
	// Adding sleep to test behaviour of read
	time.Sleep(3 * time.Second)
	queue.Add("one")
	queue.Add("two")
	queue.Add("three")
	queue.Add("four")
	queue.Add("five")
	fmt.Println("Adding to queue completed!")
	// This is needed to let the reader know writing to queue is completed
	queue.ShutDown()
}

// Read from the queue and print results
func readFromQueue(queue workqueue.Interface, stop chan int) {
	fmt.Println("Started Reader goroutine")
	for {
		fmt.Println("Reader loop")
		// Get blocks until it can return an item to be processed
		item, shutdown := queue.Get()
		fmt.Println("Trying to retrieve from queue")
		if shutdown {
			// signal that done reading queue
			stop <- -1
			return
		}
		fmt.Printf("Got item[shutdown = %t]: %s\n", shutdown, item)
		// Mark as done processing
		queue.Done(item)
	}
}

func main() {
	fmt.Println("Starting main")
	// Create a channel
	stop := make(chan int)
	// Create a queue.
	myQueue := workqueue.New()
	// Create our first worker thread.  This goroutine will
	// simply put five words into the queue after one second
	// has passed
	go writeToQueue(myQueue)
	// Create second thread that will read from the queue
	go readFromQueue(myQueue, stop)
	fmt.Println("Goroutines started, now waiting for reader to complete")
	// To avoid exiting main before finishing go routines
	<-stop
	fmt.Println("Reader signaled completion, exiting")
}
