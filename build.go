// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vugu/material"
	"github.com/vugu/vugu/distutil"
)

func main() {

	var err error

	material.NewFileWriter("static/css").MustWrite()

	os.MkdirAll("docs/css", 0755)

	distutil.MustExec("vgsassc",
		"-o", "docs/css/main.css",
		"static/css/reset.scss",
		"static/css/globals.scss",
		"static/css/elements.scss",
		"static/css/custom.scss",
	)

	// copy static files
	//distutil.MustCopyDirFiltered(".", *dist, nil)

	// find and copy wasm_exec.js
	distutil.MustCopyFile(distutil.MustWasmExecJsPath(), "docs/wasm_exec.js")

	// run go generate
	cmd := exec.Command("go", "generate")
	cmd.Dir, err = filepath.Abs("./ui")
	distutil.Must(err)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	distutil.Must(err)

	//fmt.Print(distutil.MustExec("go", "generate", "./"))

	// run go build for wasm binary
	//	fmt.Print(distutil.MustEnvExec([]string{"GOOS=js", "GOARCH=wasm"}, "go", "build", "-o", filepath.Join(*dist, "main.wasm"), "."))

	cmd = exec.Command("go", "build", "-o", "docs/main.wasm", "./wasm/main")
	//cmd.Dir, err = filepath.Abs(".")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	distutil.Must(err)
	out, err = cmd.CombinedOutput()
	fmt.Printf("%s", out)
	distutil.Must(err)

}
