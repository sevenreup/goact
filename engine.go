package goact

import (
	"encoding/json"
	"fmt"
	"github.com/sevenreup/goact/utils"
	"io"
)

type Views interface {
	Load() error
	Render(io.Writer, string, interface{}, ...string) error
}

type GoactEngine struct {
	compiler *GoactCompiler
}

func CreateGoatEngine(outDir string, workingDir string, isDebug bool) *GoactEngine {
	compiler := NewGoactCompiler(outDir, workingDir, isDebug)
	engine := GoactEngine{
		compiler: compiler,
	}
	return &engine
}

func (v GoactEngine) Load() error {
	return nil
}

func (v *GoactEngine) Render(writer io.Writer, path string, values interface{}, args ...string) error {
	actualPath := fmt.Sprintf("./%s", path)
	props, err := propsToJsonString(values)
	if err != nil {
		return err
	}
	layout := v.getLayoutPath()
	html, err := v.compiler.Compile(actualPath, layout, props)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(html))
	if err != nil {
		return err
	}
	return nil
}

func (v *GoactEngine) getLayoutPath() string {
	baseLayoutPath := fmt.Sprintf("%s/layout.tsx", v.compiler.workingDir)
	exists := utils.FileExists(baseLayoutPath)
	if exists {
		return "./layout.tsx"
	}

	return ""
}

func propsToJsonString(props interface{}) (string, error) {
	if props != nil {
		propsJSON, err := json.Marshal(props)
		return string(propsJSON), err
	}
	return "null", nil
}
