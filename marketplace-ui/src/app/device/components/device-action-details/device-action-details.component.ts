import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';

import {Device} from '../../models/device.interface';
import {User} from '../../../user/models/user.interface';

@Component(
  {
    selector: 'app-device-action-details',
    templateUrl: './device-action-details.component.html',
    styleUrls: ['device-action-details.component.scss']
  }
)
export class DeviceActionDetailsComponent implements OnInit {

  @Input()
  detail: Device;

  @Output()
  delete: EventEmitter<Device> = new EventEmitter<Device>();

  constructor() {
  }

  handleDelete(device: Device) {
    this.delete.emit(device);
  }

  ngOnInit(): void {
  }

}
