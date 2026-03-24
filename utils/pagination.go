// Package utils provides utility helpers for the Yapily SDK.
package utils

import (
	"fmt"
	"time"

	"github.com/iamkanishka/yapily-client-go/domain"
)

// PageIterator is a generic helper for iterating paginated results.
type PageIterator[T any] struct {
	pageSize int
	offset   int
	done     bool
	fetchFn  func(params domain.PaginationParams) ([]T, error)
}

// NewPageIterator creates a new PageIterator.
// fetchFn receives pagination parameters and returns the next page of results.
func NewPageIterator[T any](pageSize int, fetchFn func(params domain.PaginationParams) ([]T, error)) *PageIterator[T] {
	return &PageIterator[T]{
		pageSize: pageSize,
		fetchFn:  fetchFn,
	}
}

// Next fetches the next page. Returns nil, false when there are no more pages.
func (p *PageIterator[T]) Next() ([]T, bool) {
	if p.done {
		return nil, false
	}

	items, err := p.fetchFn(domain.PaginationParams{
		Limit:  p.pageSize,
		Offset: p.offset,
	})
	if err != nil || len(items) == 0 {
		p.done = true
		return nil, false
	}

	p.offset += len(items)
	if len(items) < p.pageSize {
		p.done = true
	}

	return items, true
}

// Collect gathers all results across all pages into a single slice.
func (p *PageIterator[T]) Collect() ([]T, error) {
	var all []T
	for {
		page, ok := p.Next()
		if !ok {
			break
		}
		all = append(all, page...)
	}
	return all, nil
}

// IdempotencyKey returns the input if non-empty, otherwise generates a
// timestamp-based key. In production, prefer a UUID library.
func IdempotencyKey(input string) string {
	if input != "" {
		return input
	}
	return fmt.Sprintf("idem-%d", time.Now().UnixNano())
}
