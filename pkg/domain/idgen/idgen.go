// Package idgen データの一意な ID を生成を生成する
package idgen

//go:generate mockgen -source=$GOFILE -destination=../../../gen/mock/$GOFILE -package=mock

type IDGenerator interface {
	Generate() string
}
