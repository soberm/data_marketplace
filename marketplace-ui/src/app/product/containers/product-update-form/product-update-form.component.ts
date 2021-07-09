import {Component, Inject, OnInit} from '@angular/core';

import {Product} from '../../models/product.interface';
import {ActivatedRoute, Params, Router} from '@angular/router';
import {ProductContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {Device} from '../../../device/models/device.interface';

@Component(
  {
    selector: 'app-product-update-form',
    templateUrl: './product-update-form.component.html',
    styleUrls: ['product-update-form.component.scss']
  }
)
export class ProductUpdateFormComponent implements OnInit {

  details: Product;

  accounts: string[];
  productContractInstance: any;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    @Inject(WEB3) private web3: Web3,
    @Inject(ProductContract) private productContract: any) {
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

    this.route.params.subscribe(async (data: Params) => {
      try {
        this.details = await this.productContractInstance.findById.call(data.id);
      } catch (e) {
        if (e instanceof Error) {
          console.log(e.message);
        } else {
          console.log(e);
        }
      }
    });
  }

  async handleSubmit(product: Product, isValid: boolean) {
    if (isValid) {
      await this.productContractInstance.update(
        this.details.id,
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
