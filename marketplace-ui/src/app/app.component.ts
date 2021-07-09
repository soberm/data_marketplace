import {Component, Inject, OnInit} from '@angular/core';
import {WEB3} from './web3';
import Web3 from 'web3';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'marketplace-ui';

  constructor(@Inject(WEB3) private web3: any) {
  }

  async ngOnInit(): Promise<void> {
    if ('enable' in this.web3.currentProvider) {
      await this.web3.currentProvider.enable();
    }
  }

}
