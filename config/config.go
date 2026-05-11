package config

import "os"

type ConfigStruct struct {
	Addr               string
	Port               string
	Environment        string
	LogLevel           string
	SecretKey          string
	DemoDeployment     bool
	SQLITEDatanasePath string
}

var config ConfigStruct

func init() {
	config = ConfigStruct{
		Addr:               "0.0.0.0",
		Port:               "3000",
		Environment:        "development",
		LogLevel:           "debug",
		DemoDeployment:     true,
		SecretKey:          "HelloWorld",
		SQLITEDatanasePath: "seesaw.db",
	}
	if addr, exsists := os.LookupEnv("ADDR"); exsists {
		config.Addr = addr
	}
	if port, exsists := os.LookupEnv("PORT"); exsists {
		config.Port = port
	}
	if env, exsists := os.LookupEnv("EVNVIRONMENT"); exsists {
		config.Environment = env
	}
	if logLevel, exsists := os.LookupEnv("LOGLEVEL"); exsists {
		config.LogLevel = logLevel
	}
	if secretKey, exsists := os.LookupEnv("SECRET_KEY"); exsists {
		config.SecretKey = secretKey
	}
	if demoDeployment, exsists := os.LookupEnv("DEMO_DEPLOYMENT"); exsists {
		if demoDeployment == "0" {
			config.DemoDeployment = false
		} else {
			config.DemoDeployment = true
		}
	}
	if SQLITEDatanasePath, exsists := os.LookupEnv("SQLITE_DATABASE_PATH"); exsists {
		config.SQLITEDatanasePath = SQLITEDatanasePath
	}
}

func GetConfig() ConfigStruct {
	return config
}
