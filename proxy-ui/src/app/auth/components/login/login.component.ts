import {Component} from '@angular/core';
import {Router} from '@angular/router';
import {AuthService} from '../../auth.service';
import {GetTokenRequest, GetTokenResponse} from '../../../proto/proxy/auth_service_pb';

@Component(
  {
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['login.component.scss']
  }
)
export class LoginComponent {

  constructor(private router: Router, private authService: AuthService) {
    console.log(authService);
  }

  async handleSubmit(form: any, isValid: boolean) {
    if (isValid) {
      const request = new GetTokenRequest();
      request.setUsername(form.username);
      request.setPassword(new TextEncoder().encode(form.password));
      const response: GetTokenResponse.AsObject = await this.authService.authenticate(request);
      this.authService.setAccessToken(response.token);
      console.log(this.authService.getPrincipal());
      console.log(this.authService.isAuthenticated());
      console.log(this.authService.hasRole(1));
      await this.router.navigateByUrl('accounts/management');
    }
  }

}
