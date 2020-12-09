package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx1, cancel := context.WithCancel(context.Background())
	g, ctx2 := errgroup.WithContext(ctx1)
	g.Go(func() error {
		server := http.Server{
			Addr:    ":8080",
			Handler: nil,
		}
		http.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
			defer cancel()
			fmt.Println("stop server 8080")
			err := server.Shutdown(context.Background())
			if err != nil {
				fmt.Println("stop err:", err)
			}
		})
		go func() {
			<-ctx2.Done()
			err := server.Shutdown(context.Background())
			if err != nil {
				fmt.Println("shutdown err:", err)
			}
			fmt.Println("shutdown server 8080")
		}()
		return server.ListenAndServe()
	})
	g.Go(func() error {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		// 接收到终止信号 返回错误终止运行
		select {
		case <-signals:
			fmt.Println("receive quit signal")
			return errors.New("receive quit signal")
		case <-ctx2.Done():
			fmt.Println("signal ctx done")
			return ctx2.Err()
		}
	})

	fmt.Println("main running")

	if err := g.Wait(); err != nil {
		fmt.Println("err group wait err:", err.Error())
	}
	fmt.Println("all stopped!")
}
