package config

import (
	"fmt"
	"log"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	loggerCustom, err := NewLogger()
	if err != nil {
		log.Fatalf("Error creating Logger: %v", err)
	}

	loggerCustom.Logger.Info("MSG 1 ")

	fmt.Println(loggerCustom.ModelLogger.Application)

	loggerCustom.Logger.Info("MSG 2 ")

	fmt.Println(loggerCustom.ModelLogger.ProcessingTime)

	loggerCustom.Logger.Info("Test")

	loggerCustom.Set("ip", "10.10.223.1")

	loggerCustom.Logger.Info("Test 2 ")
}
