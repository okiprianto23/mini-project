package config

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

var (
	AppConfig Configuration
	Logger    *LoggerCustom
)

type Configuration struct {
	Server      server      `envconfig:"server"`
	Mysql       Mysql       `envconfig:"mysql"`
	Redis       redis       `envconfig:"redis"`
	Log         logs        `envconfig:"logs"`
	Token       token       `envconfig:"token"`
	Credentials credentials `envconfig:"credentials"`
	File        file        `envconfig:"file"`
}

type server struct {
	ResourceID  string `envconfig:"main"`
	Application string `envconfig:"application" default:"xyz-company"`
	Version     string `envconfig:"version" default:"1.0.0"`
	Host        string `envconfig:"host"`
	Port        int    `envconfig:"port" default:"8871"`
}

type Mysql struct {
	Host              string `envconfig:"host" default:"localhost"`
	Username          string `envconfig:"username"`
	Password          string `envconfig:"password"`
	DBName            string `envconfig:"dbname"`
	Port              int    `envconfig:"port" default:"5432"`
	MaxOpenConnection int    `envconfig:"max-open-connection" default:"50"`
	MaxIdleConnection int    `envconfig:"max-idle-connection" default:"10"`
}

type redis struct {
	Host       string `envconfig:"host" default:"localhost"`
	Port       int    `envconfig:"port" default:"6379"`
	DB         int    `envconfig:"db" default:"0"`
	Password   string `envconfig:"password"`
	Username   string `envconfig:"username"`
	MaxRetries int    `envconfig:"max_retries"`
}

type logs struct {
	Level int8 `envconfig:"level" default:"0"`
}

type credentials struct {
	ClientID string `envconfig:"client_id"`
	UserID   int64  `envconfig:"user_id"`
}

type token struct {
	UserKey            string        `envconfig:"user_key"`
	FixedInternalToken string        `envconfig:"fixed"`
	Duration           time.Duration `envconfig:"duration" default:"24h"`
}

type file struct {
	Directory string `envconfig:"directory"`
}

func init() {

	if err := envconfig.Process(
		"main",
		&AppConfig,
	); err != nil {
		log.Fatal(err)
	}

	//Set Logger
	var err error
	Logger, err = NewLogger()
	if err != nil {
		log.Fatalf("Error creating Logger: %v", err)
	}

	//untuk check apakah envnya sudah terpasang dengan benar
	PrintConfig(AppConfig)
}

func PrintConfig(c Configuration) {
	data, _ := json.MarshalIndent(c, "", "\t")
	fmt.Println(string(data))
}
