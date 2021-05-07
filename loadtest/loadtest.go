package loadtest

import (
	"fmt"
	"github.com/spaceapegames/lambda-burst/client"
	"log"
	"sync"
	"time"
)

type LoadTest struct {
	threads  int
	duration time.Duration
	rate     int
	client   client.Client
	verbose  bool
}

type Results struct {
	Count int64
}

func NewLoadTest(threads, rate int, duration int64, client client.Client) LoadTest {
	return LoadTest{
		threads:  threads,
		rate:     rate,
		duration: time.Second * time.Duration(duration),
		client:   client,
		verbose:  true,
	}
}

func (l LoadTest) Run() (Results, error) {
	startTime := time.Now()
	sleepTime := time.Duration(1000/l.rate) * time.Millisecond

	results := Results{}
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i := 0; i < l.threads; i++ {
		wg.Add(1)
		go func() {
			for time.Now().Sub(startTime) <= l.duration {
				go func(){
					statusCode, body, duration, err := l.client.Go()
					if err != nil {
						fmt.Println(err)
					}
					l.logit(fmt.Sprintf("[%d] %s: %d ms", statusCode, body, duration.Milliseconds()))
					lock.Lock()
					results.Count++
					lock.Unlock()
				}()
				time.Sleep(sleepTime)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return results, nil
}

func (l LoadTest) logit(msg string) {
	if l.verbose {
		log.Println(msg)
	}
}
