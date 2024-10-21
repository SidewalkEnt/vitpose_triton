package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
	CollectorURL string
	ClientID     string
	TestDuration time.Duration
}

func parseFlags() Flags {
	var flags Flags
	flag.StringVar(&flags.ModelName, "m", "vitpose_ensemble", "Name of model being served. (Required)")
	flag.StringVar(&flags.ModelVersion, "x", "", "Version of model. Default: Latest Version.")
	flag.IntVar(&flags.BatchSize, "b", 4, "Batch size. Default: 4.")
	flag.StringVar(&flags.URL, "u", "34.47.107.11:8001", "Inference Server URL. Default: 34.47.107.11:8001")
	flag.StringVar(&flags.CollectorURL, "collector", "http://aggregator:8080/record-latency", "URL of the latency collector server.")
	flag.StringVar(&flags.ClientID, "id", os.Getenv("CLIENT_ID"), "Unique identifier for the client instance.")

	// TEST_DURATION 환경 변수를 읽어서 설정
	testDurationStr := os.Getenv("TEST_DURATION")
	if testDurationStr == "" {
		log.Fatal("TEST_DURATION environment variable is not set.")
	}
	testDurationSec, err := strconv.Atoi(testDurationStr)
	if err != nil || testDurationSec <= 0 {
		log.Fatalf("Invalid TEST_DURATION value: %s", testDurationStr)
	}
	flags.TestDuration = time.Duration(testDurationSec) * time.Second

	flag.Parse()
	return flags
}

func ModelInferRequest(client triton.GRPCInferenceServiceClient, modelName, modelVersion string, batchSize, imageHeight, imageWidth int, collectorURL, clientID string) time.Duration {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		log.Printf("InferRequest 처리 오류: %v", err)
		return -1
	}
	duration := time.Since(startTime)

	// 레이턴시 데이터를 집계 서버로 전송
	sendLatencyData(collectorURL, clientID, duration)

	return duration
}

func sendLatencyData(url, clientID string, latency time.Duration) {
	data := map[string]interface{}{
		"client_id": clientID,
		"latency":   latency,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling latency data: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending latency data to collector: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK response from collector: %s", resp.Status)
	}
}

func main() {
	FLAGS := parseFlags()
	fmt.Println("FLAGS:", FLAGS)

	conn, err := grpc.Dial(FLAGS.URL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to endpoint %s: %v", FLAGS.URL, err)
	}
	defer conn.Close()

	client := triton.NewGRPCInferenceServiceClient(conn)

	var wg sync.WaitGroup
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(FLAGS.TestDuration)

	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				latency := ModelInferRequest(client, FLAGS.ModelName, FLAGS.ModelVersion, FLAGS.BatchSize, 256, 192, FLAGS.CollectorURL, FLAGS.ClientID)
				if latency != -1 {
					// fmt.Printf("Inference latency: %v\n", latency)
				} else {
					fmt.Println("Inference failed.")
				}
			}()
		case <-timeout:
			fmt.Printf("%v 동안의 테스트가 완료되었습니다.\n", FLAGS.TestDuration)
			wg.Wait()
			return
		}
	}
}
