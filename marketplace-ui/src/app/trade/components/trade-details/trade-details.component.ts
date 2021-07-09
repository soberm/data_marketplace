import {Component, Input} from '@angular/core';

import {Trade} from '../../models/trade.interface';

@Component(
  {
    selector: 'app-trade-details',
    templateUrl: './trade-details.component.html',
    styleUrls: ['trade-details.component.scss']
  }
)
export class TradeDetailsComponent {

  @Input()
  detail: Trade;

  constructor() {
  }

}
