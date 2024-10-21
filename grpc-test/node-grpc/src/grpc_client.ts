import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import * as path from 'path';
import axios from 'axios';

const PROTO_PATH = path.resolve(__dirname, '../protobuf/grpc_service.proto');

const batchSize = 4;
const imageHeight = 256;
const imageWidth = 192;
const modelName = "vitpose_ensemble";

interface ModelInferRequest {
    model_name: string;
    inputs: {
        name: string;
        datatype: string;
        shape: number[];
        contents: {
            fp32_contents: number[];
        };
    }[];
}

interface GrpcError {
    error: grpc.ServiceError;
    latency: number;
}

function createModelInferRequest(): ModelInferRequest {
    const request: ModelInferRequest = {
        model_name: modelName,
        inputs: [{
            name: "input",
            datatype: "FP32",
            shape: [batchSize, 3, imageHeight, imageWidth],
            contents: {
                fp32_contents: []
            }
        }]
    };

    const dataSize = batchSize * 3 * imageHeight * imageWidth;
    const data = new Array(dataSize);
    for (let i = 0; i < dataSize; i++) {
        data[i] = Math.random();
    }

    request.inputs[0].contents.fp32_contents = data;

    return request;
}

function makeGrpcCall(client: grpc.Client, request: ModelInferRequest): Promise<{ response: any; latency: number }> {
    return new Promise((resolve, reject) => {
        const startTime = process.hrtime();
        (client as any).ModelInfer(request, (error: grpc.ServiceError | null, response: any) => {
            const endTime = process.hrtime(startTime);
            const latency = endTime[0] * 1000 + endTime[1] / 1000000; // 밀리초로 변환
            if (error) {
                reject({ error, latency });
            } else {
                resolve({ response, latency });
            }
        });
    });
}

async function simulateUser(client: grpc.Client) {
    const totalDuration = 20; // 총 실행 시간 (초)
    const startTime = Date.now();
    const latencies: number[] = [];

    while ((Date.now() - startTime) / 1000 < totalDuration) {
        const request = createModelInferRequest();
        try {
            const { response, latency } = await makeGrpcCall(client, request);
            // console.log(`요청 완료. 지연 시간: ${latency.toFixed(2)}ms`);
            latencies.push(latency);
        } catch (error) {
            const grpcError = error as GrpcError;
            console.error(`오류 발생. 지연 시간: ${grpcError.latency.toFixed(2)}ms`, grpcError.error);
            latencies.push(grpcError.latency);
        }
        await new Promise(resolve => setTimeout(resolve, 1000)); // 1초 대기
    }

    const sortedLatencies = latencies.sort((a, b) => a - b);
    const averageLatency = sortedLatencies.reduce((a, b) => a + b, 0) / sortedLatencies.length;
    const p50 = sortedLatencies[Math.floor(sortedLatencies.length * 0.5)];
    const p95 = sortedLatencies[Math.floor(sortedLatencies.length * 0.95)];

    // 결과를 중앙 서버로 전송
    try {
        const response = await axios.post('http://central-server:3000/report', {
            averageLatency,
            p50,
            p95
        });
        // console.log('결과를 중앙 서버로 전송했습니다. 서버 응답:', response.data);
    } catch (error) {
        console.error('중앙 서버로 결과 전송 실패:', error);
    }
}

async function main() {
    const packageDefinition = await protoLoader.load(PROTO_PATH, {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });

    const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
    const inferenceService = (protoDescriptor.inference as any).GRPCInferenceService as grpc.ServiceClientConstructor;

    const client = new inferenceService('34.47.107.11:8001', grpc.credentials.createInsecure());

    await simulateUser(client);
}

main().catch(console.error);
