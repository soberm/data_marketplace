import {Component, Input} from '@angular/core';

import {User} from '../../models/user.interface';

@Component(
  {
    selector: 'app-user-details',
    templateUrl: './user-details.component.html',
    styleUrls: ['user-details.component.scss']
  }
)
export class UserDetailsComponent {

  @Input()
  detail: User;

  constructor() {
  }

}
