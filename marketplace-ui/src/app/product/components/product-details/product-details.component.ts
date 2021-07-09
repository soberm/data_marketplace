import {Component, Input} from '@angular/core';

import {Product} from '../../models/product.interface';

@Component(
  {
    selector: 'app-product-details',
    templateUrl: './product-details.component.html',
    styleUrls: ['product-details.component.scss']
  }
)
export class ProductDetailsComponent {

  @Input()
  detail: Product;

  constructor() {
  }

}
