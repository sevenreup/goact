package main

import (
	"bytes"
	"fmt"
	"github.com/sevenreup/goact"
	"log"
)

func main() {
	engine := goact.CreateGoatEngine("./dist", "./view")
	var buf bytes.Buffer
	err := engine.Render(&buf, "./example/entry.jsx", "")
	if err != nil {
		log.Panic(err)
	}
	s := buf.String()
	fmt.Println(s)
}
