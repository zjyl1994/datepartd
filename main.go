package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func main() {
	logger, _ = zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "Message",
			LevelKey:    "Level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "Time",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			CallerKey:    "Caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()
	defer logger.Sync()
	logger.Info("DATEPARTD_START")
	err := loadConfig("config.toml")
	if err != nil {
		logger.Error("LOAD_CONFIG", zap.Error(err))
		return
	}
	createPartitionJob()
	deletePartitionJob()
	c := cron.New()
	_, err = c.AddFunc("0 23 * * *", createPartitionJob)
	if err != nil {
		logger.Error("CREATE_JOB", zap.Error(err))
		return
	}
	_, err = c.AddFunc("0 1 * * *", deletePartitionJob)
	if err != nil {
		logger.Error("CREATE_JOB", zap.Error(err))
		return
	}
	c.Start()
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	<-chSignal
	logger.Info("DATEPARTD_STOP")
}

func createPartitionJob() {
	for _, v := range conf.Database {
		err := createFn(v)
		if err != nil {
			logger.Error("CREATE_PARTITION_ERROR", zap.String("TaskName", v.Name), zap.Error(err))
		} else {
			logger.Info("CREATE_PARTITION_SUCCESS", zap.String("TaskName", v.Name))
		}
	}
}

func deletePartitionJob() {
	for _, v := range conf.Database {
		if v.PurgeDays > 0 {
			err := deleteFn(v)
			if err != nil {
				logger.Error("DELETE_PARTITION_ERROR", zap.String("TaskName", v.Name), zap.Error(err))
			} else {
				logger.Info("DELETE_PARTITION_SUCCESS", zap.String("TaskName", v.Name))
			}
		}
	}
}
