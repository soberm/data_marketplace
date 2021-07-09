import {Component, EventEmitter, Inject, Input, OnInit, Output} from '@angular/core';

import {Device} from '../../models/device.interface';
import {DeviceContract, UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {User} from '../../../user/models/user.interface';
import {ActivatedRoute, Params, Router} from '@angular/router';

@Component(
  {
    selector: 'app-device-update-form',
    templateUrl: './device-update-form.component.html',
    styleUrls: ['device-update-form.component.scss']
  }
)
export class DeviceUpdateFormComponent implements OnInit {

  details: Device;

  accounts: string[];
  deviceContractInstance: any;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    @Inject(WEB3) private web3: Web3,
    @Inject(DeviceContract) private deviceContract: any) {
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

    this.route.params.subscribe(async (data: Params) => {
      try {
        this.details = await this.deviceContractInstance.findByAddress.call(data.address);
      } catch (e) {
        if (e instanceof Error) {
          console.log(e.message);
        } else {
          console.log(e);
        }
      }
    });
  }

  async handleSubmit(device: Device, isValid: boolean) {
    if (isValid) {
      await this.deviceContractInstance
        .update(this.details.addr, device.name, device.description, this.web3.utils.hexToBytes(device.publicKey), {from: this.accounts[0]});
      await this.router.navigateByUrl('devices/management');
    }
  }

}
