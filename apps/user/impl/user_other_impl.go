package impl

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/user"
	"github.com/kade-chen/library/exception"
	"google.golang.org/api/googleapi"
)

func (s *service) delete(ctx context.Context, set *user.UserSet) error {
	items := set.Items

	if len(items) == 0 {
		s.log.Error().Msg("empty delete set")
		return exception.NewBadRequest("empty delete set")
	}

	// BigQuery DELETE
	deleteSQL := fmt.Sprintf(`DELETE FROM %s WHERE id IN UNNEST(@ids)`, s.bqTableFull)

	q := s.bq_client.Query(deleteSQL)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "ids",
			Value: set.UserIds(), // []string → ARRAY<STRING>
		},
	}

	job, err := q.Run(ctx)
	if err != nil {
		s.log.Error().Msgf("delete user error: %v", err)
		return exception.NewInternalServerError("delete user error: %v", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		s.log.Error().Msgf("delete user job error: %v", err)
		return exception.NewInternalServerError("delete user job error: %v", err)
	}
	if status.Err() != nil {
		s.log.Error().Msgf("delete user bq error: %v", status.Err())
		return exception.NewInternalServerError("delete user bq error: %v", status.Err())
	}
	return nil
}

func (s *service) datasetisNotFound(err error) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		return gerr.Code == 404
	}
	return false
}
