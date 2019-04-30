package lr

import (
	"context"
	"testing"
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
)

func TestStoreLR(t *testing.T) {
	repo := NewLRRepository("mongodb://root:example@localhost:27017")
	now := time.Now()
	if err := repo.StoreLR(context.Background(), &proto.LREntry{
		GlobalIp:   "192.168.10.10",
		GlobalPort: "8080",
		Latitude:   40.28030,
		Longtitude: 135.19,
		CreatedAt:  int64(now.Unix()),
		UpdatedAt:  int64(now.Unix()),
	}); err != nil {
		t.Fatal(err)
	}
}
