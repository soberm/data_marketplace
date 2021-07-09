import {Component, EventEmitter, Input, Output} from '@angular/core';

import {Broker} from '../../models/broker.interface';

@Component(
  {
    selector: 'app-broker-action-details',
    templateUrl: './broker-action-details.component.html',
    styleUrls: ['broker-action-details.component.scss']
  }
)
export class BrokerActionDetailsComponent {

  @Input()
  detail: Broker;

  @Output()
  delete: EventEmitter<Broker> = new EventEmitter<Broker>();

  locations: string[] = [
    'Brazil', 'Europe Nordic & East', 'Europe West', 'Latin America North', 'Latin America South', 'North America',
    'Oceania', 'Russia', 'Turkey', 'Japan', 'The Philippines', 'Singapore, Malaysia, and Indonesia', 'Taiwan, Hong Kong, and Macao',
    'Vietnam', 'Thailand', 'Republic of Korea', 'China',
  ];

  constructor() {
  }

  handleDelete(broker: Broker) {
    this.delete.emit(broker);
  }

}
