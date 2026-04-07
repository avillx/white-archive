package whitearchive

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Mode string

const (
	Restore Mode = "restore"
	Backup  Mode = "backup"
)

type Config struct {
	Mode        Mode
	Directory   string
	CryptoKey   []byte
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
}

func LoadConfig() (Config, error) {
	envs := NewEnvs()

	modeRaw := flag.String("mode", "", "define mode \"restore\" or \"sync\" ")
	dir := flag.String("dir", "/data", "working directory")
	flag.Parse()

	mode, err := parseMode(*modeRaw)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Mode:        mode,
		Directory:   *dir,
		CryptoKey:   []byte(envs.Get("CRYPTO_KEY")),
		S3Endpoint:  envs.Get("S3_ENDPOINT"),
		S3AccessKey: envs.Get("S3_ACCESS_KEY"),
		S3SecretKey: envs.Get("S3_SECRET_KEY"),
		S3Bucket:    envs.Get("S3_BUCKET"),
	}, envs.Err()
}

// helpers
type Envs struct {
	errc error
}

func NewEnvs() Envs {
	return Envs{}
}

func (e *Envs) Get(name string) string {
	value, exist := os.LookupEnv(name)
	if !exist {
		e.errc = errors.Join(e.errc, ErrEmptyEnvVar)
	}
	return value
}

func (e *Envs) Err() error {
	return e.errc
}

func parseMode(s string) (Mode, error) {
	switch s {
	case string(Restore):
		return Restore, nil
	case string(Backup):
		return Backup, nil
	default:
		var zero Mode
		return zero, fmt.Errorf("unknown mode: %s", s)
	}
}
