package bootstrap

import (
	"errors"
	"os"
	"strconv"
)

type ServerConfig struct {
	DatabaseROURL string
	DatabaseRWURL string
	JWTSecret     string
	ServerPort    uint16
}

func LoadConfig() (*ServerConfig, error) {
	roURL, err := requiredEnv("DATABASE_RO_URL")
	if err != nil {
		return nil, err
	}
	rwURL, err := requiredEnv("DATABASE_RW_URL")
	if err != nil {
		return nil, err
	}
	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "3000"
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, errors.New("SERVER_PORT must be a valid port number (0-65535)")
	}

	return &ServerConfig{
		DatabaseROURL: roURL,
		DatabaseRWURL: rwURL,
		JWTSecret:     jwtSecret,
		ServerPort:    uint16(port),
	}, nil
}

func requiredEnv(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", errors.New(name + " must be set")
	}
	return v, nil
}
