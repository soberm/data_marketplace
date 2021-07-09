package main

import (
	"context"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"marketplace-services/pkg/domain"
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
	userContractServiceClient := api.NewUserContractServiceClient(conn)
	deviceContractServiceClient := api.NewDeviceContractServiceClient(conn)
	productContractServiceClient := api.NewProductContractServiceClient(conn)
	brokerContractServiceClient := api.NewBrokerContractServiceClient(conn)

	getTokenResponse, err := authServiceClient.GetToken(context.Background(), &api.GetTokenRequest{
		Username: "kristina",
		Password: []byte("12345678"),
	})
	if err != nil {
		logger.Errorf("%+v", err)
		return
	}

	md := metadata.Pairs("authorization", "bearer "+getTokenResponse.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = userContractServiceClient.CreateUser(ctx, &api.CreateUserRequest{
		User: &domain.User{
			FirstName: "michael",
			LastName:  "sober",
			Company:   "gmbh",
			Email:     "michael.sober@ymail.com",
		},
	})
	if err != nil {
		logger.Errorf("%+v", err)
		return
	}

	w1, err := walletServiceClient.FindWalletById(ctx, &api.FindWalletByIdRequest{Id: 1})
	if err != nil {
		logger.Errorf("%+v", err)
		return
	}

	w2, err := walletServiceClient.FindWalletById(ctx, &api.FindWalletByIdRequest{Id: 2})
	if err != nil {
		logger.Errorf("%+v", err)
		return
	}

	devices := []*domain.Device{
		{
			Address:     "0x6da49C19d815c1c61050046456398599720716A2",
			Name:        "Raspberry Pi A",
			Description: "Small single-board computer with a temperature sensor",
			PublicKey:   w1.Wallet.PublicKey,
		},
		{
			Address:     "0x65bF93E7AC35260fd1dB8eF97E306F01fde96DD7",
			Name:        "Raspberry Pi B",
			Description: "Small single-board computer with a light sensor",
			PublicKey:   w2.Wallet.PublicKey,
		},
	}

	for _, device := range devices {
		_, err := deviceContractServiceClient.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: device,
		})
		if err != nil {
			logger.Errorf("%+v", err)
			return
		}
	}

	products := []*domain.Product{
		{
			Device:      "0x6da49C19d815c1c61050046456398599720716A2",
			Name:        "High frequency temperature data with good quality",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			DataType:    "temperature",
			Frequency:   1,
			Cost:        10,
		},
		{
			Device:      "0x65bF93E7AC35260fd1dB8eF97E306F01fde96DD7",
			Name:        "Low frequency light data",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			DataType:    "light",
			Frequency:   10,
			Cost:        20,
		},
	}

	for _, product := range products {
		_, err := productContractServiceClient.CreateProduct(ctx, &api.CreateProductRequest{
			Product: product,
		})
		if err != nil {
			logger.Errorf("%+v", err)
			return
		}
	}

	brokers := []*domain.Broker{
		{
			Address:  "0x9278Fcc1b8a086E52FB6253d1922FD9235869300",
			Name:     "Broker1",
			HostAddr: "127.0.0.1:25565",
			Location: domain.Location_OCE,
		},
		{
			Address:  "0xA2e654c259b78C6a0d7D1d7fCB84421E83e98564",
			Name:     "Broker2",
			HostAddr: "127.0.0.1:25564",
			Location: domain.Location_EUW,
		},
	}

	for _, broker := range brokers {
		_, err := brokerContractServiceClient.CreateBroker(ctx, &api.CreateBrokerRequest{
			Broker: broker,
		})
		if err != nil {
			logger.Errorf("%+v", err)
			return
		}
	}

}
