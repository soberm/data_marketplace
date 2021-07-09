import {Component, Inject, OnInit} from '@angular/core';
import {Device} from '../../models/device.interface';
import {DeviceContract, UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {User} from '../../../user/models/user.interface';

@Component(
  {
    selector: 'app-device-management',
    templateUrl: 'device-management.component.html',
    styleUrls: ['device-management.component.scss']
  }
)
export class DeviceManagementComponent implements OnInit {

  devices: Device[];
  accounts: string[];
  deviceContractInstance: any;
  userContractInstance: any;

  constructor(
    @Inject(WEB3) private web3: Web3,
    @Inject(UserContract) private userContract: any,
    @Inject(DeviceContract) private deviceContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.devices = [];
    this.accounts = await this.web3.eth.getAccounts();

    this.userContract.setProvider(this.web3.currentProvider);
    this.deviceContract.setProvider(this.web3.currentProvider);

    try {

      this.userContractInstance = await this.userContract.deployed();
      this.deviceContractInstance = await this.deviceContract.deployed();

      const user: User = await this.userContractInstance.findByAddress.call(this.accounts[0]);

      for (const address of user.devices) {
        const device: Device = await this.deviceContractInstance.findByAddress.call(address);
        if (device.deleted) {
          continue;
        }
        this.devices = [...this.devices, device];
      }
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

  async onDeleteDevice(event: Device) {
    try {
      await this.deviceContractInstance.remove(event.addr, {from: this.accounts[0]});
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
      return;
    }
    this.devices = this.devices.filter((device: Device) => {
      return device.addr !== event.addr;
    });
  }

}
