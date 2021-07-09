package api

import (
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"marketplace-services/pkg/broker/api"
	"marketplace-services/pkg/domain"
	"marketplace-services/pkg/proxy/services"
)

type discoveryServiceServer struct {
	UnimplementedDiscoveryServiceServer
	discoveryService services.DiscoveryService
}

func NewDiscoveryServiceServer(discoveryService services.DiscoveryService) *discoveryServiceServer {
	return &discoveryServiceServer{discoveryService: discoveryService}
}

func (s *discoveryServiceServer) SearchBroker(
	req *SearchBrokerRequest,
	stream DiscoveryService_SearchBrokerServer,
) error {
	ctx := stream.Context()
	sink, errc := s.discoveryService.SearchBroker(ctx, BrokerSearchQueryFromGrpcBrokerSearchQuery(req.Query))
	for b := range sink {
		broker := BrokerToGrpcBroker(b)
		if err := stream.Send(&SearchBrokerResponse{
			Broker: broker,
		}); err != nil {
			return err
		}
	}
	if err := <-errc; err != nil {
		return err
	}
	return nil
}

func (s *discoveryServiceServer) SearchProductWithBroker(
	req *SearchProductWithBrokerRequest,
	stream DiscoveryService_SearchProductWithBrokerServer,
) error {
	conn, err := grpc.Dial(req.BrokerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf("%+v", err)
		}
	}()

	discoveryServiceClient := api.NewDiscoveryServiceClient(conn)
	productStream, err := discoveryServiceClient.SearchProduct(stream.Context(), &api.SearchProductRequest{
		Query: &domain.ProductSearchQuery{
			DataType: req.Query.DataType,
			MinCost:  req.Query.MinCost,
			MaxCost:  req.Query.MaxCost,
		},
	})
	if err != nil {
		return err
	}

	for {
		response, err := productStream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		err = stream.Send(&SearchProductWithBrokerResponse{
			Product: response.Product,
		})
		if err != nil {
			return err
		}
	}
}
