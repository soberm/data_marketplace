import {Component, Inject, OnInit} from '@angular/core';
import {User} from '../../models/user.interface';
import {UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';

@Component(
  {
    selector: 'app-user-viewer',
    templateUrl: 'user-viewer.component.html',
    styleUrls: ['user-viewer.component.scss']
  }
)
export class UserViewerComponent implements OnInit {

  users: User[];
  userContractInstance: any;

  constructor(@Inject(WEB3) private web3: Web3, @Inject(UserContract) private userContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.users = [];
    try {
      this.userContract.setProvider(this.web3.currentProvider);
      this.userContractInstance = await this.userContract.deployed();
      const count: number = await this.userContractInstance.count.call();
      for (let counter = 0; counter < count; counter++) {
        const user: User = await this.userContractInstance.findByIndex.call(counter);
        if (user.deleted) {
          continue;
        }
        this.users = [...this.users, user];
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
