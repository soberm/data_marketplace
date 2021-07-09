import {Component, Inject, OnInit} from '@angular/core';

import {User} from '../../models/user.interface';
import {UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {Router} from '@angular/router';

@Component(
  {
    selector: 'app-user-create-form',
    templateUrl: './user-create-form.component.html',
    styleUrls: ['user-create-form.component.scss']
  }
)
export class UserCreateFormComponent implements OnInit {

  accounts: string[];
  userContractInstance: any;

  constructor(private router: Router, @Inject(WEB3) private web3: Web3, @Inject(UserContract) private userContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.userContract.setProvider(this.web3.currentProvider);
    try {
      this.accounts = await this.web3.eth.getAccounts();
      this.userContractInstance = await this.userContract.deployed();
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

  async handleSubmit(user: User, isValid: boolean) {
    if (isValid) {
      await this.userContractInstance
        .create(user.firstName, user.lastName, user.company, user.email, {from: this.accounts[0]});
      await this.router.navigateByUrl('users/management');
    }
  }

}
