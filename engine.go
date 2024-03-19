package goact

import (
	"encoding/json"
	"errors"
	"github.com/sevenreup/goact/html"
	"github.com/sevenreup/goact/types"
	"github.com/sevenreup/goact/utils"
	"html/template"
	"io"
	"path/filepath"
)

type Views interface {
	Load() error
	Render(io.Writer, string, interface{}, ...string) error
}

type GoactEngine struct {
	opt      *GoactEngineOpts
	compiler *GoactCompiler
}

type GoactEngineOpts struct {
	// The output directory for all the generated files
	OutputDir string
	// The root directory of the views
	WorkingDir string
	// This toggles minification of the generated files
	IsDebug bool
	// Location of structs used to pass props to the template
	StructPath string
	// The output of the auto-gen typescript types
	TsTypeOutputPath string
	// The base layout to use for rendering (This is optional, if you want a custom one make sure it matches the default)
	HtmlBaseLayout string
	// This is the default data to pass to the template on render if not provided in the render function
	DefaultRenderData RenderData
	// To stop the engine from injecting the generated CSS directly into the page (If you have a postprocess step)
	InjectCss bool
}

type RenderData struct {
	// The document title
	Title string
	// Content to insert in the head
	Head template.HTML
	// A list of meta tags to insert in the head
	MetaTags map[string]string
	// Props to pass to the rendered page
	Props interface{}
}

func CreateGoactEngine(opts *GoactEngineOpts) *GoactEngine {
	compiler := NewGoactCompiler(opts.OutputDir, opts.WorkingDir, opts.InjectCss, opts.IsDebug)
	engine := GoactEngine{
		compiler: compiler,
		opt:      opts,
	}
	if opts.IsDebug {
		types.StartTsGenerator(opts.StructPath, opts.TsTypeOutputPath)
	}
	return &engine
}

func (v GoactEngine) Load() error {
	return nil
}

func (v *GoactEngine) Render(writer io.Writer, path string, values interface{}, args ...string) error {
	renderValues, ok := values.(RenderData)
	if !ok {
		return errors.New("Value has to be of type goact.RenderData")
	}
	actualPath := utils.FormatPath(path)
	props, err := propsToJsonString(renderValues.Props)
	if err != nil {
		return err
	}

	layout := v.getLayoutPath()
	var baseLayout string
	if len(v.opt.HtmlBaseLayout) > 0 {
		baseLayout = v.opt.HtmlBaseLayout
	} else {
		baseLayout = html.BaseHtmlLayout
	}
	htmlData, err := v.compiler.Compile(actualPath, layout, props, baseLayout, renderValues, v.opt.DefaultRenderData)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(htmlData))
	if err != nil {
		return err
	}
	return nil
}

func (v *GoactEngine) getLayoutPath() string {
	baseLayoutPath := utils.FormatPath(filepath.Join(v.compiler.workingDir, "layout.tsx"))
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
