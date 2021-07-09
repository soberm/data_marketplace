import {Component} from '@angular/core';
import {AuthService} from './auth/auth.service';
import {Router} from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'proxy-ui';

  constructor(private router: Router, private authService: AuthService) {
  }

  onLogout() {
    this.authService.clearAccessToken();
    this.router.navigateByUrl('/auth/login');
  }

  isAuthenticated(): boolean {
    return this.authService.isAuthenticated();
  }

}

