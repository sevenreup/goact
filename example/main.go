package main

import (
	"bytes"
	"fmt"
	"github.com/sevenreup/goact"
	"log"
)

func main() {
	opts := goact.GoactEngineOpts{
		OutputDir:        "./dist",
		WorkingDir:       "./",
		IsDebug:          true,
		StructPath:       "./dto",
		TsTypeOutputPath: "./types",
	}
	engine := goact.CreateGoactEngine(&opts)
	var buf bytes.Buffer
	err := engine.Render(&buf, "./entry.jsx", "")
	if err != nil {
		log.Panic(err)
	}
	s := buf.String()
	fmt.Println(s)
}
