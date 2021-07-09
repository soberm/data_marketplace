import {InjectionToken} from '@angular/core';
import {AuthServiceClient} from './proto/proxy/Auth_serviceServiceClientPb';
import {JwtHelperService} from '@auth0/angular-jwt';
import {AccountServiceClient} from './proto/proxy/Account_serviceServiceClientPb';
import {WalletServiceClient} from './proto/proxy/Wallet_serviceServiceClientPb';
import {environment} from '../environments/environment';

export const AuthClient = new InjectionToken<AuthServiceClient>('authClient', {
  providedIn: 'root',
  factory: () => new AuthServiceClient(environment.host),
});

export const AccountClient = new InjectionToken<AccountServiceClient>('accountClient', {
  providedIn: 'root',
  factory: () => new AccountServiceClient(environment.host),
});

export const WalletClient = new InjectionToken<WalletServiceClient>('walletClient', {
  providedIn: 'root',
  factory: () => new WalletServiceClient(environment.host),
});

export const JwtHelper = new InjectionToken<JwtHelperService>('jwtHelper', {
  providedIn: 'root',
  factory: () => new JwtHelperService(),
});
