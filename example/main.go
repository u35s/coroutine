package main

import (
	"coroutine"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strings"
	"time"
)

func main() {
	fmt.Printf("main id %v\n", Goid())
	go func() {
		profiler("8888")
	}()
	var num uint64 = 0
	go func() {
		for {
			select {
			case <-time.Tick(5 * 1e9):
				fmt.Printf("num %v, 5s\n", num)
			}
		}
	}()
	for {
		co1 := coroutine.NewCoroutine()
		co1.Run(func() error {
			return nil
		})

		num++
		if num%10000 == 0 {
			fmt.Printf("num %v, id\n", num)
		}
	}
}

func Goid() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return idField
}

func profiler(port string) {
	println(http.ListenAndServe(":"+port, nil))
}
