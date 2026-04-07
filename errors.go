package main

import "errors"

var (
	ErrEmptyFile   = errors.New("download empty file")
	ErrEmptyEnvVar = errors.New("env var is not exist")
	ErrNotFound    = errors.New("is not found")
)
