package lr

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/Bo0km4n/claude/app/tablet/pkg/db"
	"github.com/Bo0km4n/claude/app/tablet/pkg/model"
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

func TestFetchLRByDistance(t *testing.T) {
	repo := NewLRRepository(db.Mysql)
	insertTestLocation(t, repo)

	queries := []struct {
		Longitude float32
		Latitude  float32
		Distance  float32
	}{
		{
			Longitude: 139.801278,
			Latitude:  35.652547,
			Distance:  10,
		},
	}

	for _, q := range queries {
		rows, err := repo.FetchLRsByDistance(context.Background(), q.Latitude, q.Longitude, q.Distance)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows.Entries) != 3 {
			t.Errorf("len(rows.Entries) is expected %d, got=%d", 3, len(rows.Entries))
		}
	}
}

func insertTestLocation(t *testing.T, repo LRRepository) {
	currentTime := time.Now()
	rows := []*model.LREntry{
		&model.LREntry{
			Latitude: 35.663729, Longitude: 139.744047, // Huric KamiyaChou Billding: Tokyo
			CreatedAt: currentTime, UpdatedAt: currentTime,
		},
		&model.LREntry{
			Latitude: 35.666863, Longitude: 139.74954, // Toranomonn Hills: Tokyo
			CreatedAt: currentTime, UpdatedAt: currentTime,
		},
		&model.LREntry{
			Latitude: 35.660477, Longitude: 139.729356, // Roppongi Hills: Tokyo
			CreatedAt: currentTime, UpdatedAt: currentTime,
		},
		&model.LREntry{
			Latitude: 35.689604, Longitude: 139.692305, // Tokyo Tochou: Tokyo
			CreatedAt: currentTime, UpdatedAt: currentTime,
		},
		&model.LREntry{
			Latitude: 43.064313, Longitude: 141.347255, // Hokkaidou Chou: Hokkaidou
			CreatedAt: currentTime, UpdatedAt: currentTime,
		},
		&model.LREntry{
			Latitude: 37.532225, Longitude: -122.313028, // US
			CreatedAt: currentTime, UpdatedAt: currentTime,
		},
	}

	for _, r := range rows {
		repo.StoreLR(context.Background(), r.SerializeToProto())
	}
}
