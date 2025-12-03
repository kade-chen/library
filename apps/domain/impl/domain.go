package impl

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/domain"
)

func (s *service) CreateDomain(ctx context.Context, req *domain.CreateDomainRequest) (*domain.Domain, error) {
	d, err := domain.NewDomain(req)

	// if err != nil {
	// 	return nil, exception.NewBadRequest("%s", err.Error())
	// }
	// if _, err := s.col.InsertOne(ctx, d); err != nil {
	// 	return nil, exception.NewInternalServerError("inserted a domain document error, %s", err)
	// }
	s.log.Info().Msgf("create a domain document success")
	return d, err
}

func (s *service) DescribeDomain(ctx context.Context, req *domain.DescribeDomainRequest) (*domain.Domain, error) {

	// if err := req.Validate(); err != nil {
	// 	return nil, exception.NewBadRequest("%s", err.Error())
	// }

	// filter := bson.M{}
	// switch req.DescribeBy {
	// case domain.DESCRIBE_BY_ID:
	// 	filter["_id"] = req.Id
	// case domain.DESCRIBE_BY_NAME:
	// 	filter["name"] = req.Name
	// }

	// d := domain.NewDefaultDomain()
	// if err := s.col.FindOne(ctx, filter).Decode(d); err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		s.log.Error().Msgf("domain %s not found", req)
	// 		return nil, exception.NewNotFound("domain %s not found", req)
	// 	}

	// 	return nil, exception.NewInternalServerError("find domain %s error, %s", req.Id, err)
	// }

	return nil, nil
}
