import {Component, Inject, OnInit} from '@angular/core';

import {Product} from '../../models/product.interface';
import {Router} from '@angular/router';
import {ProductContract, WEB3} from '../../../web3';
import Web3 from 'web3';

@Component(
  {
    selector: 'app-product-create-form',
    templateUrl: './product-create-form.component.html',
    styleUrls: ['product-create-form.component.scss']
  }
)
export class ProductCreateFormComponent implements OnInit {

  accounts: string[];
  productContractInstance: any;

  constructor(private router: Router, @Inject(WEB3) private web3: Web3, @Inject(ProductContract) private productContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.productContract.setProvider(this.web3.currentProvider);
    try {
      this.accounts = await this.web3.eth.getAccounts();
      this.productContractInstance = await this.productContract.deployed();
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

  async handleSubmit(product: Product, isValid: boolean) {
    if (isValid) {
      await this.productContractInstance.create(
        product.device,
        product.name,
        product.description,
        product.dataType,
        product.frequency,
        product.cost,
        {from: this.accounts[0]});
      await this.router.navigateByUrl('products/management');
    }
  }

}
