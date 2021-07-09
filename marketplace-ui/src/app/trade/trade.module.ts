import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {TradeDetailsComponent} from './components/trade-details/trade-details.component';
import {MatCardModule} from '@angular/material/card';
import {TradeViewerComponent} from './containers/trade-viewer/trade-viewer.component';
import {RouterModule, Routes} from '@angular/router';
import {ProductViewerComponent} from '../product/containers/product-viewer/product-viewer.component';
import {FlexLayoutModule} from '@angular/flex-layout';
import {MatBadgeModule} from '@angular/material/badge';
import {MatDividerModule} from '@angular/material/divider';

const routes: Routes = [
  {
    path: 'trades',
    children: [
      {path: '', component: TradeViewerComponent},
    ]
  }
];

@NgModule({
  declarations: [
    TradeDetailsComponent,
    TradeViewerComponent
  ],
  imports: [
    CommonModule,
    MatCardModule,
    RouterModule.forChild(routes),
    FlexLayoutModule,
    MatBadgeModule,
    MatDividerModule
  ]
})
export class TradeModule {
}
