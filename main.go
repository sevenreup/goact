package main

import (
	"fmt"
	"github.com/buke/quickjs-go"
	esbuild "github.com/evanw/esbuild/pkg/api"
	"strings"
)

type BuildFiles struct {
	Js string
}

func main() {
	opts := esbuild.BuildOptions{
		EntryPoints:       []string{"entry.jsx"},
		Outdir:            "./dist",
		Platform:          esbuild.PlatformNode,
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             false,
	}
	result := esbuild.Build(opts)
	if len(result.Errors) != 0 {
		fmt.Println(result.Errors)
	}
	var build BuildFiles
	for _, file := range result.OutputFiles {
		if strings.HasSuffix(file.Path, ".js") {
			build.Js = string(file.Contents)
			fmt.Println("js")
		} else {
			fmt.Println("not js", file.Path)
		}
	}

	html, err := renderReactToHTMLNew(build.Js)
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	fmt.Println(html)
}

func renderReactToHTMLNew(js string) (string, error) {
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
