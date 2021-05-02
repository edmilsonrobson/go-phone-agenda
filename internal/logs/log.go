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
	InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix|log.LUTC)
	WarningLogger = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix|log.LUTC)
	ErrorLogger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix|log.LUTC)
}
