package tests

import (
	"strings"
	"testing"

	"github.com/iamkanishka/yapily-client-go/domain"
	"github.com/iamkanishka/yapily-client-go/utils"
)

func TestPageIteratorCollect(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e"}
	callCount := 0

	iter := utils.NewPageIterator[string](2, func(p domain.PaginationParams) ([]string, error) {
		callCount++
		start := p.Offset
		end := start + p.Limit
		if start >= len(data) {
			return nil, nil
		}
		if end > len(data) {
			end = len(data)
		}
		return data[start:end], nil
	})

	all, err := iter.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 5 {
		t.Errorf("expected 5 items, got %d", len(all))
	}
	// Should have taken 3 calls: [0:2], [2:4], [4:5]
	if callCount != 3 {
		t.Errorf("expected 3 page fetches, got %d", callCount)
	}
}

func TestPageIteratorNext(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	iter := utils.NewPageIterator[int](3, func(p domain.PaginationParams) ([]int, error) {
		start := p.Offset
		end := start + p.Limit
		if start >= len(data) {
			return nil, nil
		}
		if end > len(data) {
			end = len(data)
		}
		return data[start:end], nil
	})

	page1, ok1 := iter.Next()
	if !ok1 || len(page1) != 3 {
		t.Errorf("page1: expected 3 items and ok=true, got %d ok=%v", len(page1), ok1)
	}

	page2, ok2 := iter.Next()
	if !ok2 || len(page2) != 3 {
		t.Errorf("page2: expected 3 items and ok=true, got %d ok=%v", len(page2), ok2)
	}

	page3, ok3 := iter.Next()
	if ok3 {
		t.Errorf("page3: expected ok=false (exhausted), got items=%v", page3)
	}
}

func TestPageIteratorEmpty(t *testing.T) {
	iter := utils.NewPageIterator[string](10, func(_ domain.PaginationParams) ([]string, error) {
		return nil, nil
	})

	all, err := iter.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("expected 0 items, got %d", len(all))
	}
}

func TestIdempotencyKeyPassthrough(t *testing.T) {
	key := utils.IdempotencyKey("my-key")
	if key != "my-key" {
		t.Errorf("expected my-key, got %s", key)
	}
}

func TestIdempotencyKeyGenerated(t *testing.T) {
	key := utils.IdempotencyKey("")
	if key == "" {
		t.Error("expected non-empty generated key")
	}
	if !strings.HasPrefix(key, "idem-") {
		t.Errorf("expected key to start with 'idem-', got %s", key)
	}
}

func TestIdempotencyKeyUnique(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		k := utils.IdempotencyKey("")
		if seen[k] {
			t.Errorf("duplicate idempotency key generated: %s", k)
		}
		seen[k] = true
	}
}
