package goact

import (
	"io"
	"os"
)

type Views interface {
	Load() error
	Render(io.Writer, string, interface{}, ...string) error
}

type GoactEngine struct {
	compiler *GoactCompiler
}

func CreateGoatEngine(outDir string, workingDir string) *GoactEngine {
	compiler := NewGoactCompiler(outDir, workingDir)
	engine := GoactEngine{
		compiler: compiler,
	}
	return &engine
}

func (v GoactEngine) Load() error {
	return nil
}

func (v *GoactEngine) Render(writer io.Writer, path string, values interface{}, args ...string) error {
	dat, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// TODO: Read the layout files
	// TODO: Cache the reads
	html, err := v.compiler.Compile(string(dat), "")
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(html))
	if err != nil {
		return err
	}
	return nil
}
