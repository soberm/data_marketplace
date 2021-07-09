import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {AuthModule} from './auth/auth.module';
import {LoginComponent} from './auth/components/login/login.component';


const routes: Routes = [
  {path: '', redirectTo: 'accounts/management', pathMatch: 'full'},
  {path: '**', component: LoginComponent}
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes),
    AuthModule],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
