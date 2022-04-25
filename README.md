Tinbu XML Lottery feed client
-----------------------------

[![Godoc](https://godoc.org/bitbucket.org/advbet/tinbu?status.svg)](https://godoc.org/bitbucket.org/advbet/tinbu)

This is a go library for accessing Tinbu XML Lottery feed.

Usage example (full state):

```go
package main

import (
        "context"
        "fmt"

        "github.com/advbet/tinbu"
)

func main() {
        c := tinbu.Client{
                URL: "http://example.net/lotterydata/lottery.xml",
        }

        games, err := c.Load(context.TODO())
        fmt.Println(games, err)
}
```

Usage example (update stream):

```go
package main

import (
        "context"
        "fmt"
        "net/http"
        "time"

        "github.com/advbet/tinbu"
)

func main() {
        c := tinbu.Client{
                URL: "http://example.net/lotterydata/lottery.xml",
                HTTPClient: http.Client{
                        Timeout: 10 * time.Second,
                },
        }

        stream := c.StreamUpdates(context.TODO(), time.Minute)
        for msg := range stream {
                if msg.Error != nil {
                        fmt.Println(msg.Error)
                        continue
                }
                fmt.Println(msg.Game)
        }
}
```
