package testshlp

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type EnvVar struct {
	Name              string
	DefaultVal        string
	DontUseDefaultVal bool
}

func SetEnvFromFile(envPath string) error {
	if err := godotenv.Overload(envPath); err != nil {
		return errors.Wrap(errors.WithStack(err), "load env file error")
	}

	return nil
}

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

func (ev EnvVar) GetBool() bool {
	return strings.EqualFold("true", ev.GetValOrDefault())
}
