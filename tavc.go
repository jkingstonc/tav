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

func NewFile(filename string) (*src.File, error) {
	bytes, err := ioutil.ReadFile("F:\\OneDrive\\Programming\\GO\\src\\tav\\test.tv")
	if err == nil {
		s := string(bytes)
		return &src.File{
			Filename: filename,
			Source:   &s,
		}, nil
	}
	return nil, err
}

func main() {
	src.Log("tav v_a_0_1")
	args := os.Args[1:]

	// read the file into a byte array
	file, err := NewFile("F:\\OneDrive\\Programming\\GO\\src\\tav\\test.tv")
	if err == nil {
		// build to an executable
		if args[0] == "build" {
			program := src.AheadCompile(file)
			src.BuildExe(program)
		} else if args[1] == "run" {
			src.JITCompile(file)
		}
	} else {
		src.Log(err)
	}
}
