package main

import (
	"io/ioutil"
	"os"
	"tav/src"
)

const (
	EXTENSION = ".tv"
)

func NewFile(filename string) (*src.File, error) {
	bytes, err := ioutil.ReadFile(filename)
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
	file, err := NewFile(args[1]+EXTENSION)
	if err == nil {
		// build to an executable
		if args[0] == "build" {
			program := src.AheadCompile(file)
			src.BuildExe(args[1], program)
		} else if args[0] == "run" {
			src.JITCompile(file)
		}
	} else {
		src.Log(err)
	}
}
