import {Component, Input} from '@angular/core';

import * as domain_account_pb from '../../../proto/domain/account_pb';

@Component(
  {
    selector: 'app-account-details',
    templateUrl: './account-details.component.html',
    styleUrls: ['account-details.component.scss']
  }
)
export class AccountDetailsComponent {

  @Input()
  detail: domain_account_pb.Account.AsObject;

  constructor() {
  }

}
