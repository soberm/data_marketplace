import {Component, Inject, OnInit} from '@angular/core';
import {Product} from '../../models/product.interface';
import {DeviceContract, ProductContract, UserContract, WEB3} from '../../../web3';
import Web3 from 'web3';
import {User} from '../../../user/models/user.interface';

@Component(
  {
    selector: 'app-product-management',
    templateUrl: 'product-management.component.html',
    styleUrls: ['product-management.component.scss']
  }
)
export class ProductManagementComponent implements OnInit {

  products: Product[];
  accounts: string[];
  userContractInstance: any;
  deviceContractInstance: any;
  productContractInstance: any;

  constructor(
    @Inject(WEB3) private web3: Web3,
    @Inject(UserContract) private userContract: any,
    @Inject(DeviceContract) private deviceContract: any,
    @Inject(ProductContract) private productContract: any) {
  }

  async ngOnInit(): Promise<void> {
    this.products = [];
    this.accounts = await this.web3.eth.getAccounts();

    this.userContract.setProvider(this.web3.currentProvider);
    this.deviceContract.setProvider(this.web3.currentProvider);
    this.productContract.setProvider(this.web3.currentProvider);

    try {

      this.userContractInstance = await this.userContract.deployed();
      this.deviceContractInstance = await this.deviceContract.deployed();
      this.productContractInstance = await this.productContract.deployed();

      const user: User = await this.userContractInstance.findByAddress.call(this.accounts[0]);

      for (const address of user.devices) {
        const products = await this.deviceContractInstance.findProductsByAddress.call(address);
        for (const productId of products) {
          const product = await this.productContractInstance.findById.call(productId.toNumber());
          if (product.deleted) {
            continue;
          }
          this.products = [...this.products, product];
        }
      }
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
    }
  }

  async onDeleteProduct(event: Product) {
    try {
      await this.productContractInstance.remove(event.id, {from: this.accounts[0]});
    } catch (e) {
      if (e instanceof Error) {
        console.log(e.message);
      } else {
        console.log(e);
      }
      return;
    }
    this.products = this.products.filter((product: Product) => {
      return product.id !== event.id;
    });
  }

}
