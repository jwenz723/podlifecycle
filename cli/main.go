package main

import (
	"context"
	podlifecycle "github.com/jwenz723/podlifecycle/server/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:8080"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := podlifecycle.NewStufferClient(conn)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.DoStuff(ctx, &podlifecycle.StuffRequest{Name: "test"})
		if err != nil {
			log.Fatalf("could not DoStuff: %v", err)
		}
		log.Printf("Response: %s", r.GetName())
		time.Sleep(1 * time.Second)
	}
}
