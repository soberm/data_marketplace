import {Component, Inject, OnInit} from '@angular/core';

import {UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {User} from '../../models/user.interface';
import {ActivatedRoute, Params, Router} from '@angular/router';

@Component(
  {
    selector: 'app-user-update-form',
    templateUrl: './user-update-form.component.html',
    styleUrls: ['user-update-form.component.scss']
  }
)
export class UserUpdateFormComponent implements OnInit {

  details: User;

  accounts: string[];
  userContractInstance: any;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    @Inject(WEB3) private web3: Web3,
    @Inject(UserContract) private userContract: any) {
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

    this.route.params.subscribe(async (data: Params) => {
      try {
        this.details = await this.userContractInstance.findByAddress.call(data.address);
      } catch (e) {
        if (e instanceof Error) {
          console.log(e.message);
        } else {
          console.log(e);
        }
      }
    });
  }

  async handleSubmit(user: User, isValid: boolean) {
    if (isValid) {
      await this.userContractInstance
        .update(user.firstName, user.lastName, user.company, user.email, {from: this.accounts[0]});
      await this.router.navigateByUrl('users/management');
    }
  }

}
