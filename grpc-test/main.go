package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc_test/gen"
	genconnect "grpc_test/gen/genconnect"

	"google.golang.org/grpc"
)

func main() {
	// gRPC 서버 주소 설정
	address := "34.47.107.11:8001"

	// gRPC 연결 설정
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// 헬스 체크 클라이언트 생성
	healthClient := genconnect.NewHealthClient(conn)

	// 헬스 체크 요청 생성
	request := &pb.HealthCheckRequest{
		Service: "", // 기본 서비스 헬스 체크를 수행하기 위해 빈 문자열 사용
	}

	// 헬스 체크 요청 수행
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := healthClient.Check(ctx, request)
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}

	// 헬스 체크 응답 상태 출력
	switch response.Status {
	case pb.HealthCheckResponse_SERVING:
		fmt.Println("The server is serving.")
	case pb.HealthCheckResponse_NOT_SERVING:
		fmt.Println("The server is not serving.")
	case pb.HealthCheckResponse_UNKNOWN:
		fmt.Println("The server status is unknown.")
	case pb.HealthCheckResponse_SERVICE_UNKNOWN:
		fmt.Println("The specified service is unknown.")
	default:
		fmt.Println("Unknown status received.")
	}
}
