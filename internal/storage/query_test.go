package storage_test

import (
	"testing"
	"time"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

func populatedStore(t *testing.T) *storage.Store {
	t.Helper()
	s, err := storage.New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	now := time.Now()
	entries := []*logentry.Entry{
		{Job: "backup", Message: "ok", Timestamp: now.Add(-3 * time.Hour), IsError: false},
		{Job: "backup", Message: "fail", Timestamp: now.Add(-2 * time.Hour), IsError: true},
		{Job: "sync", Message: "ok", Timestamp: now.Add(-1 * time.Hour), IsError: false},
		{Job: "sync", Message: "fail", Timestamp: now, IsError: true},
	}
	for _, e := range entries {
		_ = s.Add(e)
	}
	return s
}

func TestQueryByJob(t *testing.T) {
	s := populatedStore(t)
	res := s.Query(storage.Filter{Job: "backup"})
	if len(res) != 2 {
		t.Fatalf("expected 2, got %d", len(res))
	}
}

func TestQueryOnlyErrors(t *testing.T) {
	s := populatedStore(t)
	res := s.Query(storage.Filter{OnlyErr: true})
	if len(res) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(res))
	}
}

func TestQuerySince(t *testing.T) {
	s := populatedStore(t)
	res := s.Query(storage.Filter{Since: time.Now().Add(-90 * time.Minute)})
	if len(res) != 1 {
		t.Fatalf("expected 1 entry since 90m ago, got %d", len(res))
	}
}

func TestQueryMaxCount(t *testing.T) {
	s := populatedStore(t)
	res := s.Query(storage.Filter{MaxCount: 2})
	if len(res) != 2 {
		t.Fatalf("expected 2, got %d", len(res))
	}
}

func TestQueryCombined(t *testing.T) {
	s := populatedStore(t)
	res := s.Query(storage.Filter{Job: "sync", OnlyErr: true})
	if len(res) != 1 {
		t.Fatalf("expected 1, got %d", len(res))
	}
	if res[0].Message != "fail" {
		t.Fatalf("unexpected message: %s", res[0].Message)
	}
}
