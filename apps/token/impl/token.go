package impl

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/library/exception"
)

// This function is used to issue a token based on the given request and return the token.
func (s *service) IssueToken(ctx context.Context, req *token.IssueTokenRequest) (*token.Token, error) {
	// 1.issuer token
	tk, err := s.issuer_token(ctx, req)
	if err != nil {
		return nil, err
	}
	// 2.login security after loading

	return tk, nil
}

func (s *service) ValicateToken(ctx context.Context, req *token.ValicateTokenRequest) (*token.Token, error) {
	// 1.query wether token exist

	tk, err := s.get(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	// // 2.verify wether the token is expired
	// time, err := tk.CheckIssue_Token_expried(tk.IssueAt, tk.AccessExpiredAt, tk.AccessToken)
	// if err != nil {
	// 	return nil, err
	// }
	// // 3.print info log
	// s.log.Info().Msg(time)
	return tk, err
}

func (s *service) UpdateToken(ctx context.Context, req *token.UpdateTokenRequest) (string, error) {
	// 用 fmt.Sprintf 拼 SQL
	sql := fmt.Sprintf(`
		MERGE %v T
		USING (
			SELECT access_token
			FROM %v
			WHERE refresh_token = @refresh_token
			LIMIT 1
		) S
		ON T.access_token = S.access_token
		WHEN MATCHED THEN
			UPDATE SET access_token = @new_access_token;
	`, s.bqTableFull, s.bqTableFull)

	q := s.bq_client.Query(sql)
	q.UseLegacySQL = false
	// 4️⃣ 绑定参数
	q.Parameters = []bigquery.QueryParameter{
		{Name: "new_access_token", Value: req.AccessToken},
		{Name: "refresh_token", Value: req.RefreshToken},
	}
	job, err := q.Run(ctx)
	if err != nil {
		return "", exception.NewInternalServerError("failed to run query: %v", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return "", exception.NewInternalServerError("job wait failed: %v", err)
	}

	if err := status.Err(); err != nil {
		return "", exception.NewInternalServerError("job execution error: %v", err)
	}
	// ✅ 检查更新行数
	// 1️⃣ 获取 Statistics.Details
	if stats, ok := status.Statistics.Details.(*bigquery.QueryStatistics); ok {
		if stats.DMLStats != nil {
			updatedRows := stats.DMLStats.UpdatedRowCount
			if updatedRows == 0 {
				return "", exception.NewInternalServerError("no row updated")
			}
			fmt.Println("Updated rows:", updatedRows)
		}
	}

	fmt.Println("Update successful!")
	return req.AccessToken, nil
}

// remove Token
func (s *service) RevolkToken(ctx context.Context, req *token.RevolkTokenRequest) (*token.Token, int64, error) {
	//1.query wether the token exist
	tk, err := s.get(ctx, req.AccessToken)
	if err != nil {
		return nil, 0, err
	}

	//2.delete token
	ins, row, err := s.delete(ctx, tk)
	if err != nil {
		return nil, 0, err
	}
	//3.delete ui cookie
	s.log.Info().Msgf("revolk token success. Number of rows affecred %d", row)
	return ins, row, nil
}
