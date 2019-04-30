package lr

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/tablet/db"
)

type LRRepository interface {
	StoreLR(ctx context.Context, lr *proto.LREntry) error
}

type lrRepository struct {
	db *mongo.Client
}

func NewLRRepository(host string) LRRepository {
	cli, _ := db.NewMongo(host)
	return &lrRepository{
		db: cli,
	}
}

func (lrr *lrRepository) StoreLR(ctx context.Context, lr *proto.LREntry) error {
	collection := lrr.db.Database("claude_tablet").Collection("lr")
	_, err := collection.InsertOne(ctx, lr)
	return err
}
