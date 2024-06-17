package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()
	gr, _ := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i := 0; i < 4; i++ {
		gr.Go(func() error {
			makeWork(ctx, i)
			println("Good")
			if i == 0 {
				return errors.New("New")
			}
			return nil
		})
	}
	err := gr.Wait()
	if err != nil {
		panic("group error " + err.Error())
	}
}

func makeWork(ctx context.Context, i int) {
	println("start ", i)
	ti := time.NewTimer(time.Second * time.Duration(i*2))
	select {
	case <-ti.C:
		println("Goodbye ", strconv.Itoa(i))
	case <-ctx.Done():
		println("Goodbye from ctx ", strconv.Itoa(i))
	}

}
