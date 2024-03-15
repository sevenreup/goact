package types

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func StartTsGenerator(structPath string, typeOutputFolder string) {
	structs, err := readGoStructs(structPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	build := strings.Builder{}
	for _, str := range structs {
		raw := GenerateTs(&str)
		build.WriteString(raw)
	}
	file, err := os.Create(fmt.Sprintf("%s/generate.d.ts", typeOutputFolder))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(build.String())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readGoStructs(structPath string) (structs []ast.TypeSpec, err error) {
	info, err := os.Stat(structPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !info.IsDir() {
		if strings.HasSuffix(info.Name(), ".go") {
			structs, err = parseAndExtract(structPath)
			if err != nil {
				return nil, err
			}
			return structs, nil
		} else {
			fmt.Println("Error: The provided path is not a Go file.")
		}
		return
	}

	err = filepath.Walk(structPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			strs, err := parseAndExtract(path)
			if err == nil {
				err = nil
				structs = append(structs, strs...)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return
}

func parseAndExtract(filePath string) (structs []ast.TypeSpec, err error) {
	fSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	for _, decl := range astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					structs = append(structs, *typeSpec)
				}
			}
		}
	}
	return
}
