package src

import (
	"github.com/llir/llvm/ir"
	"io/ioutil"
	"os/exec"
)

const (
	// emit an exe
	AHEAD_COMPILE uint8 = 0x0
	// JIT compile using LLVM execution engine (used for compile time execution)
	JIT_COMPILE uint8 = 0x1

	SUCCESS_COMP uint8 = 0x0
	SUCCESS_JIT  uint8 = 0x1
	FAIL_COMP    uint8 = 0x2

	TAV_OUT = "tavout/"
)

func BuildExe(filename string, module *ir.Module) uint8 {
	ioutil.WriteFile(TAV_OUT+filename+".ll", []byte(module.String()), 0644)
	c := exec.Command("llc", TAV_OUT+filename+".ll")
	err := c.Run()
	Log("llc err",err)
	c = exec.Command("gcc", "-c", TAV_OUT+filename+".s", "-o", TAV_OUT+filename+".o")
	err = c.Run()
	Log("gcc o err",err)
	c = exec.Command("gcc", TAV_OUT+filename+".o", "-o", TAV_OUT+filename+".exe")
	err = c.Run()
	Log("gcc exe err",err)
	return SUCCESS_COMP
}
