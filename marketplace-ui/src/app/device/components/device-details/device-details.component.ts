import {Component, Input} from '@angular/core';

import {Device} from '../../models/device.interface';

@Component(
  {
    selector: 'app-device-details',
    templateUrl: './device-details.component.html',
    styleUrls: ['device-details.component.scss']
  }
)
export class DeviceDetailsComponent {

  @Input()
  detail: Device;

  constructor() {
  }

}
