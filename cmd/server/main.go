package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"qwerty/config"
	"qwerty/logger"
	"qwerty/server"
	"syscall"

	"github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load("etc")
	if err != nil {
		log.Fatal("Could not read config file: ", err)
	}

	args := os.Args
	args[0] = cfg.Server.Name

	dctx := &daemon.Context{
		PidFileName: cfg.Server.PidFile,
		PidFilePerm: 0644,
		LogFileName: cfg.Logger.File,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        os.Args,
	}

	d, err := dctx.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer dctx.Release()

	ctx := context.Background()
	flag.StringVar(&cfg.Server.BindAddress, "l", cfg.Server.BindAddress, "server listening host/port")
	flag.StringVar(&cfg.Logger.File, "c", cfg.Logger.File, "log file")
	flag.Parse()

	logger.Setup(cfg.Logger)
	s := server.NewServer(cfg.Server.BindAddress, log.WithField("logger", "main"))

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		for {
			select {
			case sig := <-signalChan:
				switch sig {
				case syscall.SIGINT, syscall.SIGTERM:
					log.Printf("Got SIGINT/SIGTERM, exiting.")
					if err := s.Stop(); err != nil {
						log.Fatal(err)
					}
					os.Exit(1)
				case syscall.SIGUSR1:
					log.Printf("Got SIGHUP, print debug.")
					s.Debug()
				}
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()
	<-ctx.Done()
}
