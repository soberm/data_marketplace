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
import {AccountManagementComponent} from './containers/account-management/account-management.component';
import {AuthGuard} from '../auth/auth.guard';
import {AccountCreateFormComponent} from './containers/account-create-form/account-create-form.component';
import {AccountService} from './account.service';
import {WalletService} from "./wallet.service";
import {AccountDetailsComponent} from "./components/account-details/account-details.component";
import {HexPipe} from "./hex.pipe";

const routes: Routes = [
  {
    path: 'accounts',
    children: [
      {path: 'management', component: AccountManagementComponent, canActivate: [AuthGuard]},
      {path: 'management/create', component: AccountCreateFormComponent, canActivate: [AuthGuard]},
    ]
  }
];

@NgModule({
  declarations: [
    AccountManagementComponent,
    AccountCreateFormComponent,
    AccountDetailsComponent,
    HexPipe
  ],
  exports: [
    AccountManagementComponent
  ],
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
  providers: [
    AccountService,
    WalletService
  ],
})
export class AccountModule {
}
