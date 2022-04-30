package common

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/TheCodeTeam/goodbye"
	"github.com/sirupsen/logrus"

	"github.com/draeron/gof1/pkg/device"
	"github.com/draeron/gof1/pkg/f1"
	"github.com/draeron/gof1/pkg/layout"
	"github.com/draeron/gopkgs/logger"
)

func Setup() {
	logrus.SetLevel(logrus.DebugLevel)
	device.SetLogger(logger.NewLogrus("device"))
	layout.SetLogger(logger.NewLogrus("layout"))
	f1.SetLogger(logger.NewLogrus("f1"))
}

func WaitExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	ctx := context.Background()
	goodbye.Notify(ctx)
	// defer goodbye.Exit(ctx, -1)

	signal.Notify(sigs)
	go func() {
		sig := <-sigs
		fmt.Printf("receive signal %v\n", sig)
		done <- true
	}()

	<-done
}

func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
