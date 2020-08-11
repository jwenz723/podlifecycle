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
	for {
		time.Sleep(1 * time.Second)
		// Set up a connection to the server.
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Printf("did not connect: %v", err)
			continue
		}
		defer conn.Close()
		c := podlifecycle.NewStufferClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.DoStuff(ctx, &podlifecycle.StuffRequest{Name: "test"})
		if err != nil {
			log.Printf("could not DoStuff: %v", err)
			continue
		}
		log.Printf("Response: %s", r.GetName())
	}
}
