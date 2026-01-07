package logger

import (
	"fmt"
	"time"
)

func Log(level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	formattedMsg := fmt.Sprintf(msg, args...)
	fmt.Printf("[%s] [%s] %s\n", timestamp, level, formattedMsg)
}