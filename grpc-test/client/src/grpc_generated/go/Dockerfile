# 베이스 이미지 설정
FROM golang:1.23-alpine

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 파일 복사 및 빌드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# aggregator 빌드
RUN go build -o aggregator ./aggregator.go

# client-app 빌드
RUN go build -o client-app ./client.go

# 컨테이너 실행 시 entrypoint 설정
CMD ["sh", "-c", "if [ \"$APP_TYPE\" = \"aggregator\" ]; then ./aggregator; else ./client-app -collector http://aggregator:8080/record-latency; fi"]
