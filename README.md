# Golang goroutine pool(thread pool)

## Example
 

```go
package main

import (
	"fmt"
	"time"

	"github.com/wuyongjia/threadpool"
)

type request struct {
	n int
}

func testfun(r *request) {
	fmt.Printf("id: %d, time: %d\n", r.n, time.Now().Unix())
	time.Sleep(time.Second)
}

func main() {
	var pool = threadpool.NewWithFunc(50, 100, func(payload interface{}) {
		var r, ok = payload.(*request)
		if ok {
			testfun(r)
		}
	})
	defer pool.Close()

	for i := 0; i < 200; i++ {
		r := &request{n: i}
		pool.Invoke(r)
	}

	fmt.Println("done")
	time.Sleep(time.Second * 5)
}
```