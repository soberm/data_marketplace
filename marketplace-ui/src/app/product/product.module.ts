import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ProductViewerComponent} from './containers/product-viewer/product-viewer.component';
import {ProductDetailsComponent} from './components/product-details/product-details.component';
import {MatCardModule} from '@angular/material/card';
import {RouterModule, Routes} from '@angular/router';
import {FlexLayoutModule} from '@angular/flex-layout';
import {ProductCreateFormComponent} from './containers/product-create-form/product-create-form.component';
import {MatFormFieldModule} from '@angular/material/form-field';
import {FormsModule} from '@angular/forms';
import {MatInputModule} from '@angular/material/input';
import {MatIconModule} from '@angular/material/icon';
import {MatButtonModule} from '@angular/material/button';
import {ProductManagementComponent} from './containers/product-management/product-management.component';
import {MatTooltipModule} from '@angular/material/tooltip';
import {ProductActionDetailsComponent} from './components/product-action-details/product-action-details.component';
import {ProductUpdateFormComponent} from './containers/product-update-form/product-update-form.component';
import {DeviceUpdateFormComponent} from '../device/containers/device-update-form/device-update-form.component';
import {MatDividerModule} from '@angular/material/divider';

const routes: Routes = [
  {
    path: 'products',
    children: [
      {path: 'viewer', component: ProductViewerComponent},
      {path: 'management', component: ProductManagementComponent},
      {path: 'management/create', component: ProductCreateFormComponent},
      {path: 'management/update/:id', component: ProductUpdateFormComponent},
    ]
  }
];

@NgModule({
  declarations: [
    ProductCreateFormComponent,
    ProductUpdateFormComponent,
    ProductViewerComponent,
    ProductDetailsComponent,
    ProductManagementComponent,
    ProductActionDetailsComponent
  ],
  imports: [
    CommonModule,
    MatCardModule,
    RouterModule.forChild(routes),
    FlexLayoutModule,
    MatFormFieldModule,
    FormsModule,
    MatInputModule,
    MatIconModule,
    MatButtonModule,
    MatTooltipModule,
    MatDividerModule,
  ]
})
export class ProductModule {
}
