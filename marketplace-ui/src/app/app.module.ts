import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatMenuModule} from '@angular/material/menu';
import {MatIconModule} from '@angular/material/icon';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatButtonModule} from '@angular/material/button';
import {FlexLayoutModule} from '@angular/flex-layout';
import {MatCardModule} from '@angular/material/card';
import {UserModule} from './user/user.module';
import {DeviceModule} from './device/device.module';
import {ProductModule} from './product/product.module';
import {TradeModule} from './trade/trade.module';
import {BrokerModule} from './broker/broker.module';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    UserModule,
    ProductModule,
    TradeModule,
    BrokerModule,
    FlexLayoutModule,
    BrowserAnimationsModule,
    MatMenuModule,
    MatIconModule,
    MatToolbarModule,
    MatButtonModule,
    MatCardModule,
    DeviceModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
