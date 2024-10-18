package main

import (
	"context"
	"fmt"
	"grpc_test/gen"
	"grpc_test/gen/genconnect"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	address     = "http://a100.sidewalkplay.top:8001"
	numUsers    = 1024
	numImages   = 32
	imageHeight = 192
	imageWidth  = 256
	duration    = 25 * time.Second
	modelName   = "vitpose_ensemble"
)

func createModelInferRequest() *gen.ModelInferRequest {
	request := &gen.ModelInferRequest{
		ModelName: modelName,
		Inputs: []*gen.ModelInferRequest_InferInputTensor{
			{
				Name:     "input",
				Datatype: "FP32",
				Shape:    []int64{1, 3, int64(imageHeight), int64(imageWidth)},
			},
		},
	}

	data := make([]float32, 3*imageHeight*imageWidth)
	for i := range data {
		data[i] = rand.Float32()
	}

	request.Inputs[0].Contents = &gen.InferTensorContents{
		Fp32Contents: data,
	}

	return request
}

func main() {
	client := genconnect.NewGRPCInferenceServiceClient(
		http.DefaultClient,
		address,
	)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for i := 0; i < numUsers; i++ {
		go func(userID int) {
			for j := 0; j < numImages; j++ {
				request := createModelInferRequest()
				response, err := client.ModelInfer(ctx, request)
				if err != nil {
					log.Printf("추론 요청 실패: %v", err)
					continue
				}
				log.Printf("사용자 %d, 이미지 %d 추론 완료: %v", userID, j, response.ModelName)
			}
		}(i)
	}

	time.Sleep(duration)
	fmt.Println("추론 완료")
}
