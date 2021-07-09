import {Component, Inject, OnInit} from '@angular/core';

import {BrokerContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {ActivatedRoute, Params, Router} from '@angular/router';
import {Broker} from '../../models/broker.interface';

@Component(
  {
    selector: 'app-broker-update-form',
    templateUrl: './broker-update-form.component.html',
    styleUrls: ['broker-update-form.component.scss']
  }
)
export class BrokerUpdateFormComponent implements OnInit {

  details: Broker;

  accounts: string[];
  brokerContractInstance: any;

  locations: string[] = [
    'Brazil', 'Europe Nordic & East', 'Europe West', 'Latin America North', 'Latin America South', 'North America',
    'Oceania', 'Russia', 'Turkey', 'Japan', 'The Philippines', 'Singapore, Malaysia, and Indonesia', 'Taiwan, Hong Kong, and Macao',
    'Vietnam', 'Thailand', 'Republic of Korea', 'China',
  ];

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    @Inject(WEB3) private web3: Web3,
    @Inject(BrokerContract) private brokerContract: any) {
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

    this.route.params.subscribe(async (data: Params) => {
      try {
        this.details = await this.brokerContractInstance.findByAddress.call(data.address);
        console.log(this.details);
      } catch (e) {
        if (e instanceof Error) {
          console.log(e.message);
        } else {
          console.log(e);
        }
      }
    });
  }

  async handleSubmit(broker: Broker, isValid: boolean) {
    if (isValid) {
      await this.brokerContractInstance
        .update(this.details.addr, broker.name, broker.hostAddr, broker.location, {from: this.accounts[0]});
      await this.router.navigateByUrl('brokers/management');
    }
  }

}
