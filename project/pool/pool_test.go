package pool

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type test struct {
}

func (t *test) Do() error {
	time.Sleep(time.Second * 3)
	fmt.Println("hello world")
	return nil
}
func tt() (int, error) {
	return 2, nil
}
func TestPool(t *testing.T) {
	//p := NewPool(14)
	//p.Run()
	//for i := 0; i < 15; i++ {
	//	p.SubMit(func() {
	//		//time.Sleep(time.Second * 3)
	//		fmt.Println("hello world")
	//	})
	//}

	i := 0
	runtime.GOMAXPROCS(1)
	go func() {
		i++
		for {
		}
		fmt.Println("hello1")
	}()
	go func() {
		i++
		//for {
		//}
		fmt.Println("hello2")
	}()
	go func() {
		i++
		for {
		}
		fmt.Println("hello3")
	}()
	runtime.LockOSThread()
	runtime.GC()

	fmt.Println(i, runtime.NumGoroutine())
}
