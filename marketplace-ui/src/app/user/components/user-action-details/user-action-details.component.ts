import {Component, EventEmitter, Input, Output} from '@angular/core';

import {User} from '../../models/user.interface';

@Component(
  {
    selector: 'app-user-action-details',
    templateUrl: './user-action-details.component.html',
    styleUrls: ['user-action-details.component.scss']
  }
)
export class UserActionDetailsComponent {

  @Input()
  detail: User;

  @Output()
  delete: EventEmitter<User> = new EventEmitter<User>();

  constructor() {
  }

  handleDelete(user: User) {
    this.delete.emit(user);
  }

}
