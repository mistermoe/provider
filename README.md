# Provider

Smol lib that can be used to create handles for resource initialization. 

# Why does this exist?
In most go services, several resources are initialized and then used throughout various parts of the code. This can be a database connection, loggers, caches, api clients that need to be initialized with secrets etc.

Developers either initialize all resources and then pass them around or use global variables. the latter is not recommended as a means to avoid global state pollution. The former makes it quite difficult to manage the resource from a single location e.g. Secrets were rolled and now all resources using said secret need to be re-initialized.


This library aims to provide a way to manage the initialization of resources from a single location and then provide a way to access these resources from anywhere in the codebase. It also provides helper functions to make testing easier.


# Usage

Here's an overly simplified example of how to use the library.

```go
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

client := provider.Provide(func(ctx context.Context) (Client, error) {
  return &RealClient{}, nil
})

func someFunc() {
  // ... stuff ...
  client.Get(ctx).DoThing()

  // ... stuff ...
}

func SomeOtherFunc() {
  // ... other stuff ...
  client.Get(ctx).DoThing()
  // ... other stuff ...
}
```

## In Tests
provider comes with functions that can make testing easier. If you're testing a function that uses a provider and you want to temporarily return a fake resource, you can use `Temporarily` to do so.

```go
func TestSomeOtherFunc(t *testing.T) {
  client.Temporarily(func (ctx context.Context) (Client, error) {
    return &FakeClient{}, nil
  })

  defer client.Reset()

  idk := SomeOtherFunc()
  // ... assertions ...
}
```