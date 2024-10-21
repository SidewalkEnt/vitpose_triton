package main

import (
	"bytes"
	"context"
	"encoding/binary"

	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"

	// "net/http"
	"os"
	"strconv"
	"sync"
	"time"

	triton "github.com/triton-inference-server/client/src/grpc_generated/go/grpc-client"
	"google.golang.org/grpc"
)

type Flags struct {
	ModelName    string
	ModelVersion string
	BatchSize    int
	URL          string
	TestDuration time.Duration
}

func parseFlags() Flags {
	var flags Flags
	flag.StringVar(&flags.ModelName, "m", "vitpose_ensemble", "Name of model being served. (Required)")
	flag.StringVar(&flags.ModelVersion, "x", "", "Version of model. Default: Latest Version.")
	flag.IntVar(&flags.BatchSize, "b", 4, "Batch size. Default: 4.")
	flag.StringVar(&flags.URL, "u", "34.47.107.11:8001", "Inference Server URL. Default: 34.47.107.11:8001")

	// TEST_DURATION 환경 변수를 읽어서 설정
	testDurationStr := os.Getenv("TEST_DURATION")
	if testDurationStr == "" {
		testDurationStr = "60" // 기본값: 60초
	}
	testDurationSec, err := strconv.Atoi(testDurationStr)
	if err != nil || testDurationSec <= 0 {
		log.Fatalf("Invalid TEST_DURATION value: %s", testDurationStr)
	}
	flags.TestDuration = time.Duration(testDurationSec) * time.Second

	flag.Parse()
	return flags
}

func ModelInferRequest(client triton.GRPCInferenceServiceClient, modelName, modelVersion string, batchSize, imageHeight, imageWidth int, clientID string) time.Duration {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dataSize := batchSize * 3 * imageHeight * imageWidth
	data := make([]float32, dataSize)
	for i := range data {
		data[i] = rand.Float32()
	}

	inferInputs := []*triton.ModelInferRequest_InferInputTensor{
		{
			Name:     "input",
			Datatype: "FP32",
			Shape:    []int64{int64(batchSize), 3, int64(imageHeight), int64(imageWidth)},
		},
	}

	modelInferRequest := triton.ModelInferRequest{
		ModelName:    modelName,
		ModelVersion: modelVersion,
		Inputs:       inferInputs,
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		log.Fatalf("바이너리 데이터 쓰기 오류: %v", err)
	}
	modelInferRequest.RawInputContents = append(modelInferRequest.RawInputContents, buf.Bytes())

	startTime := time.Now()
	_, err = client.ModelInfer(ctx, &modelInferRequest)
	if err != nil {
		log.Printf("Client %s: InferRequest 처리 오류: %v", clientID, err)
		return -1
	}
	duration := time.Since(startTime)

	return duration
}

func simulateClient(clientID string, client triton.GRPCInferenceServiceClient, flags Flags, wg *sync.WaitGroup, stopCh chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-ticker.C:
			latency := ModelInferRequest(client, flags.ModelName, flags.ModelVersion, flags.BatchSize, 256, 192, clientID)
			if latency != -1 {
				fmt.Printf("Client %s: Inference latency: %v\n", clientID, latency)
			} else {
				fmt.Printf("Client %s: Inference failed.\n", clientID)
			}
		case <-stopCh:
			return
		}
	}
}

func main() {
	FLAGS := parseFlags()

	conn, err := grpc.Dial(FLAGS.URL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to endpoint %s: %v", FLAGS.URL, err)
	}
	defer conn.Close()

	client := triton.NewGRPCInferenceServiceClient(conn)

	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	// 128명의 유저를 시뮬레이션합니다.
	for i := 0; i < 1; i++ {
		clientID := fmt.Sprintf("client-%d", i+1)
		wg.Add(1)
		go simulateClient(clientID, client, FLAGS, &wg, stopCh)
	}

	// 테스트 지속 시간 동안 대기합니다.
	time.Sleep(FLAGS.TestDuration)

	// 모든 클라이언트 고루틴을 종료합니다.
	close(stopCh)
	wg.Wait()

	fmt.Println("모든 클라이언트의 테스트가 완료되었습니다.")
}
