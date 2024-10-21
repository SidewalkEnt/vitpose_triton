package main

import (
	"grpc_test/gen"
	"math/rand"
	"os"

	"google.golang.org/protobuf/proto"
)

const (
	batchSize   = 1
	imageHeight = 256
	imageWidth  = 192
	modelName   = "vitpose_ensemble"
)

func createModelInferRequest() *gen.ModelInferRequest {
	request := &gen.ModelInferRequest{
		ModelName: modelName,
		Inputs: []*gen.ModelInferRequest_InferInputTensor{
			{
				Name:     "input",
				Datatype: "FP32",
				Shape:    []int64{batchSize, 3, int64(imageHeight), int64(imageWidth)},
			},
		},
	}

	data := make([]float32, batchSize*3*imageHeight*imageWidth)
	for i := range data {
		data[i] = rand.Float32()
	}

	request.Inputs[0].Contents = &gen.InferTensorContents{
		Fp32Contents: data,
	}

	return request
}

func main() {
	// ModelInferRequest 생성
	request := createModelInferRequest()

	// 프로토콜 버퍼 메시지를 바이너리로 직렬화
	data, err := proto.Marshal(request)
	if err != nil {
		panic(err)
	}

	// 바이너리 데이터를 파일에 저장
	file, err := os.Create("input_data.bin")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}

	println("Binary data has been written to input_data.bin")
}
