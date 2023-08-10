// Package idgen データの一意な ID 生成を抽象化する
package idgen

type IDGenerator interface {
	Generate() string
}
