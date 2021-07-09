import {Component, Inject, OnInit} from '@angular/core';
import {User} from '../../models/user.interface';
import {UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {Product} from '../../../product/models/product.interface';

@Component(
  {
    selector: 'app-user-management',
    templateUrl: 'user-management.component.html',
    styleUrls: ['user-management.component.scss']
  }
)
export class UserManagementComponent implements OnInit {

  user: User;
  accounts: string[];
  userContractInstance: any;

  constructor(@Inject(WEB3) private web3: Web3, @Inject(UserContract) private userContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.accounts = await this.web3.eth.getAccounts();
    this.userContract.setProvider(this.web3.currentProvider);

    try {
      this.userContractInstance = await this.userContract.deployed();
      const exists = this.userContractInstance.existsByAddress.call(this.accounts[0]);
      if (exists) {
        this.user = await this.userContractInstance.findByAddress.call(this.accounts[0]);
      }
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

  async onDeleteUser(event: User) {
    try {
      await this.userContractInstance.remove(event.addr, {from: this.accounts[0]});
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
      return;
    }
    this.user = undefined;
  }

}
