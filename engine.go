package goact

import (
	"io"
)

type Views interface {
	Load() error
	Render(io.Writer, string, interface{}, ...string) error
}

type GoactEngine struct {
	compiler *GoactCompiler
}

func CreateGoatEngine(outDir string, viewFolder string) *GoactEngine {
	compiler := NewGoactCompiler(outDir, viewFolder)
	engine := GoactEngine{
		compiler: compiler,
	}
	return &engine
}

func (v GoactEngine) Load() error {
	return nil
}

func (v *GoactEngine) Render(writer io.Writer, path string, values interface{}, args ...string) error {
	html, err := v.compiler.Compile(path)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(html))
	if err != nil {
		return err
	}
	return nil
}
