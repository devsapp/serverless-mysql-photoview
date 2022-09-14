package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"serverless-terraform-mysql-creator/code/api"
	"serverless-terraform-mysql-creator/code/tool"
	"time"
)

func Init() {
	viper.SetDefault("log-level", logrus.InfoLevel)
	_ = viper.BindEnv("log-level", "SERVERLESS_TERRAFORM_LOG_LEVEL")
}

var logger *logrus.Entry

func start() {
	r := gin.Default()
	r.POST("/invoke", func(ctx *gin.Context) {
		logger = tool.GetLoggerByRequestID(ctx.GetHeader("x-fc-request-id"))
		jobStop := make(chan int)
		timerStop := make(chan int)

		go timingStop(jobStop, timerStop, logger)

		api.Invoke(ctx, logger, jobStop)

		// stop the timer
		if timerStop != nil {
			timerStop <- 1
		}

		jobStop = nil
	})

	if err := r.Run(":9000"); err != nil {
		panic(err)
	}
}

// we must gracefully stop terraform job in case of not upload terraform backend state into oss.
// gracefully stop terraform job:
func timingStop(jobStop chan int, timerStop chan int, logger *logrus.Entry) {
	to := time.NewTimer(28 * time.Minute)
	logger.Info("Begin the timer")
	defer to.Stop()
	select {
	case <-timerStop:
		return
	case <-to.C:

		timerStop = nil
		if jobStop != nil {
			logger.Warn("Begin to gracefully stop the terraform job")
			jobStop <- 1
		}
		time.After(60 * time.Second)
		if jobStop != nil {
			logger.Warn("Begin to force the terraform job to stop")
			jobStop <- 1
		}
		return
	}

}
func main() {
	Init()
	start()
}
