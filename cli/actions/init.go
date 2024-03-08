package actions

import (
	"log"
	"os"
	"os/exec"
)

type InitAction struct {
	PackageManager string
}

func (a *InitAction) HandleAction() {
	if _, err := os.Stat("./package.json"); err == nil {
		a.runNpmAddPackage()
		return
	}
	a.runNpmInit()
	a.runNpmAddPackage()
}

func (a *InitAction) runNpmInit() {
	a.execCommand(true, "init", "-y")
}

func (a *InitAction) runNpmAddPackage() {
	a.execCommand(false, "install", "-D", "tailwindcss", "postcss", "autoprefixer")
}

func (a *InitAction) execCommand(useNpm bool, arg ...string) {
	var cmd string
	if useNpm {
		cmd = "npm"
	} else {
		cmd = a.PackageManager
	}

	execCMD(cmd, arg...)
}

func execCMD(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	err := cmd.Run()
	if err != nil {
		log.Panic(err)
	}
}
