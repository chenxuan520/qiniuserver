package utils

import (
	"math/rand"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func CreateFileName(format string, fileOriginName string) string {
	// replace all space to underscore
	fileOriginName = strings.ReplaceAll(fileOriginName, " ", "_")
	if format == "" {
		return fileOriginName
	}

	date := time.Now().Format("20060112150405")
	randomStr := GenerateRandomString(6)

	format = strings.ReplaceAll(format, "%d", date)
	format = strings.ReplaceAll(format, "%r", randomStr)
	format = strings.ReplaceAll(format, "%f", fileOriginName)
	return format
}

func GenerateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
