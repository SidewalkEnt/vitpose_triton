import express, { Request, Response } from 'express';
import bodyParser from 'body-parser';

const app = express();
app.use(bodyParser.json());

interface ReportData {
  replicaId: string;
  averageLatency: number;
  p50: number;
  p95: number;
}

let results: ReportData[] = [];
let expectedReplicas = parseInt(process.env.EXPECTED_REPLICAS || '1', 10);
let completedReplicas = 0;

app.post('/report', (req: Request, res: Response) => {
    const reportData: ReportData = req.body;
    results.push(reportData);
    completedReplicas++;
    // console.log(`Replica ${reportData.replicaId} 결과 수신:`, reportData);
    // console.log(`완료된 replica: ${completedReplicas}/${expectedReplicas}`);
    
    if (completedReplicas === expectedReplicas) {
        console.log('모든 replica가 완료되었습니다. 최종 통계를 계산합니다.');
        calculateFinalStats();
    }
    
    // res.json({ message: '결과가 성공적으로 저장되었습니다.' });
});

app.get('/stats', (req: Request, res: Response) => {
    if (completedReplicas < expectedReplicas) {
        res.json({ message: `아직 모든 replica가 완료되지 않았습니다. (${completedReplicas}/${expectedReplicas})` });
    } else {
        res.json(calculateFinalStats());
    }
});

function calculateFinalStats() {
    const totalReplicas = results.length;
    const overallP50 = calculatePercentile(results.map(r => r.p50), 50);
    const overallP95 = calculatePercentile(results.map(r => r.p95), 95);
    const overallAverage = results.reduce((sum, r) => sum + r.averageLatency, 0) / totalReplicas;

    const stats = {
        totalReplicas,
        overallAverage,
        overallP50,
        overallP95,
        replicaResults: results
    };

    console.log('최종 통계:', stats);
    return stats;
}

function calculatePercentile(values: number[], percentile: number): number {
    const sorted = values.sort((a, b) => a - b);
    const index = Math.ceil((percentile / 100) * sorted.length) - 1;
    return sorted[index];
}

const port = 3000;
app.listen(port, () => {
    console.log(`중앙 서버가 ${port} 포트에서 실행 중입니다.`);
    console.log(`예상되는 총 replica 수: ${expectedReplicas}`);
});
