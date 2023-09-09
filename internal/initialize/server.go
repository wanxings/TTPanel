package initialize

import (
	"TTPanel/internal/global"
	"context"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run() {
	Router := Routers()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", global.Config.System.PanelPort),
		Handler:        Router,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// 服务连接
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("run failed: %s", err))
		}
		// 加载自签名证书
		//s.TLSConfig = &tls.Config{
		//	MinVersion: tls.VersionTLS12, // 最小支持TLS 1.2
		//}
		//err := s.ListenAndServeTLS("cert.pem", "key.pem")
		//if err != nil {
		//	panic(err)
		//}
	}()
	_, _ = fmt.Fprintf(color.Output, "TTPanel service listen on %s,Version:%s\n",
		color.GreenString(fmt.Sprintf(":%d", global.Config.System.PanelPort)),
		color.GreenString(global.Version),
	)

	// 等待中断信号关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	_, _ = fmt.Fprintf(color.Output, "%s",
		color.RedString("Shutdown TTPanel...   "),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("TTPanel Shutdown Error:%s", err))
	}
	_, _ = fmt.Fprintf(color.Output, "%s\n",
		color.GreenString("done"),
	)
}
