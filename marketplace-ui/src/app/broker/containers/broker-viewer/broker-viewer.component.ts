import {Component, Inject, OnInit} from '@angular/core';
import {Broker} from '../../models/broker.interface';
import {Location} from '../../models/location.enum';
import {BrokerContract, UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {User} from '../../../user/models/user.interface';

@Component(
  {
    selector: 'app-broker-viewer',
    templateUrl: 'broker-viewer.component.html',
    styleUrls: ['broker-viewer.component.scss']
  }
)
export class BrokerViewerComponent implements OnInit {

  brokers: Broker[];
  brokerContractInstance: any;

  constructor(@Inject(WEB3) private web3: Web3, @Inject(BrokerContract) private brokerContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.brokers = [];
    try {
      this.brokerContract.setProvider(this.web3.currentProvider);
      this.brokerContractInstance = await this.brokerContract.deployed();
      const count: number = await this.brokerContractInstance.count.call();
      for (let counter = 0; counter < count; counter++) {
        const broker: Broker = await this.brokerContractInstance.findByIndex.call(counter);
        if (broker.deleted) {
          continue;
        }
        this.brokers = [...this.brokers, broker];
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
