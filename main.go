package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/maragudk/goqite"
	"go-learn/db"
	"go-learn/helpers"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "goqite.db"
	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// create the instance of database migration
	d := db.NewDatabase(dbConn)
	// run the schema of goqite table, but check the database first
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("Database file %s does not exist", dbPath)
		// then migrate or create it
		d.Migration()
	}

	// Initialize the queue
	q := goqite.New(goqite.NewOpts{
		DB:   dbConn,
		Name: "jobs",
	})

	/**
	Perform queue operations
	Send this dummy message
	*/
	body := []byte("This is a queue message")
	for range make([]struct{}, 3) {
		err = helpers.SendToQueue(context.Background(), q, body)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	// receive the queue message
	bodyReceived, err := helpers.ReceiveFromQueue(context.Background(), q)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("Received message:", string(bodyReceived))

	// this is example of how to extend the timeout of the queue message
	err = helpers.ExtendMessageTimeout(context.Background(), q, "message_id", time.Second)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// this is the example of deleting message if we already receive it, so it doesn't get redelivered
	err = helpers.DeleteMessageFromQueue(context.Background(), q, "message_id")
	if err != nil {
		log.Fatalln(err.Error())
	}

	/**
	This is example of how background service keep stay alive to process the queue
	For that, we can create a background goroutine to continuously process queue messages
	*/
	go func() {
		for {
			// Receive a message from the queue
			m, err := q.Receive(context.Background())
			if err != nil {
				log.Println(err)
				continue
			}

			/**
			Then we can process the message, like calling an API or other business logic
			*/
			if m != nil {
				fmt.Println("Processing message:", string(m.Body))

				// Delete the message after processing
				if err := q.Delete(context.Background(), m.ID); err != nil {
					log.Println("Error deleting message:", err)
					continue
				}
			}

			time.Sleep(time.Second) // Add some delay before processing the next message
		}
	}()

	// Keep the main goroutine alive
	select {}
}
