package impl

import (
	"context"

	"github.com/kade-chen/google-billing-console/apps/user"
	"github.com/kade-chen/library/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *service) delete(ctx context.Context, set *user.UserSet) error {

	// 专门优化单个删除
	var result *mongo.DeleteResult
	var err error
	if len(set.Items) == 1 {
		result, err = s.col.DeleteMany(ctx, bson.M{"_id": set.Items[0].Meta.Id})
	} else {
		result, err = s.col.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": set.UserIds()}})
	}

	if err != nil {
		return exception.NewInternalServerError("This user is empty")
	}

	if result.DeletedCount == 0 {
		return exception.NewNotFound("user %s not found", set)
	}

	return nil
}
