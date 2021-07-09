import {Component, Inject, OnInit} from '@angular/core';
import {Trade} from '../../models/trade.interface';
import {SettlementContract, TradingContract, WEB3} from '../../../web3';
import Web3 from 'web3';

@Component(
  {
    selector: 'app-trade-viewer',
    templateUrl: 'trade-viewer.component.html',
    styleUrls: ['trade-viewer.component.scss']
  }
)
export class TradeViewerComponent implements OnInit {

  trades: Trade[];
  tradingContractInstance: any;

  constructor(
    @Inject(WEB3) private web3: Web3,
    @Inject(TradingContract) private tradingContract: any,
    @Inject(SettlementContract) private settlementContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.trades = [];
    try {
      this.tradingContract.setProvider(this.web3.currentProvider);
      this.settlementContract.setProvider(this.web3.currentProvider);
      this.tradingContractInstance = await this.tradingContract.deployed();
      const count: number = await this.tradingContractInstance.countTrades.call();
      for (let counter = 0; counter < count; counter++) {
        const trade: Trade = await this.tradingContractInstance.findTradeByIndex.call(counter);
        const settlementContractInstance = await this.settlementContract.at(trade.settlementContract);
        trade.settled = await settlementContractInstance.settled.call();
        this.trades = [...this.trades, trade];
      }
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

}
