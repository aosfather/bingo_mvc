package bingo_mvc

type BingoError interface {
	error
	Code() int
}
