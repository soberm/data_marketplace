import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {BrokerDetailsComponent} from './components/broker-details/broker-details.component';
import {MatCardModule} from '@angular/material/card';
import {BrokerViewerComponent} from './containers/broker-viewer/broker-viewer.component';
import {FlexLayoutModule} from '@angular/flex-layout';
import {RouterModule, Routes} from '@angular/router';
import {BrokerActionDetailsComponent} from './components/broker-action-details/broker-action-details.component';
import {MatIconModule} from '@angular/material/icon';
import {MatTooltipModule} from '@angular/material/tooltip';
import {MatButtonModule} from '@angular/material/button';
import {BrokerManagementComponent} from './containers/broker-management/broker-management.component';
import {BrokerCreateFormComponent} from './containers/broker-create-form/broker-create-form.component';
import {FormsModule} from '@angular/forms';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {MatSelectModule} from '@angular/material/select';
import {BrokerUpdateFormComponent} from './containers/broker-update-form/broker-update-form.component';
import {MatDividerModule} from '@angular/material/divider';

const routes: Routes = [
  {
    path: 'brokers',
    children: [
      {path: 'viewer', component: BrokerViewerComponent},
      {path: 'management', component: BrokerManagementComponent},
      {path: 'management/create', component: BrokerCreateFormComponent},
      {path: 'management/update/:address', component: BrokerUpdateFormComponent},
    ]
  }
];

@NgModule({
  declarations: [
    BrokerDetailsComponent,
    BrokerCreateFormComponent,
    BrokerUpdateFormComponent,
    BrokerManagementComponent,
    BrokerActionDetailsComponent,
    BrokerViewerComponent
  ],
  imports: [
    CommonModule,
    MatCardModule,
    FlexLayoutModule,
    RouterModule.forChild(routes),
    MatIconModule,
    MatTooltipModule,
    MatButtonModule,
    FormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatDividerModule
  ]
})
export class BrokerModule {
}
