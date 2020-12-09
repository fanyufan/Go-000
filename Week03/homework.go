/*
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。
 */
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	g, errCtx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(errCtx)
	defer cancel()

	http1 := http.Server{Addr:"127.0.0.1:8080"}
	g.Go(func() error {
		if err:=http1.ListenAndServe();err !=nil {
			cancel()
			fmt.Printf("http1 err:%v\n",err)
			return err
		}
		return nil
	})

	http2 := http.Server{Addr:"127.0.0.1:8081"}
	g.Go(func() error {
		if err:=http2.ListenAndServe();err !=nil {
			cancel()
			fmt.Printf("http2 err:%v\n",err)
			return err
		}
		return nil
	})
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			select {
			case s:= <-c:
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					fmt.Printf("cancel:%v\n", s)
					cancel()
				default:

				}
			default:

			}
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("context cancel")
		default:

		}
	}()
	// Wait for all HTTP fetches to complete.
	if err := g.Wait(); err == nil {
		fmt.Println("exit")
	}

}

