// Package logging はロギングに関するパッケージ
package logging

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minguu42/opepe/pkg/ttime"
)

var (
	debugLogger = log.New(os.Stdout, "\x1b[32mDBG\x1b[0m ", 0)
	infoLogger  = log.New(os.Stdout, "\x1b[36mINF\x1b[0m ", 0)
	errorLogger = log.New(os.Stderr, "\x1b[31mERR\x1b[0m ", 0)
)

// Debugf はDebugレベルのログを出力する
func Debugf(ctx context.Context, format string, v ...any) {
	prefix := fmt.Sprintf("\u001B[33m%s\u001B[0m ", ttime.Now(ctx).Format("2006/01/02 15:04:05"))
	debugLogger.Printf(prefix+format, v...)
}

// Infof はInformationレベルのログを出力する
func Infof(ctx context.Context, format string, v ...any) {
	prefix := fmt.Sprintf("\u001B[33m%s\u001B[0m ", ttime.Now(ctx).Format("2006/01/02 15:04:05"))
	infoLogger.Printf(prefix+format, v...)
}

// Errorf はErrorレベルのログを出力する
func Errorf(ctx context.Context, format string, v ...any) {
	prefix := fmt.Sprintf("\u001B[33m%s\u001B[0m ", ttime.Now(ctx).Format("2006/01/02 15:04:05"))
	errorLogger.Printf(prefix+format, v...)
}

// Fatalf はErrorレベルのログを出力した後、プログラムを終了する
func Fatalf(ctx context.Context, format string, v ...any) {
	prefix := fmt.Sprintf("\u001B[33m%s\u001B[0m ", ttime.Now(ctx).Format("2006/01/02 15:04:05"))
	errorLogger.Fatalf(prefix+format, v...)
}
