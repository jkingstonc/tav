package main

import (
	"io/ioutil"
	"os"
	"tav/src"
)

const (
	// build and emit executable
	BUILD = 0x0
)

func main() {
	src.Log("tav v_a_0_1")
	args := os.Args[1:]


	bytes, err := ioutil.ReadFile("F:\\OneDrive\\Programming\\GO\\src\\tav\\test.tv");
	if err == nil{
		program := string(bytes)
		if args[0] == "build" {
			src.BuildExe(src.AheadCompile("test.tv", &program))
		}
	}else{
		src.Log("couldn't open test.tv")
	}
}