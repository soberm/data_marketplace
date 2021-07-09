import {Component, Inject, OnInit} from '@angular/core';

import {Router} from '@angular/router';
import {CreateAccountRequest, CreateAccountResponse} from '../../../proto/proxy/account_service_pb';
import * as domain_account_pb from '../../../proto/domain/account_pb';
import {Role} from '../../../proto/domain/role_pb';
import {AccountService} from '../../account.service';
import {WalletService} from '../../wallet.service';
import {CreateWalletRequest, CreateWalletResponse} from '../../../proto/proxy/wallet_service_pb';
import * as domain_wallet_pb from '../../../proto/domain/wallet_pb';

@Component(
  {
    selector: 'app-account-create-form',
    templateUrl: './account-create-form.component.html',
    styleUrls: ['account-create-form.component.scss']
  }
)
export class AccountCreateFormComponent implements OnInit {

  error: string;

  constructor(private router: Router, private accountService: AccountService, private walletService: WalletService) {
  }

  async ngOnInit(): Promise<void> {
  }

  async handleSubmit(form: any, isValid: boolean) {
    if (isValid) {
      try {
        const account: domain_account_pb.Account = new domain_account_pb.Account();
        account.setName(form.username);
        account.setPassword(new TextEncoder().encode(form.password));
        account.setRole(Role.USER);

        const createAccountRequest = new CreateAccountRequest();
        createAccountRequest.setAccount(account);

        const createAccountResponse: CreateAccountResponse.AsObject = await this.accountService.createAccount(createAccountRequest);
        const wallet: domain_wallet_pb.Wallet = new domain_wallet_pb.Wallet();
        wallet.setUserid(createAccountResponse.account.id);
        wallet.setPassphrase(form.passphrase);

        const createWalletRequest = new CreateWalletRequest();
        createWalletRequest.setWallet(wallet);

        await this.walletService.createWallet(createWalletRequest);
        await this.router.navigateByUrl('accounts/management');
      } catch (e) {
        this.error = e.message;
      }

    }
  }

}
