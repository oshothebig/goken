package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	n   int
	c   int
	url string
)

func init() {
	log.SetFlags(log.Lshortfile)
	flag.IntVar(&n, "n", 1, "number of requests")
	flag.IntVar(&c, "c", 1, "number of clients")
	flag.Parse()
	url = os.Args[len(os.Args)-1]
}

func seq(max int) <-chan int {
	i := make(chan int)
	go func() {
		for {
			i <- max
			max = max - 1
			if max <= 0 {
				close(i)
				break
			}
		}
	}()
	return i
}

func main() {
	var wg sync.WaitGroup
	s := seq(n)

	start := time.Now()
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			for _ = range s {
				resp, err := http.Get(url)
				if err != nil {
					log.Println(resp, err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	total := time.Since(start)
	avg := (float64(total.Nanoseconds()) / float64(n)) / (1000 * 1000)
	rps := (float64(n) / float64(total.Nanoseconds())) * (1000 * 1000)
	log.Println(rps)

	format := `
total time: %.4f [s]
average time: %.4f [ms]
req per sec: %.4f [#/sec]
`

	fmt.Printf(format, total.Seconds(), avg, rps)
}
