package contracts

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"marketplace-services/pkg/contracts/bindings"
	"math/big"
)

type Product struct {
	Id          *big.Int
	Device      common.Address
	Name        string
	Description string
	DataType    string
	Frequency   *big.Int
	Cost        *big.Int
	Deleted     bool
}

func ProductStructToProduct(result struct {
	Id          *big.Int
	Device      common.Address
	Name        string
	Description string
	DataType    string
	Frequency   *big.Int
	Cost        *big.Int
	Deleted     bool
}) *Product {
	return &Product{
		Id:          result.Id,
		Device:      result.Device,
		Name:        result.Name,
		Description: result.Description,
		DataType:    result.DataType,
		Frequency:   result.Frequency,
		Cost:        result.Cost,
		Deleted:     result.Deleted,
	}
}

type ProductContract interface {
	CreateProduct(opts *bind.TransactOpts, product *Product) (*types.Transaction, error)
	UpdateProduct(opts *bind.TransactOpts, product *Product) (*types.Transaction, error)
	RemoveProduct(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
	FindProductByIndex(opts *bind.CallOpts, index *big.Int) (*Product, error)
	FindProductById(opts *bind.CallOpts, id *big.Int) (*Product, error)
	FindTradingRequestsOfProductById(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error)
	FindNegotiationRequestsOfProductById(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error)
	FindTradesOfProductById(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error)
	FindNegotiationsOfProductById(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error)
	FindCostOfProductById(opts *bind.CallOpts, id *big.Int) (*big.Int, error)
	IsProductOwnedByDevice(opts *bind.CallOpts, id *big.Int, device common.Address) (bool, error)
	CountProducts(opts *bind.CallOpts) (*big.Int, error)
	ExistsProductById(opts *bind.CallOpts, id *big.Int) (bool, error)
	ExistsProductByIdAndDeleted(opts *bind.CallOpts, id *big.Int, deleted bool) (bool, error)
	WatchCreatedProductEvent(opts *bind.WatchOpts, sink chan<- *bindings.ProductContractCreatedProduct, user []common.Address) (event.Subscription, error)
	WatchUpdatedProductEvent(opts *bind.WatchOpts, sink chan<- *bindings.ProductContractUpdatedProduct, id []*big.Int, user []common.Address) (event.Subscription, error)
	WatchRemovedProductEvent(opts *bind.WatchOpts, sink chan<- *bindings.ProductContractRemovedProduct, id []*big.Int, user []common.Address) (event.Subscription, error)
}

type productContractImpl struct {
	binding *bindings.ProductContract
}

func NewProductContractImpl(address common.Address, backend bind.ContractBackend) (*productContractImpl, error) {
	binding, err := bindings.NewProductContract(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new product contract binding %s: %w", address.Hex(), err)
	}
	return &productContractImpl{binding: binding}, err
}

func (p *productContractImpl) CreateProduct(opts *bind.TransactOpts, product *Product) (*types.Transaction, error) {
	return p.binding.Create(opts, product.Device, product.Name, product.Description, product.DataType, product.Frequency, product.Cost)
}

func (p *productContractImpl) UpdateProduct(opts *bind.TransactOpts, product *Product) (*types.Transaction, error) {
	return p.binding.Update(opts, product.Id, product.Name, product.Description, product.DataType, product.Frequency, product.Cost)
}

func (p *productContractImpl) RemoveProduct(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return p.binding.Remove(opts, id)
}

func (p *productContractImpl) FindProductByIndex(opts *bind.CallOpts, index *big.Int) (*Product, error) {
	result, err := p.binding.FindByIndex(opts, index)
	return ProductStructToProduct(result), err
}

func (p *productContractImpl) FindProductById(opts *bind.CallOpts, id *big.Int) (*Product, error) {
	result, err := p.binding.FindById(opts, id)
	return ProductStructToProduct(result), err
}

func (p *productContractImpl) FindTradingRequestsOfProductById(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error) {
	return p.binding.FindTradingRequestsById(opts, id)
}

func (p *productContractImpl) FindNegotiationRequestsOfProductById(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error) {
	return p.binding.FindNegotiationRequestsById(opts, id)
}

func (p *productContractImpl) FindTradesOfProductById(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error) {
	return p.binding.FindTradesById(opts, id)
}

func (p *productContractImpl) FindNegotiationsOfProductById(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error) {
	return p.binding.FindNegotiationsById(opts, id)
}

func (p *productContractImpl) FindCostOfProductById(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	return p.binding.FindCostById(opts, id)
}

func (p *productContractImpl) IsProductOwnedByDevice(opts *bind.CallOpts, id *big.Int, device common.Address) (bool, error) {
	return p.binding.IsOwnedByDevice(opts, id, device)
}

func (p *productContractImpl) CountProducts(opts *bind.CallOpts) (*big.Int, error) {
	return p.binding.Count(opts)
}

func (p *productContractImpl) ExistsProductById(opts *bind.CallOpts, id *big.Int) (bool, error) {
	return p.binding.ExistsById(opts, id)
}

func (p *productContractImpl) ExistsProductByIdAndDeleted(opts *bind.CallOpts, id *big.Int, deleted bool) (bool, error) {
	return p.binding.ExistsByIdAndDeleted(opts, id, deleted)
}

func (p *productContractImpl) WatchCreatedProductEvent(opts *bind.WatchOpts, sink chan<- *bindings.ProductContractCreatedProduct, user []common.Address) (event.Subscription, error) {
	return p.binding.WatchCreatedProduct(opts, sink, user)
}

func (p *productContractImpl) WatchUpdatedProductEvent(opts *bind.WatchOpts, sink chan<- *bindings.ProductContractUpdatedProduct, id []*big.Int, user []common.Address) (event.Subscription, error) {
	return p.binding.WatchUpdatedProduct(opts, sink, id, user)
}

func (p *productContractImpl) WatchRemovedProductEvent(opts *bind.WatchOpts, sink chan<- *bindings.ProductContractRemovedProduct, id []*big.Int, user []common.Address) (event.Subscription, error) {
	return p.binding.WatchRemovedProduct(opts, sink, id, user)
}
