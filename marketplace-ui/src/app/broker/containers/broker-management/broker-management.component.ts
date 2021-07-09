import {Component, Inject, OnInit} from '@angular/core';
import {Broker} from '../../models/broker.interface';
import {BrokerContract, DeviceContract, UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {User} from '../../../user/models/user.interface';

@Component(
  {
    selector: 'app-broker-management',
    templateUrl: 'broker-management.component.html',
    styleUrls: ['broker-management.component.scss']
  }
)
export class BrokerManagementComponent implements OnInit {

  brokers: Broker[];
  accounts: string[];
  brokerContractInstance: any;
  userContractInstance: any;

  constructor(
    @Inject(WEB3) private web3: Web3,
    @Inject(UserContract) private userContract: any,
    @Inject(BrokerContract) private brokerContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.brokers = [];
    this.accounts = await this.web3.eth.getAccounts();

    this.userContract.setProvider(this.web3.currentProvider);
    this.brokerContract.setProvider(this.web3.currentProvider);

    try {

      this.userContractInstance = await this.userContract.deployed();
      this.brokerContractInstance = await this.brokerContract.deployed();

      const user: User = await this.userContractInstance.findByAddress.call(this.accounts[0]);

      for (const address of user.brokers) {
        const broker: Broker = await this.brokerContractInstance.findByAddress.call(address);
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

  async onDeleteBroker(event: Broker) {
    try {
      await this.brokerContractInstance.remove(event.addr, {from: this.accounts[0]});
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
      return;
    }
    this.brokers = this.brokers.filter((broker: Broker) => {
      return broker.addr !== event.addr;
    });
  }

}
