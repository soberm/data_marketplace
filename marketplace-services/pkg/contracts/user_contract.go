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

type User struct {
	Addr      common.Address
	FirstName string
	LastName  string
	Company   string
	Email     string
	Deleted   bool
	Devices   []common.Address
	Brokers   []common.Address
}

func UserStructToUser(result struct {
	Addr      common.Address
	FirstName string
	LastName  string
	Company   string
	Email     string
	Deleted   bool
	Devices   []common.Address
	Brokers   []common.Address
}) *User {
	return &User{
		Addr:      result.Addr,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Company:   result.Company,
		Email:     result.Email,
		Deleted:   result.Deleted,
		Devices:   result.Devices,
		Brokers:   result.Brokers,
	}
}

type UserContract interface {
	CreateUser(opts *bind.TransactOpts, user *User) (*types.Transaction, error)
	UpdateUser(opts *bind.TransactOpts, user *User) (*types.Transaction, error)
	RemoveUser(opts *bind.TransactOpts) (*types.Transaction, error)
	FindUserByIndex(opts *bind.CallOpts, index *big.Int) (*User, error)
	FindUserByAddress(opts *bind.CallOpts, addr common.Address) (*User, error)
	CountUsers(opts *bind.CallOpts) (*big.Int, error)
	ExistsUserByAddress(opts *bind.CallOpts, addr common.Address) (bool, error)
	ExistsUserByAddressAndDeleted(opts *bind.CallOpts, addr common.Address, deleted bool) (bool, error)
	WatchCreatedUserEvent(opts *bind.WatchOpts, sink chan<- *bindings.UserContractCreatedUser, addr []common.Address) (event.Subscription, error)
	WatchUpdatedUserEvent(opts *bind.WatchOpts, sink chan<- *bindings.UserContractUpdatedUser, addr []common.Address) (event.Subscription, error)
	WatchRemovedUserEvent(opts *bind.WatchOpts, sink chan<- *bindings.UserContractRemovedUser, addr []common.Address) (event.Subscription, error)
}

type userContractImpl struct {
	binding *bindings.UserContract
}

func NewUserContractImpl(address common.Address, backend bind.ContractBackend) (*userContractImpl, error) {
	binding, err := bindings.NewUserContract(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new user contract binding %s: %w", address.Hex(), err)
	}
	return &userContractImpl{binding: binding}, err
}

func (u *userContractImpl) CreateUser(opts *bind.TransactOpts, user *User) (*types.Transaction, error) {
	return u.binding.Create(opts, user.FirstName, user.LastName, user.Company, user.Email)
}

func (u *userContractImpl) UpdateUser(opts *bind.TransactOpts, user *User) (*types.Transaction, error) {
	return u.binding.Update(opts, user.FirstName, user.LastName, user.Company, user.Email)
}

func (u *userContractImpl) RemoveUser(opts *bind.TransactOpts) (*types.Transaction, error) {
	return u.binding.Remove(opts)
}

func (u *userContractImpl) FindUserByIndex(opts *bind.CallOpts, index *big.Int) (*User, error) {
	result, err := u.binding.FindByIndex(opts, index)
	return UserStructToUser(result), err
}

func (u *userContractImpl) FindUserByAddress(opts *bind.CallOpts, addr common.Address) (*User, error) {
	result, err := u.binding.FindByAddress(opts, addr)
	return UserStructToUser(result), err
}

func (u *userContractImpl) CountUsers(opts *bind.CallOpts) (*big.Int, error) {
	return u.binding.Count(opts)
}

func (u *userContractImpl) ExistsUserByAddress(opts *bind.CallOpts, addr common.Address) (bool, error) {
	return u.binding.ExistsByAddress(opts, addr)
}

func (u *userContractImpl) ExistsUserByAddressAndDeleted(opts *bind.CallOpts, addr common.Address, deleted bool) (bool, error) {
	return u.binding.ExistsByAddressAndDeleted(opts, addr, deleted)
}

func (u *userContractImpl) WatchCreatedUserEvent(opts *bind.WatchOpts, sink chan<- *bindings.UserContractCreatedUser, addr []common.Address) (event.Subscription, error) {
	return u.binding.WatchCreatedUser(opts, sink, addr)
}

func (u *userContractImpl) WatchUpdatedUserEvent(opts *bind.WatchOpts, sink chan<- *bindings.UserContractUpdatedUser, addr []common.Address) (event.Subscription, error) {
	return u.binding.WatchUpdatedUser(opts, sink, addr)
}

func (u *userContractImpl) WatchRemovedUserEvent(opts *bind.WatchOpts, sink chan<- *bindings.UserContractRemovedUser, addr []common.Address) (event.Subscription, error) {
	return u.binding.WatchRemovedUser(opts, sink, addr)
}
