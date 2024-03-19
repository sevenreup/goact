package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type InitAction struct {
	PackageManager string
	UseTailwind    bool
	ViewDir        string
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
	a.execCommand(false, "install", "-D", "react", "react-dom")
	if a.UseTailwind {
		a.setupTailwind()
	}
}

func (a *InitAction) setupTailwind() {
	a.execCommand(false, "install", "-D", "tailwindcss", "postcss", "autoprefixer")

	if _, err := os.Stat("postcss.config.js"); os.IsNotExist(err) {
		fmt.Println("Creating postcss.config.js")
		createPostCSSConfig()
	} else {
		fmt.Println("postcss config is present.")
	}

	if _, err := os.Stat("tailwind.config.js"); os.IsNotExist(err) {
		fmt.Println("Creating tailwind.config.js")
		createTailwindConfig()
	} else {
		fmt.Println("tailwind.config.js is present.")
	}

	mainCSSPath := filepath.Join(a.ViewDir, "main.css")
	if _, err := os.Stat(mainCSSPath); os.IsNotExist(err) {
		fmt.Println("Creating main.css in viewDir.")
		createMainCSS(mainCSSPath)
	} else {
		fmt.Println("main.css is present in viewDir.")
	}
}

func createPostCSSConfig() {
	file, err := os.Create("postcss.config.js")
	if err != nil {
		fmt.Println("Error creating postcss.config.js:", err)
		return
	}
	defer file.Close()

	content := `module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}`
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to postcss.config.js:", err)
		return
	}
}

func createTailwindConfig() {
	file, err := os.Create("tailwind.config.js")
	if err != nil {
		fmt.Println("Error creating tailwind.config.js:", err)
		return
	}
	defer file.Close()

	content := `module.exports = {
  content: ["./**/*.{html,jsx,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [],
}`
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to tailwind.config.js:", err)
		return
	}
}

func createMainCSS(filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating main.css:", err)
		return
	}
	defer file.Close()

	content := `@tailwind base;
@tailwind components;
@tailwind utilities;`
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to main.css:", err)
		return
	}
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
