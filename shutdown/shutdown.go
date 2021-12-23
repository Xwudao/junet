package shutdown

import (
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	wg     sync.WaitGroup
	closer []io.Closer
	mu     sync.Mutex
	once   sync.Once
	closed bool
)

func init() {
	once.Do(func() {
		closer = make([]io.Closer, 0)
		closed = false
		wg.Add(1)
		go handleSignals()
	})
}

func handleSignals() {
	defer wg.Done()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
			mu.Lock()
			closed = true
			for _, c := range closer {
				c.Close()
			}
			mu.Unlock()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func Add(c io.Closer) {
	mu.Lock()
	defer mu.Unlock()
	if closed {
		return
	}
	closer = append(closer, c)
}

func Wait() {
	wg.Wait()
}
