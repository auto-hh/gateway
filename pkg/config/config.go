package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func MustGet(key string) string { //используем для проверки наличия системной переменной
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("config.MustGet: %s required variable is not set", val))
	}

	return val
}

func Get(key string, defaultVal string) string { //само получение переменной
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func LoadDotEnv(path string) error { //загрузка виртуального окружения из .env файла
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("config.LoadDotEnv: %w", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")

		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("config.LoadDotEnv: %w", err)
			}
		}
	}
	return nil
}
