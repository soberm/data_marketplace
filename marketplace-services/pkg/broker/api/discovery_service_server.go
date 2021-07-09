package api

import (
	"marketplace-services/pkg/broker/services"
)

type discoveryServiceServer struct {
	UnimplementedDiscoveryServiceServer
	discoveryService services.DiscoveryService
}

func NewDiscoveryServiceServer(discoveryService services.DiscoveryService) *discoveryServiceServer {
	return &discoveryServiceServer{discoveryService: discoveryService}
}

func (s *discoveryServiceServer) SearchProduct(
	req *SearchProductRequest,
	stream DiscoveryService_SearchProductServer,
) error {
	ctx := stream.Context()
	sink, errc := s.discoveryService.SearchProduct(ctx, ProductSearchQueryFromGrpcProductSearchQuery(req.Query))
	for p := range sink {
		product := ProductToGrpcProduct(p)
		if err := stream.Send(&SearchProductResponse{
			Product: product,
		}); err != nil {
			return err
		}
	}
	if err := <-errc; err != nil {
		return err
	}
	return nil
}
