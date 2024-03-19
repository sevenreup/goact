package goact

import (
	"bytes"
	"fmt"
	"github.com/buke/quickjs-go"
	esbuild "github.com/evanw/esbuild/pkg/api"
	html2 "github.com/sevenreup/goact/html"
	"html/template"
	"log"
	"strings"
	textTemplate "text/template"
)

type BuildFiles struct {
	Js   string
	Css  esbuild.OutputFile
	Html string
}

type GoactCompiler struct {
	outputDir  string
	workingDir string
	isDebug    bool
	injectCss  bool
}

type Params struct {
	Title string
	Meta  map[string]string
	JS    template.JS
	CSS   template.CSS
	Body  template.HTML
	Head  template.HTML
}

var textEncoderPolyfill = `function TextEncoder(){}TextEncoder.prototype.encode=function(string){var octets=[];var length=string.length;var i=0;while(i<length){var codePoint=string.codePointAt(i);var c=0;var bits=0;if(codePoint<=0x0000007F){c=0;bits=0x00}else if(codePoint<=0x000007FF){c=6;bits=0xC0}else if(codePoint<=0x0000FFFF){c=12;bits=0xE0}else if(codePoint<=0x001FFFFF){c=18;bits=0xF0}octets.push(bits|(codePoint>>c));c-=6;while(c>=0){octets.push(0x80|((codePoint>>c)&0x3F));c-=6}i+=codePoint>=0x10000?2:1}return octets};function TextDecoder(){}TextDecoder.prototype.decode=function(octets){var string="";var i=0;while(i<octets.length){var octet=octets[i];var bytesNeeded=0;var codePoint=0;if(octet<=0x7F){bytesNeeded=0;codePoint=octet&0xFF}else if(octet<=0xDF){bytesNeeded=1;codePoint=octet&0x1F}else if(octet<=0xEF){bytesNeeded=2;codePoint=octet&0x0F}else if(octet<=0xF4){bytesNeeded=3;codePoint=octet&0x07}if(octets.length-i-bytesNeeded>0){var k=0;while(k<bytesNeeded){octet=octets[i+k+1];codePoint=(codePoint<<6)|(octet&0x3F);k+=1}}else{codePoint=0xFFFD;bytesNeeded=octets.length-i}string+=String.fromCodePoint(codePoint);i+=bytesNeeded+1}return string};`
var consolePolyfill = `var console = {log: function(){}};`
var shims = "\nvar process = { env: new Proxy({}, { get: () => '', }) }"

var loaders = map[string]esbuild.Loader{
	".png":   esbuild.LoaderFile,
	".svg":   esbuild.LoaderFile,
	".jpg":   esbuild.LoaderFile,
	".jpeg":  esbuild.LoaderFile,
	".gif":   esbuild.LoaderFile,
	".bmp":   esbuild.LoaderFile,
	".woff2": esbuild.LoaderFile,
	".woff":  esbuild.LoaderFile,
	".ttf":   esbuild.LoaderFile,
	".eot":   esbuild.LoaderFile,
}

func NewGoactCompiler(outDir string, workingDir string, injectCss bool, isDebug bool) *GoactCompiler {
	compiler := GoactCompiler{
		outputDir:  outDir,
		workingDir: workingDir,
		isDebug:    isDebug,
		injectCss:  injectCss,
	}
	return &compiler
}

func (g *GoactCompiler) Compile(content string, layout string, props string, baseHtmlLayout string, renderData RenderData, defaultRenderData RenderData) (string, error) {
	build, err := g.compileReactToHtml(content, layout, props)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(baseHtmlLayout).Parse(baseHtmlLayout)
	if err != nil {
		return "", err
	}

	var title string
	if len(renderData.Title) > 0 {
		title = renderData.Title
	} else {
		title = defaultRenderData.Title
	}

	css := template.CSS("")
	if g.injectCss {
		css = template.CSS(build.Css.Contents)
	}

	var head template.HTML
	if len(renderData.Head) > 0 {
		head = renderData.Head
	} else {
		head = defaultRenderData.Head
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, Params{
		Body:  template.HTML(build.Html),
		Title: title,
		CSS:   css,
		JS:    "",
		Meta:  renderData.MetaTags,
		Head:  head,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *GoactCompiler) compileReactToHtml(content string, layout string, props string) (*BuildFiles, error) {
	var tmpl *textTemplate.Template

	if len(layout) > 0 {
		temp, err := textTemplate.New(html2.BaseReactRenderWithLayout).Parse(html2.BaseReactRenderWithLayout)
		if err != nil {
			return nil, err
		}
		tmpl = temp
	} else {
		temp, err := textTemplate.New(html2.BaseReactRenderNoLayout).Parse(html2.BaseReactRenderNoLayout)
		if err != nil {
			return nil, err
		}
		tmpl = temp
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, map[string]string{
		"FilePath":   content,
		"Content":    "",
		"LayoutPath": layout,
		"Props":      fmt.Sprintf("const props = %s;", props),
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(buf.String())
	std := esbuild.StdinOptions{
		Contents:   buf.String(),
		Loader:     esbuild.LoaderTSX,
		ResolveDir: g.workingDir,
	}
	opts := esbuild.BuildOptions{
		Outdir:            g.outputDir,
		Platform:          esbuild.PlatformNode,
		Metafile:          false,
		Bundle:            true,
		MinifyWhitespace:  !g.isDebug,
		MinifyIdentifiers: !g.isDebug,
		MinifySyntax:      !g.isDebug,
		Write:             true,
		Stdin:             &std,
		JSXDev:            true,
		JSX:               esbuild.JSXAutomatic,
		Loader:            loaders,
		Banner: map[string]string{
			"js": textEncoderPolyfill + consolePolyfill + shims,
		},
	}
	result := esbuild.Build(opts)
	if len(result.Errors) != 0 {
		fmt.Println(result.Errors)
	}
	build := BuildFiles{}

	var jsContent string
	for _, file := range result.OutputFiles {
		if strings.HasSuffix(file.Path, ".js") {
			build.Js = file.Path
			jsContent = string(file.Contents)
		} else if strings.HasSuffix(file.Path, ".css") {
			build.Css = file
		}
	}

	html, err := renderReactToHTMLNewQuick(jsContent)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	build.Html = html

	return &build, nil
}

func renderReactToHTMLNewQuick(js string) (string, error) {
	rt := quickjs.NewRuntime()
	defer rt.Close()
	ctx := rt.NewContext()
	defer ctx.Close()
	res, err := ctx.Eval(js)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
