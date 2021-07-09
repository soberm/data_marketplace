import {Component, Inject, OnInit} from '@angular/core';
import {Product} from '../../models/product.interface';
import {ProductContract, WEB3} from '../../../web3';
import Web3 from 'web3';

@Component(
  {
    selector: 'app-product-viewer',
    templateUrl: 'product-viewer.component.html',
    styleUrls: ['product-viewer.component.scss']
  }
)
export class ProductViewerComponent implements OnInit {

  products: Product[];
  productContractInstance: any;

  constructor(@Inject(WEB3) private web3: Web3, @Inject(ProductContract) private productContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.products = [];
    try {
      this.productContract.setProvider(this.web3.currentProvider);
      this.productContractInstance = await this.productContract.deployed();
      const count: number = await this.productContractInstance.count.call();
      for (let counter = 0; counter < count; counter++) {
        const product: Product = await this.productContractInstance.findByIndex.call(counter);
        if (product.deleted) {
          continue;
        }
        this.products = [...this.products, product];
      }
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

}
