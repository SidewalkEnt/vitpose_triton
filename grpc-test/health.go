package health

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"

	pb "grpc_test/gen"
	greetv1 "grpc_test/gen/genconnect"
)

func healthCheck() {
	client := greetv1.NewHealthClient(
		http.DefaultClient,
		"http://34.47.107.11:8001",
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var res *connect.Response[pb.HealthCheckResponse]
	var err error

	for retries := 0; retries < 3; retries++ {
		res, err = client.Check(
			ctx,
			connect.NewRequest(&pb.HealthCheckRequest{}),
		)
		if err == nil {
			break
		}
		log.Printf("Error on attempt %d: %v", retries+1, err)
		if retries < 2 {
			log.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("Failed after 3 attempts: %v", err)
	}

	log.Printf("Response: %+v", res.Msg)
}
