import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FlexLayoutModule} from '@angular/flex-layout';
import {MatCardModule} from '@angular/material/card';
import {RouterModule, Routes} from '@angular/router';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {MatIconModule} from '@angular/material/icon';
import {MatButtonModule} from '@angular/material/button';
import {FormsModule} from '@angular/forms';
import {MatTooltipModule} from '@angular/material/tooltip';
import {MatDividerModule} from '@angular/material/divider';
import {LoginComponent} from './components/login/login.component';
import {AuthService} from './auth.service';
import {AuthGuard} from "./auth.guard";

const routes: Routes = [
  {
    path: 'auth',
    children: [
      {path: 'login', component: LoginComponent},
    ]
  }
];

@NgModule({
  declarations: [
    LoginComponent
  ],
  exports: [],
  imports: [
    CommonModule,
    FlexLayoutModule,
    MatCardModule,
    RouterModule.forChild(routes),
    MatFormFieldModule,
    MatInputModule,
    MatIconModule,
    MatButtonModule,
    FormsModule,
    MatTooltipModule,
    MatDividerModule
  ],
  providers: [AuthService, AuthGuard],
})
export class AuthModule {
}
