import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Product} from '../../models/product.interface';

@Component(
  {
    selector: 'app-product-action-details',
    templateUrl: './product-action-details.component.html',
    styleUrls: ['product-action-details.component.scss']
  }
)
export class ProductActionDetailsComponent implements OnInit {

  @Input()
  detail: Product;

  @Output()
  delete: EventEmitter<Product> = new EventEmitter<Product>();

  constructor() {
  }

  handleDelete(product: Product) {
    this.delete.emit(product);
  }

  ngOnInit(): void {
  }

}
