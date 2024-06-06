package provider

import (
	"context"
	"fmt"
	"sync"
)

type Provider[T any] struct {
	fn     func(ctx context.Context) (T, error)
	tempFn func(ctx context.Context) (T, error)
	result T
	once   sync.Once
	mu     sync.Mutex
}

// Provide returns a provider that will return the result of the provided function.
// The function will only be executed once, and the result will be cached.
// The result can be reset by calling the Reset method. This function is thread-safe.
func Provide[T any](fn func(ctx context.Context) (T, error)) *Provider[T] {
	return &Provider[T]{fn: fn}
}

// Get returns the result of the function provided to the Provide function.
// The function will only be executed once, and the result will be cached.
// The result can be reset by calling the Reset method. This function is thread-safe.
func (p *Provider[T]) Get(ctx context.Context) T {
	p.once.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		fn := p.fn
		if p.tempFn != nil {
			fn = p.tempFn
		}

		result, err := fn(ctx)
		if err != nil {
			panic(fmt.Errorf("initialization failed: %w", err))
		}
		p.result = result
	})
	return p.result
}

// Reset resets the provider, causing the next call to Get to re-execute the provided function. Resetting is useful in two scenarios:
// 1. When something about the environment has changed and the result of the function may be different e.g. rolled secrets
// 2. In tests to provide a new instance of the result for each test
func (p *Provider[T]) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.once = sync.Once{}
	var zeroValue T
	p.result = zeroValue
	p.tempFn = nil
}

// Temporarily sets a new function to be used for initialization until Reset is called.
// This is only useful when testing functions that use the provider. e.g. returning a mock instead of a real value
func (p *Provider[T]) Temporarily(fn func(ctx context.Context) (T, error)) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.tempFn = fn
	p.once = sync.Once{}
}
