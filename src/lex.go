package src

import "time"

func Lex() []Token {
	start := time.Now()
	end := time.Since(start)
	Log("front end took ", end.Seconds(), "ms")
	return nil
}
