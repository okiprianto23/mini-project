package context

import (
	"context"
	"log"
	"main-xyz/config"
	"main-xyz/constanta"
	"time"
)

func (c *ContextModel) ToContext() context.Context {
	return context.WithValue(
		context.Background(),
		constanta.ApplicationContextConstanta,
		c,
	)
}

func NewContextModel() *ContextModel {
	ctx := new(ContextModel)
	Logger, err := config.NewLogger()
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	ctx.ClientAccess.Logger = *Logger
	return ctx
}

type ContextModel struct {
	AuthAccessTokenModel AuthAccessTokenModel
	Limitation           Limitation
	ClientAccess         ClientAccess
}

type ClientAccess struct {
	IdempotencyKey    string
	Timestamp         time.Time
	ClientAccount     int64
	ClientAccountName string
	Logger            config.LoggerCustom
	Headers           map[string]string
	Path              string
}

type Limitation struct {
	UserID int64
	Other  map[string]interface{}
}
