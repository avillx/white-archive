package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM)

	defer stop()

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cipher := NewCipher([]byte(cfg.CryptoKey))
	fileService := NewFileService(cfg.Directory)
	storageClient, err := NewStorageClient(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	syncer := NewSyncer(fileService, storageClient, cipher)

	switch cfg.Mode {
	case Restore:
		if err := syncer.Restore(ctx); err != nil {
			log.Fatal(err)
			return
		}
		log.Print("files restored")
	case Backup:
		if err := syncer.Backup(ctx); err != nil {
			log.Fatal(err)
			return
		}
		log.Print("files backuped")
	}
}
