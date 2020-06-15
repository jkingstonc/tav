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
)

func BuildExe(module *ir.Module) uint8 {
	Log("fuckballs")
	ioutil.WriteFile("tmp/test.ll", []byte(module.String()), 0644)
	c := exec.Command("llc", "tmp/test.ll")
	err := c.Run()
	Log("llc err",err)
	c = exec.Command("gcc", "-c", "tmp/test.s", "-o", "tmp/test.o")
	err = c.Run()
	Log("gcc o err",err)
	c = exec.Command("gcc", "tmp/test.o", "-o", "tmp/test")
	err = c.Run()
	Log("gcc exe err",err)
	return SUCCESS_COMP
}
