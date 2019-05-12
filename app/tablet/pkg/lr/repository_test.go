package lr

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/tablet/pkg/db"
)

func TestMain(m *testing.M) {
	db.InitMysql("claude_test")
	db.MigrateMysql()
	code := m.Run()
	db.CloseMysql()
	os.Exit(code)
}

func TestStoreLR(t *testing.T) {
	repo := NewLRRepository(db.Mysql)
	if _, err := repo.StoreLR(context.Background(), &proto.LREntry{
		GlobalIp:   "100.10.10.10",
		GlobalPort: "7000",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}); err != nil {
		t.Fatal(err)
	}
}

func TestLoadLR(t *testing.T) {
	repo := NewLRRepository(db.Mysql)
	if _, err := repo.StoreLR(context.Background(), &proto.LREntry{
		GlobalIp:   "200.10.10.10",
		GlobalPort: "7000",
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}); err != nil {
		t.Fatal(err)
	}

	result, err := repo.LoadLRs(context.Background(), &proto.LREntry{
		GlobalIp: "200.10.10.10",
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Entries) != 1 {
		t.Errorf("result.Entries length is expected 1. got=%d", len(result.Entries))
	}
}
