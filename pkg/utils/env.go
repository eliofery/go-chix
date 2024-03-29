package utils

import (
	"log/slog"
	"os"
)

const (
	Local = "local"
	Prod  = "prod"
)

// GetEnv возвращает окружение
// Пример: go run main.go => local
// Пример: go run main.go prod => prod
// Пример: go run main.go foobar => foobar
func GetEnv() string {
	args := os.Args[1:]
	if len(args) > 0 {
		return args[0]
	}

	return Local
}

// PrintEnv выводит используемое окружение в лог
func PrintEnv(log *slog.Logger) {
	log.Info("Используемое окружение", slog.String("env", GetEnv()))
}
