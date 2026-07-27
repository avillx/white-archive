package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	whitearchive "white-archive/internal"
)

type Mode string

func run(
	ctx context.Context,
	mode string,
	directory string,
	cryptoKey string,
	s3Endpoint string,
	s3AccessKey string,
	s3SecretKey string,
	s3Bucket string,
) error {

	// validate
	switch {
	case mode == "":
		return fmt.Errorf("missing mode parameter")
	case directory == "":
		return fmt.Errorf("missing directory parameter")
	case cryptoKey == "":
		return fmt.Errorf("missing cryptoKey parameter")
	case s3Endpoint == "":
		return fmt.Errorf("missing s3Endpoint parameter")
	case s3AccessKey == "":
		return fmt.Errorf("missing s3AccessKey parameter")
	case s3SecretKey == "":
		return fmt.Errorf("missing s3SecretKey parameter")
	case s3Bucket == "":
		return fmt.Errorf("missing s3Bucket parameter")
	}

	// context
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	// di
	storageClient, err := whitearchive.NewStorageClient(
		s3Endpoint,
		s3AccessKey,
		s3SecretKey,
		s3Bucket,
	)
	if err != nil {
		return err
	}

	syncer := whitearchive.NewSyncer(
		whitearchive.NewFileService(directory),
		whitearchive.NewCipher([]byte(cryptoKey)),
		storageClient,
	)

	// run
	switch mode {
	case "restore":
		if err := syncer.Restore(ctx); err != nil {
			return err
		}
		slog.Info("files restored")
	case "backup":
		if err := syncer.Backup(ctx); err != nil {
			return err
		}
		slog.Info("files backuped")
	default:
		return fmt.Errorf("mode must be 'restore' or 'backup', got '%s' ", mode)
	}

	return nil
}

func main() {

	modeFlag := flag.String("mode", "", "working mode <restore|backup> ")
	directoryFlag := flag.String("directory", "", "working directory")
	cryptoKeyFlag := flag.String("crypto-key", "", "key for data encryption")
	s3EndpointFlag := flag.String("s3-endpoint", "", "s3 endpoint")
	s3AccessKeyFlag := flag.String("s3-access-key", "", "s3 access key")
	s3SecretKeyFlag := flag.String("s3-secret-key", "", "s3 secret key")
	s3BucketFlag := flag.String("s3-bucket", "", "s3 bucket")
	flag.Parse()

	mode := *modeFlag
	directory := *directoryFlag
	cryptoKey := *cryptoKeyFlag
	s3Endpoint := *s3EndpointFlag
	s3AccessKey := *s3AccessKeyFlag
	s3SecretKey := *s3SecretKeyFlag
	s3Bucket := *s3BucketFlag

	if mode == "" {
		mode = os.Getenv("MODE")
	}

	if directory == "" {
		directory = os.Getenv("DIRECTORY")
	}

	if cryptoKey == "" {
		cryptoKey = os.Getenv("CRYPTO_KEY")
	}

	if s3Endpoint == "" {
		s3Endpoint = os.Getenv("S3_ENDPOINT")
	}

	if s3AccessKey == "" {
		s3AccessKey = os.Getenv("S3_ACCESS_KEY")
	}

	if s3SecretKey == "" {
		s3SecretKey = os.Getenv("S3_SECRET_KEY")
	}

	if s3Bucket == "" {
		s3Bucket = os.Getenv("S3_BUCKET")
	}

	if err := run(
		context.Background(),
		mode,
		directory,
		cryptoKey,
		s3Endpoint,
		s3AccessKey,
		s3SecretKey,
		s3Bucket,
	); err != nil {
		slog.Error("run error", "error", err)
		os.Exit(1)
	}
	slog.Warn("done")
}
