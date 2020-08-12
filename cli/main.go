package main

import (
	"context"
	example "github.com/jwenz723/podlifecycle/server/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "192.168.64.24:30194"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := example.NewExampleClient(conn)

	for {
		time.Sleep(100 * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		r, err := c.Work(ctx, &example.WorkItem{
			Name: "test",
			Size: 0,
		})
		if err != nil {
			log.Printf("could not Work: %v", err)
			continue
		}
		log.Printf("Response: %s", r.GetName())
	}
}
