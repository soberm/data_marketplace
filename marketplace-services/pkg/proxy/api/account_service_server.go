package api

import (
	"context"
	"marketplace-services/pkg/proxy/services"
)

type accountServiceServer struct {
	UnimplementedAccountServiceServer
	accountService services.AccountService
}

func NewAccountServiceServer(service services.AccountService) *accountServiceServer {
	return &accountServiceServer{accountService: service}
}

func (s *accountServiceServer) CreateAccount(
	ctx context.Context,
	req *CreateAccountRequest,
) (*CreateAccountResponse, error) {
	account, err := s.accountService.CreateAccount(ctx, AccountFromGrpcAccount(req.Account))
	if err != nil {
		return nil, err
	}
	return &CreateAccountResponse{Account: AccountToGrpcAccount(account)}, err
}

func (s *accountServiceServer) FindAccountById(
	ctx context.Context,
	req *FindAccountByIdRequest,
) (*FindAccountByIdResponse, error) {
	account, err := s.accountService.FindAccountById(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &FindAccountByIdResponse{Account: AccountToGrpcAccount(account)}, err
}

func (s *accountServiceServer) FindAccountByName(
	ctx context.Context,
	req *FindAccountByNameRequest,
) (*FindAccountByNameResponse, error) {
	account, err := s.accountService.FindAccountByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &FindAccountByNameResponse{Account: AccountToGrpcAccount(account)}, err
}

func (s *accountServiceServer) FindAccounts(req *FindAccountsRequest, stream AccountService_FindAccountsServer) error {
	accounts, err := s.accountService.FindAccounts(stream.Context())
	if err != nil {
		return err
	}
	for _, account := range accounts {
		err = stream.Send(&FindAccountsResponse{Account: AccountToGrpcAccount(account)})
		if err != nil {
			return err
		}
	}
	return nil
}
