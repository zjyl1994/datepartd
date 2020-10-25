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
	err := loadConfig("config.toml")
	if err != nil {
		logger.Error("LOAD_CONFIG_ERROR", zap.Error(err))
		return
	}
	if !conf.Cron {
		logger.Info("DATEPARTD_RUN_IMMEDIATELY")
		createPartitionJob()
		deletePartitionJob()
		logger.Info("DATEPARTD_DONE")
	} else {
		logger.Info("DATEPARTD_DAEMON_START")
		var tz *time.Location
		if len(conf.Timezone) > 0 {
			tz, err = time.LoadLocation(conf.Timezone)
			if err != nil {
				logger.Error("LOAD_TIMEZONE_FAILED", zap.Error(err))
				tz = time.Local
			}
		} else {
			tz = time.Local
		}
		c := cron.New(cron.WithLocation(tz))
		var createCron, deleteCron string
		if len(conf.CreateCron) > 0 {
			createCron = conf.CreateCron
		} else {
			createCron = "0 23 * * *"
		}
		if len(conf.DeleteCron) > 0 {
			deleteCron = conf.DeleteCron
		} else {
			deleteCron = "0 1 * * *"
		}
		logger.Info("CREATE_PARTITION_JOB", zap.String("Cron", createCron))
		logger.Info("DELETE_PARTITION_JOB", zap.String("Cron", deleteCron))
		_, err = c.AddFunc(createCron, createPartitionJob)
		if err != nil {
			logger.Error("CREATE_JOB", zap.Error(err))
			return
		}
		_, err = c.AddFunc(deleteCron, deletePartitionJob)
		if err != nil {
			logger.Error("CREATE_JOB", zap.Error(err))
			return
		}
		c.Start()
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, os.Interrupt)
		<-chSignal
		logger.Info("DATEPARTD_DAEMON_STOP")
	}
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
