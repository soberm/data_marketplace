import {Component, Inject, OnInit} from '@angular/core';
import {Device} from '../../models/device.interface';
import {DeviceContract, UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';

@Component(
  {
    selector: 'app-device-viewer',
    templateUrl: 'device-viewer.component.html',
    styleUrls: ['device-viewer.component.scss']
  }
)
export class DeviceViewerComponent implements OnInit {

  devices: Device[];
  deviceContractInstance: any;

  constructor(@Inject(WEB3) private web3: Web3, @Inject(DeviceContract) private deviceContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.devices = [];
    try {
      this.deviceContract.setProvider(this.web3.currentProvider);
      this.deviceContractInstance = await this.deviceContract.deployed();
      const count: number = await this.deviceContractInstance.count.call();
      for (let counter = 0; counter < count; counter++) {
        const device: Device = await this.deviceContractInstance.findByIndex.call(counter);
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
}
