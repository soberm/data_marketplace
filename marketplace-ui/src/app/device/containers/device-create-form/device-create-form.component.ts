import {Component, EventEmitter, Inject, OnInit, Output} from '@angular/core';

import {Device} from '../../models/device.interface';
import {DeviceContract, UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {User} from '../../../user/models/user.interface';
import {Router} from '@angular/router';

@Component(
  {
    selector: 'app-device-create-form',
    templateUrl: './device-create-form.component.html',
    styleUrls: ['device-create-form.component.scss']
  }
)
export class DeviceCreateFormComponent implements OnInit {

  accounts: string[];
  deviceContractInstance: any;

  constructor(private router: Router, @Inject(WEB3) private web3: Web3, @Inject(DeviceContract) private deviceContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.deviceContract.setProvider(this.web3.currentProvider);
    try {
      this.accounts = await this.web3.eth.getAccounts();
      this.deviceContractInstance = await this.deviceContract.deployed();
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

  async handleSubmit(device: Device, isValid: boolean) {
    if (isValid) {
      await this.deviceContractInstance
        .create(device.addr, device.name, device.description, this.web3.utils.hexToBytes(device.publicKey), {from: this.accounts[0]});
      await this.router.navigateByUrl('devices/management');
    }
  }

}
