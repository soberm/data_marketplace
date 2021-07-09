import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {UserViewerComponent} from './containers/user-viewer/user-viewer.component';
import {UserDetailsComponent} from './components/user-details/user-details.component';
import {FlexLayoutModule} from '@angular/flex-layout';
import {MatCardModule} from '@angular/material/card';
import {RouterModule, Routes} from '@angular/router';
import {UserCreateFormComponent} from './containers/user-create-form/user-create-form.component';
import {UserManagementComponent} from './containers/user-management/user-management.component';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {MatIconModule} from '@angular/material/icon';
import {MatButtonModule} from '@angular/material/button';
import {FormsModule} from '@angular/forms';
import {MatTooltipModule} from '@angular/material/tooltip';
import {UserActionDetailsComponent} from './components/user-action-details/user-action-details.component';
import {DeviceUpdateFormComponent} from '../device/containers/device-update-form/device-update-form.component';
import {UserUpdateFormComponent} from './containers/user-update-form/user-update-form.component';
import {MatDividerModule} from '@angular/material/divider';

const routes: Routes = [
  {
    path: 'users',
    children: [
      {path: 'viewer', component: UserViewerComponent},
      {path: 'management', component: UserManagementComponent},
      {path: 'management/create', component: UserCreateFormComponent},
      {path: 'management/update/:address', component: UserUpdateFormComponent},
    ]
  }
];

@NgModule({
  declarations: [
    UserManagementComponent,
    UserCreateFormComponent,
    UserUpdateFormComponent,
    UserViewerComponent,
    UserDetailsComponent,
    UserActionDetailsComponent
  ],
  exports: [
    UserViewerComponent,
    UserManagementComponent
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
  ]
})
export class UserModule {
}
