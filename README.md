<div align="center">
  <img class="logo" src="https://raw.githubusercontent.com/amerkurev/gcache/master/logo.svg" alt="gcache | Go caching library"/>
</div>

<div align="center">
    Сoncurrency-safe Go caching library with type safety, multiple cache stores and collecting statistics.
</div>

---
<div align="center">

[![GoVersion](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)&nbsp;
[![Build](https://github.com/amerkurev/gcache/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/amerkurev/gcache/actions/workflows/ci.yml)&nbsp;
[![GoReportCard](https://goreportcard.com/badge/github.com/amerkurev/gcache)](https://goreportcard.com/report/github.com/amerkurev/gcache)&nbsp;
[![CoverageStatus](https://coveralls.io/repos/github/amerkurev/gcache/badge.svg)](https://coveralls.io/github/amerkurev/gcache)&nbsp;
[![GoDoc](https://godoc.org/github.com/amerkurev/gcache?status.svg)](https://godoc.org/github.com/amerkurev/gcache)&nbsp;
[![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/amerkurev/gcache/master/LICENSE)&nbsp;

</div>

## Features

* Multiple cache stores: actually in memory, Redis, SQLite or [your own custom store](#write-your-own-custom-store)
* High concurrent thread-safe access
* A metric cache to let you store metrics about your caches usage (hits, miss, set success, set error, ...)
* An efficient binary marshaler to automatically marshal/unmarshal your cache values
* A well tested and adaptable lightweight pure Go code
* Use of Generics

## Install
gcache requires a Go version with Generics support (Go 1.18 or newer). To install gcache, use `go get`:

    go get github.com/amerkurev/gcache
    
## Quickstart
See it in action:
```go
import (
	"fmt"
	"github.com/amerkurev/gcache"
	"github.com/amerkurev/gcache/store"
)

func main() {
	c := gcache.New[int, string](store.MapStore(0))

	c.Set(1, "Hello World")
	v, _ := c.Get(1)

	fmt.Println(v) // Hello World
}
```

A more complex example:
```go
import (
	"fmt"
	"github.com/amerkurev/gcache"
	"github.com/amerkurev/gcache/store"
	"time"
)

type Key struct {
	ID int
}

type Employee struct {
	Key
	Name string
	DoB  time.Time
}

func main() {
	c := gcache.New[Key, *Employee](store.MapStore(0))

	key := Key{1001}

	c.Set(key, &Employee{
		Key:  key,
		Name: "Amelia",
		DoB:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	})

	v, err := c.Get(key)
	if err == nil {
		fmt.Println(v.ID)   // 1001
		fmt.Println(v.Name) // Amelia
		fmt.Println(v.DoB)  // 2009-11-11 02:00:00 +0300 MSK
	}
}
```

## Built-in stores

### MapStore 
Go builtin map with mutex lock.
```go
import (
	"github.com/amerkurev/gcache"
	"github.com/amerkurev/gcache/store"
)

func main() {
	c := gcache.New[int, string](store.MapStore(0))
	// ...
}
```

### BigcacheStore
[Bigcache](https://github.com/allegro/bigcache) is a fast, concurrent, evicting in-memory cache written to keep big number of entries.
```go
import (
	"github.com/allegro/bigcache/v3"
	"github.com/amerkurev/gcache"
	"github.com/amerkurev/gcache/store"
	"time"
)

func main() {
	bc, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}
	c := gcache.New[int, string](store.BigcacheStore(bc))
	// ...
}
```

### RedisStore
[Redis](https://github.com/go-redis/redis) is an in-memory database that persists on disk.
```go
import (
	"github.com/amerkurev/gcache"
	"github.com/amerkurev/gcache/store"
	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	c := gcache.New[int, string](store.RedisStore(rdb))
	// ...
}
```

### SQLiteStore
SQLite is a lightweight disk-based database that doesn’t require a separate server process.

### Write your own custom store
You also have the ability to write your own custom store by implementing the following interface:
```go
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}
```
Let's do this together:
```go
import (
	"context"
	"github.com/amerkurev/gcache"
)

type MySuperStore struct {
	m map[string][]byte // non-concurrent
}

func (s *MySuperStore) Get(_ context.Context, key string) ([]byte, error) {
	v, ok := s.m[key]
	if !ok {
		return nil, gcache.ErrNotFound
	}
	return v, nil
}

func (s *MySuperStore) Set(_ context.Context, key string, data []byte) error {
	s.m[key] = data
	return nil
}

func (s *MySuperStore) Delete(_ context.Context, key string) error {
	delete(s.m, key)
	return nil
}

func (s *MySuperStore) Clear(_ context.Context) error {
	s.m = make(map[string][]byte)
	return nil
}

func main() {
	store := &MySuperStore{m: make(map[string][]byte, 0)}
	c := gcache.New[int, string](store)
	// ...
}
```

## Example of using metrics
```go
import (
	"fmt"
	"github.com/amerkurev/gcache"
	"github.com/amerkurev/gcache/store"
	"sync"
	"time"
)

func main() {
	c := gcache.New[int, int](store.MapStore(10_000))
	c.UseStats() // enable to store metrics about caches usage

	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < 10_000; k++ {
				c.Set(k, k)
				c.Get(k + 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Took: %v\n", time.Since(start))

	s, _ := c.Stats()
	fmt.Printf("Hits: %d\n", s.Hits)
	fmt.Printf("Miss: %d\n", s.Miss)
	fmt.Printf("ReadCount: %d\n", s.ReadCount)
	fmt.Printf("WriteCount: %d\n", s.WriteCount)
	fmt.Printf("ReadBytes: %d\n", s.ReadBytes)
	fmt.Printf("WriteBytes: %d\n", s.WriteBytes)

	// Our results:
	// Took: 1.295712084s
	// Hits: 989805
	// Miss: 10195
	// ReadCount: 1000000
	// WriteCount: 1000000
	// ReadBytes: 4895973
	// WriteBytes: 4946000
}
```

## Status
The project is under active development and may have breaking changes till v1 is released. However, we are trying our best not to break things unless there is a good reason.

## License
The MIT License
