import {Inject, Injectable} from '@angular/core';
import {WalletClient} from '../injection';
import {AuthService} from '../auth/auth.service';
import {WalletServiceClient} from '../proto/proxy/Wallet_serviceServiceClientPb';
import {CreateWalletRequest, CreateWalletResponse} from '../proto/proxy/wallet_service_pb';

@Injectable()
export class WalletService {

  constructor(@Inject(WalletClient) private walletClient: WalletServiceClient, private authService: AuthService) {
  }

  createWallet(request: CreateWalletRequest): Promise<CreateWalletResponse.AsObject> {
    const metadata = {authorization: 'bearer ' + this.authService.getToken()};
    return new Promise((resolve, reject) => {
      this.walletClient.createWallet(request, metadata, (err, response: CreateWalletResponse) => {
        if (err) {
          return reject(err);
        }
        resolve(response.toObject());
      });
    });
  }

}

