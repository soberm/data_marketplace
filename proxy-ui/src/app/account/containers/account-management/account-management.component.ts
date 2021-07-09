import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {AccountService} from '../../account.service';
import * as domain_account_pb from '../../../proto/domain/account_pb';
import {FindAccountsResponse} from '../../../proto/proxy/account_service_pb';
import {Role} from '../../../proto/domain/role_pb';

@Component(
  {
    selector: 'app-account-management',
    templateUrl: './account-management.component.html',
    styleUrls: ['account-management.component.scss']
  }
)
export class AccountManagementComponent implements OnInit {

  accounts: domain_account_pb.Account.AsObject[];

  constructor(private router: Router, private accountService: AccountService) {
  }

  ngOnInit(): void {
    this.accounts = [];
    this.accountService.findAccounts().subscribe((response: FindAccountsResponse.AsObject) => {
        console.log(response.account);
        if (response.account.role === Role.USER) {
          this.accounts = [...this.accounts, response.account];
        }
      },
      (error) => {
        console.log(error);
      },
      () => {
        console.log(this.accounts);
      }
    );
  }

}
