export interface Trade {
  id: number;
  provider: string;
  consumer: string;
  broker: string;
  product: number;
  startTime: number;
  endTime: number;
  cost: number;
  settlementContract: string;
  settled: boolean;
}
