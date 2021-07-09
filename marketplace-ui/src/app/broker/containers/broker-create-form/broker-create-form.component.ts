import {Component, Inject, OnInit} from '@angular/core';

import {Broker} from '../../models/broker.interface';
import {BrokerContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {Router} from '@angular/router';

@Component(
  {
    selector: 'app-broker-create-form',
    templateUrl: './broker-create-form.component.html',
    styleUrls: ['broker-create-form.component.scss']
  }
)
export class BrokerCreateFormComponent implements OnInit {

  accounts: string[];
  brokerContractInstance: any;

  locations: string[] = [
    'Brazil', 'Europe Nordic & East', 'Europe West', 'Latin America North', 'Latin America South', 'North America',
    'Oceania', 'Russia', 'Turkey', 'Japan', 'The Philippines', 'Singapore, Malaysia, and Indonesia', 'Taiwan, Hong Kong, and Macao',
    'Vietnam', 'Thailand', 'Republic of Korea', 'China',
  ];

  constructor(private router: Router, @Inject(WEB3) private web3: Web3, @Inject(BrokerContract) private brokerContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.brokerContract.setProvider(this.web3.currentProvider);
    try {
      this.accounts = await this.web3.eth.getAccounts();
      this.brokerContractInstance = await this.brokerContract.deployed();
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

  async handleSubmit(broker: Broker, isValid: boolean) {
    if (isValid) {
      await this.brokerContractInstance
        .create(broker.addr, broker.name, broker.hostAddr, broker.location, {from: this.accounts[0]});
      await this.router.navigateByUrl('brokers/management');
    }
  }

}
