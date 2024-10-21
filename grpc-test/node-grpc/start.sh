#!/bin/bash

if [ "$SERVICE_TYPE" = "central" ]; then
    echo "중앙 서버 시작"
    npx ts-node src/central_server.ts
else
    # echo "gRPC 클라이언트 시작"
    npx ts-node src/grpc_client.ts
fi
