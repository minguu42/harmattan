// Package idgen データの一意な ID 生成を抽象化する
package idgen

//go:generate mockgen -source=$GOFILE -destination=../../gen/mock/$GOFILE -package=mock

type IDGenerator interface {
	Generate() string
}
