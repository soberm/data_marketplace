import {Inject, Injectable} from '@angular/core';
import {AccountClient} from '../injection';
import {AccountServiceClient} from '../proto/proxy/Account_serviceServiceClientPb';
import {
  CreateAccountRequest,
  CreateAccountResponse,
  FindAccountsRequest,
  FindAccountsResponse
} from '../proto/proxy/account_service_pb';
import {AuthService} from '../auth/auth.service';
import {Observable} from 'rxjs';
import {Status} from 'grpc-web';
import {Error} from 'grpc-web';

@Injectable()
export class AccountService {

  constructor(@Inject(AccountClient) private accountClient: AccountServiceClient, private authService: AuthService) {
  }

  createAccount(request: CreateAccountRequest): Promise<CreateAccountResponse.AsObject> {
    const metadata = {authorization: 'bearer ' + this.authService.getToken()};
    return new Promise((resolve, reject) => {
      this.accountClient.createAccount(request, metadata, (err, response: CreateAccountResponse) => {
        if (err) {
          return reject(err);
        }
        resolve(response.toObject());
      });
    });
  }

  findAccounts(): Observable<FindAccountsResponse.AsObject> {
    const metadata = {authorization: 'bearer ' + this.authService.getToken()};
    return new Observable<FindAccountsResponse.AsObject>(obs => {
      const stream = this.accountClient.findAccounts(new FindAccountsRequest(), metadata);
      stream.on('error', (err: Error) => {
        obs.error(err);
      });
      stream.on('status', (status: Status) => {
        console.log(status);
      });
      stream.on('data', (response: any) => {
        obs.next(response.toObject());
      });
      stream.on('end', () => {
        obs.complete();
      });
    });
  }

}

