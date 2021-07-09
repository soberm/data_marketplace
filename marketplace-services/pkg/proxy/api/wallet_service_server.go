package api

import (
	"context"
	"marketplace-services/pkg/proxy/services"
)

type walletServiceServer struct {
	UnimplementedWalletServiceServer
	walletService services.WalletService
}

func NewWalletServiceServer(service services.WalletService) *walletServiceServer {
	return &walletServiceServer{walletService: service}
}

func (s *walletServiceServer) CreateWallet(
	ctx context.Context,
	req *CreateWalletRequest,
) (*CreateWalletResponse, error) {
	w, err := s.walletService.CreateWallet(ctx, WalletFromGrpcWallet(req.Wallet))
	if err != nil {
		return nil, err
	}
	return &CreateWalletResponse{Wallet: WalletToGrpcWallet(w)}, err
}

func (s *walletServiceServer) FindWalletById(
	ctx context.Context,
	req *FindWalletByIdRequest,
) (*FindWalletByIdResponse, error) {
	w, err := s.walletService.FindWalletById(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &FindWalletByIdResponse{Wallet: WalletToGrpcWallet(w)}, err
}

func (s *walletServiceServer) FindWalletByAccountId(
	ctx context.Context,
	req *FindWalletByAccountIdRequest,
) (*FindWalletByAccountIdResponse, error) {
	w, err := s.walletService.FindWalletByAccountId(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &FindWalletByAccountIdResponse{Wallet: WalletToGrpcWallet(w)}, err
}
