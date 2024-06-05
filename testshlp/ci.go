package testshlp

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// EnvVar is a type for environment variables
type EnvVar struct {
	Name              string
	DefaultVal        string
	DontUseDefaultVal bool
}

// SetEnvFromFile sets environment variables from file
func SetEnvFromFile(envPath string) error {
	if err := godotenv.Overload(envPath); err != nil {
		return fmt.Errorf("load env file error: %w", err)
	}

	return nil
}

// GetValOrDefault gets environment variable value or default value
func (ev EnvVar) GetValOrDefault() string {
	v, ok := os.LookupEnv(ev.Name)
	if ok {
		return v
	}
	if ev.DontUseDefaultVal {
		return ""
	}
	return ev.DefaultVal
}

// GetBool gets environment variable value as bool
func (ev EnvVar) GetBool() bool {
	return strings.EqualFold("true", ev.GetValOrDefault())
}
