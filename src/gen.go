package src

import "time"

func Generate(ast *AST) bool {
	start := time.Now()
	end := time.Since(start)
	Log("back end took ", end.Seconds(), "ms")
	return true
}
