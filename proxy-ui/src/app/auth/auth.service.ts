import {Inject, Injectable} from '@angular/core';
import {AuthClient, JwtHelper} from '../injection';
import {AuthServiceClient} from '../proto/proxy/Auth_serviceServiceClientPb';
import {JwtHelperService} from '@auth0/angular-jwt';
import {GetTokenRequest, GetTokenResponse} from '../proto/proxy/auth_service_pb';
import {DecodedAccessToken} from './model/decoded-access-token';

@Injectable()
export class AuthService {

  constructor(
    @Inject(AuthClient) private authClient: AuthServiceClient,
    @Inject(JwtHelper) private jwtHelper: JwtHelperService) {
  }

  authenticate(request: GetTokenRequest): Promise<GetTokenResponse.AsObject> {
    return new Promise((resolve, reject) => {
      this.authClient.getToken(request, null, (err, response: GetTokenResponse) => {
        if (err) {
          return reject(err);
        }
        resolve(response.toObject());
      });
    });
  }

  isAuthenticated(): boolean {
    const accessToken: string = localStorage.getItem('access_token');

    if (accessToken == null) {
      return false;
    }

    return !this.jwtHelper.isTokenExpired(accessToken);
  }

  getPrincipal(): DecodedAccessToken {
    if (!this.isAuthenticated()) {
      return null;
    }

    const accessToken: string = localStorage.getItem('access_token');
    return this.jwtHelper.decodeToken(accessToken);
  }

  getToken(): string {
    return localStorage.getItem('access_token');
  }

  hasRole(role: number): boolean {
    if (!this.isAuthenticated()) {
      return false;
    }

    const accessToken: string = localStorage.getItem('access_token');
    const decodedToken: DecodedAccessToken = this.jwtHelper.decodeToken(accessToken);

    return decodedToken.role === role;
  }

  setAccessToken(accessToken: string): void {
    localStorage.setItem('access_token', accessToken);
  }

  clearAccessToken(): void {
    localStorage.removeItem('access_token');
  }

}

