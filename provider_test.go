package provider_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/mistermoe/provider"
)

type Client interface {
	DoThing() error
}

type RealClient struct{}

func (c *RealClient) DoThing() error {
	return nil
}

type FakeClient struct{}

func (c *FakeClient) DoThing() error {
	return errors.New("fake client error")
}

func TestProvider_Get(t *testing.T) {
	ctx := context.Background()

	clientProvider := provider.Provide(func(ctx context.Context) (Client, error) {
		return &RealClient{}, nil
	})

	client := clientProvider.Get(ctx)
	assert.NotZero(t, client)
	assert.NoError(t, client.DoThing())

	client2 := clientProvider.Get(ctx)
	assert.Equal(t, client, client2)
}

func TestProvider_Temporarily(t *testing.T) {
	ctx := context.Background()

	clientProvider := provider.Provide(func(ctx context.Context) (Client, error) {
		return &RealClient{}, nil
	})

	clientProvider.Temporarily(func(ctx context.Context) (Client, error) {
		return &FakeClient{}, nil
	})

	client := clientProvider.Get(ctx)
	assert.NotZero(t, client)
	assert.Error(t, client.DoThing())

	clientProvider.Reset()
	client = clientProvider.Get(ctx)
	assert.NotZero(t, client)
	assert.NoError(t, client.DoThing())
}

func TestProvider_Reset(t *testing.T) {
	ctx := context.Background()

	clientProvider := provider.Provide(func(ctx context.Context) (Client, error) {
		return &RealClient{}, nil
	})

	client := clientProvider.Get(ctx)
	assert.NotZero(t, client)
	assert.NoError(t, client.DoThing())

	clientProvider.Reset()
	client = clientProvider.Get(ctx)
	assert.NotZero(t, client)
	assert.NoError(t, client.DoThing())
}

func TestProvider_Concurrency(t *testing.T) {
	ctx := context.Background()

	clientProvider := provider.Provide(func(ctx context.Context) (Client, error) {
		return &RealClient{}, nil
	})

	var wg sync.WaitGroup
	numGoroutines := 10
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			client := clientProvider.Get(ctx)
			assert.NotZero(t, client)
			assert.NoError(t, client.DoThing())
		}()
	}

	wg.Wait()
}

func TestProvider_ErrorHandling(t *testing.T) {
	ctx := context.Background()

	clientProvider := provider.Provide(func(ctx context.Context) (Client, error) {
		return nil, errors.New("failed to initialize")
	})

	assert.Panics(t, func() {

		clientProvider.Get(ctx)
	})
}
