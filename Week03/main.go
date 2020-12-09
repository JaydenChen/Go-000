/**
 * @Author chenjun
 * @Date 2020/12/9
 **/
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	//ctx:=context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	group, _ := errgroup.WithContext(ctx)
	//正常退出
	group.Go(func() error {
		return normalExit(ctx.Done(),cancel)
	})
	//信号退出
	group.Go(func() error {
		return signalExit(cancel)
	})

	if err := group.Wait(); err != nil {
		fmt.Println("Successfully exit.")
	}
}

func normalExit(exitChan <-chan struct{},cancel context.CancelFunc) error {
	defer func() {
		log.Print("关闭连接")
		cancel()
	}()
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: mux,
	}
	go func() {
		select {
		case <-exitChan:
			//_, cancel := context.WithTimeout(context.TODO(),time.Second)
			cancel()
			fmt.Errorf("服务关闭")
		}
	}()
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}
	mux.HandleFunc("/hello", helloHandler)
	server.ListenAndServe()

	return fmt.Errorf("http服务关闭")

}
func signalExit(cancel context.CancelFunc) error {
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	select {
	case sign:=<-interrupt:
		cancel()
		err:=fmt.Errorf("异常退出:%v",sign)
		return err
	}

}
