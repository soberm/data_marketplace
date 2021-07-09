package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"marketplace-services/pkg/proxy/api"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:25566", grpc.WithInsecure())
	if err != nil {
		logger.Errorf("%+v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf("%+v", err)
		}
	}()

	authServiceClient := api.NewAuthServiceClient(conn)
	walletServiceClient := api.NewWalletServiceClient(conn)

	request := &api.GetTokenRequest{Username: "kristina", Password: []byte("12345678")}
	response, err := authServiceClient.GetToken(context.Background(), request)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	md := metadata.Pairs("authorization", "bearer "+response.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	w, err := walletServiceClient.FindWalletById(ctx, &api.FindWalletByIdRequest{Id: 2})
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	logger.Infof("PublicKey: %s", hexutil.Encode(w.Wallet.PublicKey))
	logger.Infof("%s", hexutil.Encode(w.Wallet.Address))
	_, err = crypto.UnmarshalPubkey(w.Wallet.PublicKey)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
}
