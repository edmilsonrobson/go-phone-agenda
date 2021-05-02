package logs

import (
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	logPath := "agendajuju.log"

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}

	InfoLogger = log.New(file, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix|log.LUTC)
	WarningLogger = log.New(file, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix|log.LUTC)
	ErrorLogger = log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix|log.LUTC)
}
