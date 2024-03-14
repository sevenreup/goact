package goact

import (
	"bytes"
	"fmt"
	"github.com/buke/quickjs-go"
	esbuild "github.com/evanw/esbuild/pkg/api"
	html2 "github.com/sevenreup/goact/html"
	"log"
	"strings"
	"text/template"
)

type BuildFiles struct {
	Js string
}

var textEncoderPolyfill = `function TextEncoder(){}TextEncoder.prototype.encode=function(string){var octets=[];var length=string.length;var i=0;while(i<length){var codePoint=string.codePointAt(i);var c=0;var bits=0;if(codePoint<=0x0000007F){c=0;bits=0x00}else if(codePoint<=0x000007FF){c=6;bits=0xC0}else if(codePoint<=0x0000FFFF){c=12;bits=0xE0}else if(codePoint<=0x001FFFFF){c=18;bits=0xF0}octets.push(bits|(codePoint>>c));c-=6;while(c>=0){octets.push(0x80|((codePoint>>c)&0x3F));c-=6}i+=codePoint>=0x10000?2:1}return octets};function TextDecoder(){}TextDecoder.prototype.decode=function(octets){var string="";var i=0;while(i<octets.length){var octet=octets[i];var bytesNeeded=0;var codePoint=0;if(octet<=0x7F){bytesNeeded=0;codePoint=octet&0xFF}else if(octet<=0xDF){bytesNeeded=1;codePoint=octet&0x1F}else if(octet<=0xEF){bytesNeeded=2;codePoint=octet&0x0F}else if(octet<=0xF4){bytesNeeded=3;codePoint=octet&0x07}if(octets.length-i-bytesNeeded>0){var k=0;while(k<bytesNeeded){octet=octets[i+k+1];codePoint=(codePoint<<6)|(octet&0x3F);k+=1}}else{codePoint=0xFFFD;bytesNeeded=octets.length-i}string+=String.fromCodePoint(codePoint);i+=bytesNeeded+1}return string};`
var consolePolyfill = `var console = {log: function(){}};`
var shims = "\nvar process = { env: new Proxy({}, { get: () => '', }) }"

type GoactCompiler struct {
	OutputDir  string
	workingDir string
	IsDebug    bool
}

func NewGoactCompiler(outDir string, workingDir string, isDebug bool) *GoactCompiler {
	compiler := GoactCompiler{
		OutputDir:  outDir,
		workingDir: workingDir,
		IsDebug:    isDebug,
	}
	return &compiler
}

func (g *GoactCompiler) Compile(content string, layout string, props string) (string, error) {
	var tmpl *template.Template

	if len(layout) > 0 {
		temp, err := template.New(html2.BaseHtmlWithLayout).Parse(html2.BaseHtmlWithLayout)
		if err != nil {
			return "", err
		}
		tmpl = temp
	} else {
		temp, err := template.New(html2.BaseHtmlNoLayout).Parse(html2.BaseHtmlNoLayout)
		if err != nil {
			return "", err
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
		return "", err
	}
	fmt.Println(buf.String())
	std := esbuild.StdinOptions{
		Contents:   buf.String(),
		Loader:     esbuild.LoaderTSX,
		ResolveDir: g.workingDir,
	}
	opts := esbuild.BuildOptions{
		Outdir:            g.OutputDir,
		Platform:          esbuild.PlatformNode,
		Metafile:          false,
		Bundle:            true,
		MinifyWhitespace:  !g.IsDebug,
		MinifyIdentifiers: !g.IsDebug,
		MinifySyntax:      !g.IsDebug,
		Write:             true,
		Stdin:             &std,
		JSXDev:            true,
		JSX:               esbuild.JSXAutomatic,
		Loader: map[string]esbuild.Loader{
			".js": esbuild.LoaderJSX,
		},
		Banner: map[string]string{
			"js": textEncoderPolyfill + consolePolyfill + shims,
		},
	}
	result := esbuild.Build(opts)
	if len(result.Errors) != 0 {
		fmt.Println(result.Errors)
	}
	var build BuildFiles
	for _, file := range result.OutputFiles {
		if strings.HasSuffix(file.Path, ".js") {
			build.Js = string(file.Contents)
		} else {
			fmt.Println("not js", file.Path)
		}
	}

	html, err := renderReactToHTMLNewQuick(build.Js)
	log.Print(html)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return html, nil
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
