package config

import "os"

func GetEnvVar(envName, _default string) string {
	config, exist := os.LookupEnv(envName)
	if !exist {
		return _default
	}
	return config
}
