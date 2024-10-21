package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type LatencyData struct {
	ClientID string        `json:"client_id"`
	Latency  time.Duration `json:"latency"`
}

var (
	latencyRecords  = make(map[string][]time.Duration)
	totalRequests   int
	currentRequests int
	mu              sync.Mutex
	done            = make(chan bool)
)

func init() {
	totalRequestsStr := os.Getenv("TOTAL_REQUESTS")
	if totalRequestsStr == "" {
		log.Fatal("TOTAL_REQUESTS environment variable is not set.")
	}
	count, err := strconv.Atoi(totalRequestsStr)
	if err != nil || count <= 0 {
		log.Fatalf("Invalid TOTAL_REQUESTS value: %s", totalRequestsStr)
	}
	totalRequests = count
}

// recordLatency 수집된 레이턴시 데이터를 저장하는 핸들러
func recordLatency(w http.ResponseWriter, r *http.Request) {
	var data LatencyData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mu.Lock()
	latencyRecords[data.ClientID] = append(latencyRecords[data.ClientID], data.Latency)
	currentRequests++
	if currentRequests >= totalRequests {
		go showResults()
	}
	mu.Unlock()

	fmt.Fprintf(w, "Latency recorded for client: %s", data.ClientID)
}

// calculatePercentile 퍼센타일 계산 함수
func calculatePercentile(latencies []time.Duration, percentile float64) time.Duration {
	index := int(float64(len(latencies)) * percentile)
	if index >= len(latencies) {
		index = len(latencies) - 1
	}
	return latencies[index]
}

// showResults 수집된 결과를 읽기 쉽게 출력하는 함수
func showResults() {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Final Results:")
	fmt.Println("---------------------------------------------------------")
	fmt.Printf("%-10s | %-10s | %-10s | %-10s | %-10s\n", "ClientID", "Min", "Max", "P50", "P95")
	fmt.Println("---------------------------------------------------------")

	for clientID, latencies := range latencyRecords {
		if len(latencies) == 0 {
			continue
		}
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

		min := latencies[0]
		max := latencies[len(latencies)-1]
		p50 := calculatePercentile(latencies, 0.50)
		p95 := calculatePercentile(latencies, 0.95)

		fmt.Printf("%-10s | %-10d | %-10d | %-10d | %-10d\n", clientID, min.Milliseconds(), max.Milliseconds(), p50.Milliseconds(), p95.Milliseconds())
	}
	fmt.Println("---------------------------------------------------------")

	done <- true
}

func main() {
	http.HandleFunc("/record-latency", recordLatency)

	fmt.Println("Starting latency collection server on port 8080...")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-done
	fmt.Println("Server has finished collecting data and displayed results.")
}
