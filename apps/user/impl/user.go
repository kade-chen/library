package impl

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/user"
)

func (s *service) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	// 1.create user strust
	u1, _ := user.NewUser(req)
	// schema, err := bigquery.InferSchema(u1)
	// if err != nil {
	// 	log.Fatal("infer schema failed:", err)
	// }
	// 2.insert to db
	// _, err := s.col.InsertOne(ctx, u1)
	// if err != nil {
	// 	return nil, exception.NewInternalServerError("user %s error, %s", req, err.Error())
	// }

	return u1, nil
}

func (s *service) DescribeUser(ctx context.Context, req *user.DescribeUserRequest) (*user.User, error) {
	// filter := bson.M{}

	// switch req.DescribeBy {
	// case user.DESCRIBE_BY_USER_ID:
	// 	filter["_id"] = req.Id
	// case user.DESCRIBE_BY_USER_NAME:
	// 	filter["username"] = req.Username
	// case user.DESCRIBE_BY_FEISHU_USER_ID:
	// 	filter["feishu.user_id"] = req.Id
	// default:
	// 	return nil, exception.NewBadRequest("unknow desribe by %s", req.DescribeBy)
	// }

	// u, _ := user.NewUser(nil)
	// s.log.Debug().Msgf("query filter: %s", filter)

	// // 2.find user
	// err := s.col.FindOne(ctx, filter).Decode(u)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		return nil, exception.NewNotFound("user %s No found in the document", req.Username)
	// 	}
	// 	return nil, exception.NewInternalServerError("user %s error, %s", req, err.Error())
	// }
	// // res, _ := bson.MarshalExtJSON(u, false, false)
	// // fmt.Println(string(res))

	// //3.initalize null value
	// u.InitializeNullValue()

	return nil, nil
}

// 查询用户列表
func (s *service) ListUser(ctx context.Context, req *user.QueryUserRequest) (*user.UserSet, error) {
	// r := user.NewQueryUserRequest(req)
	// resp, err := s.col.Find(ctx, r.FindFilter(), r.FindOptions().SetSort(bson.D{{"create_at", 1}}))
	// if err != nil {
	// 	s.log.Error().Msgf("find user error, error is %s", err)
	// 	return nil, exception.NewInternalServerError("find user error, error is %s", err)
	// }
	// set := user.NewUserSet()
	// // 循环
	// if !req.SkipItems {
	// 	for resp.Next(ctx) {
	// 		ins := user.NewDefaultUser()
	// 		if err := resp.Decode(ins); err != nil {
	// 			s.log.Error().Msgf("decode user error, error is %s", err)
	// 			return nil, exception.NewInternalServerError("decode user error, error is %s", err)
	// 		}
	// 		ins.Desensitization()
	// 		s.log.Info().Msgf("this %s passwd has been desensitized", ins.Spec.Username)
	// 		set.Add(ins)
	// 	}
	// }
	// // count
	// count, err := s.col.CountDocuments(ctx, r.FindFilter())
	// if err != nil {
	// 	s.log.Error().Msgf("get user count error, error is %s", err)
	// 	return nil, exception.NewInternalServerError("get user count error, error is %s", err)
	// }
	// set.Total = count
	return nil, nil
}

// 查询用户列表
func (s *service) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.UserSet, error) {
	// // 0.判断这些要删除的用户是否存在
	// queryReq := user.NewQueryUserDeleteRequest()
	// queryReq.UserIds = req.UserIds
	// s.log.Info().Msgf("this is user ids to delete: %s", req.UserIds)
	// listuser, err := s.ListUser(ctx, queryReq)
	// if err != nil {
	// 	s.log.Error().Msgf("find user error, error is %s", err)
	// 	return nil, exception.NewInternalServerError("find user error, error is %s", err)
	// }
	// //1.delete users from the user_ids
	// if err := s.delete(ctx, listuser); err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
