package main

import (
	"os"

	"example.com/user/ng-go/web"
	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	handler := web.New(logger, &web.Options{ListenAddress: ":8080"})
	handler.Run()
}
