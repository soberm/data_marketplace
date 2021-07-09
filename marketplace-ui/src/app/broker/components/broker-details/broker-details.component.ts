import {Component, Input} from '@angular/core';

import {Broker} from '../../models/broker.interface';

@Component(
  {
    selector: 'app-broker-details',
    templateUrl: './broker-details.component.html',
    styleUrls: ['broker-details.component.scss']
  }
)
export class BrokerDetailsComponent {

  @Input()
  detail: Broker;

  locations: string[] = [
    'Brazil', 'Europe Nordic & East', 'Europe West', 'Latin America North', 'Latin America South', 'North America',
    'Oceania', 'Russia', 'Turkey', 'Japan', 'The Philippines', 'Singapore, Malaysia, and Indonesia', 'Taiwan, Hong Kong, and Macao',
    'Vietnam', 'Thailand', 'Republic of Korea', 'China',
  ];

  constructor() {
  }

}
