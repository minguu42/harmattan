package logging

import (
	"log"
	"os"
)

var (
	debugLogger = log.New(os.Stdout, "\x1b[32mDBG\x1b[0m ", log.LstdFlags)
	infoLogger  = log.New(os.Stdout, "\x1b[36mINF\x1b[0m ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "\x1b[31mERR\x1b[0m ", log.LstdFlags)
)

// Debugf はDebugレベルのログを出力する
func Debugf(format string, v ...any) {
	debugLogger.Printf(format, v...)
}

// Infof はInformationレベルのログを出力する
func Infof(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

// Errorf はErrorレベルのログを出力する
func Errorf(format string, v ...any) {
	errorLogger.Printf(format, v...)
}

// Fatalf はErrorレベルのログを出力した後、プログラムを終了する
func Fatalf(format string, v ...any) {
	errorLogger.Fatalf(format, v...)
}
