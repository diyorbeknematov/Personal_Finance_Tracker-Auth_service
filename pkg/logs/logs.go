package logs

import (
	"log"
	"log/slog"
	"os"
)

func InitLogger() *slog.Logger {
	logFile, err := os.OpenFile("pkg/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	handler := slog.NewJSONHandler(logFile, nil)

	return slog.New(handler)
}
