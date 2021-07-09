import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DeviceDetailsComponent} from './components/device-details/device-details.component';
import {MatCardModule} from '@angular/material/card';
import {MatButtonModule} from '@angular/material/button';
import {DeviceViewerComponent} from './containers/device-viewer/device-viewer.component';
import {FlexLayoutModule} from '@angular/flex-layout';
import {RouterModule, Routes} from '@angular/router';
import {DeviceManagementComponent} from './containers/device-management/device-management.component';
import {DeviceCreateFormComponent} from './containers/device-create-form/device-create-form.component';
import {FormsModule} from '@angular/forms';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';
import {DeviceActionDetailsComponent} from './components/device-action-details/device-action-details.component';
import {MatTooltipModule} from '@angular/material/tooltip';
import {DeviceUpdateFormComponent} from './containers/device-update-form/device-update-form.component';
import {MatDividerModule} from '@angular/material/divider';

const routes: Routes = [
  {
    path: 'devices',
    children: [
      {path: 'viewer', component: DeviceViewerComponent},
      {path: 'management', component: DeviceManagementComponent},
      {path: 'management/create', component: DeviceCreateFormComponent},
      {path: 'management/update/:address', component: DeviceUpdateFormComponent},
    ]
  }
];

@NgModule({
  declarations: [
    DeviceViewerComponent,
    DeviceDetailsComponent,
    DeviceActionDetailsComponent,
    DeviceManagementComponent,
    DeviceCreateFormComponent,
    DeviceUpdateFormComponent
  ],
  exports: [
    DeviceViewerComponent,
    DeviceManagementComponent
  ],
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    FlexLayoutModule,
    RouterModule.forChild(routes),
    FormsModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatTooltipModule,
    MatDividerModule
  ]
})
export class DeviceModule {
}
