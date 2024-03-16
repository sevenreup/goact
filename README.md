# Goact
A simple framework for using React as a templating engine in Golang.
## How it works
We are compiling React in SSR mode, this generates Javascript code. 
When the user call the render function we add the props to the compiled js the run the js. 
This step generates the required html then we can return that as our response.
## Simple setup
Add the package to your project
```bash
go get github.com/sevenreup/goact
```
Create the GoactEngine instance and use it
```go
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
		WorkingDir:       "./views",
		IsDebug:          true,
		StructPath:       "./dto",
		TsTypeOutputPath: "./views/types",
	}
	engine := goact.CreateGoactEngine(&opts)

	var buf bytes.Buffer
	err := engine.Render(&buf, "./entry.tsx", map[string]string{
		"title": "Hello World",
	})
	if err != nil {
		log.Panic(err)
	}
	s := buf.String()
	fmt.Println(s)
}

```
## Setup React
Make sure you install the required packages for React in the folder of your project, you can use the provided cli tool for a quick setup or manually.

## Manual setup
Make sure you have node installed on your device.

Init the package.json 
```bash
npm init -y
```
Install React using your favorite package manager
```bash
npm install -D react react-dom
```
## Using the cli
Install the cli
```bash
go install github.com/sevenreup/goact/goact-cli@latest
```
Then you can init the project
```bash
goact-cli init
```
You can pass extra params like

`--packageManger`: To specify the package manager to use (npm is default)

`--tailwind`: To setup tailwind in your project

`--viewDir`: If you are setting up tailwind pass this to point to the location where the main css file should be created

For example
```bash
goact-cli init --packageManger="pnpm" --tailwind --viewDir="./"
```
## Working with React
Below are some of the special rules (🥲 Limitations) of Pages and Layouts, these only matter for pages and Layouts, other components in your project can be used normally.
### Pages
The engine expects all your pages to have a default export of `Page`;

For Example
```tsx
const Page = () => {}
export default Page;
```
### Layouts
The engine also supports layouts. It will check for layouts at the base of the WorkingDir.
The engine expects all layouts to have a default export of `Layout`;

For example
```tsx
const Layout = ({ children }) => {}
export default Layout;
```


## Generating types for Props to use in React
You need to make sure you have a special folder or file that you put the structs that you are going to use when passing data to the render function.
The struct is used to generate a typescript file that you can import in your tsx code.
Pass the `StructPath` option that points to your struct or struct folder.
Pass the `TsTypeOutputPath` which is where the ts file will be created in.

The type generation should only be done in dev mode so make sure you have a way to tell the Engine the current environment because the generation happens at start.

For example using Environment variables
```go
opts := goact.GoactEngineOpts {
// rest of the code
IsDebug:          os.Getenv("Environment") == "Dev",
// rest of the code
}
```

## What is remaining to make this framework production ready?
- [ ] Prebuild all the view files so that they can be used in production 
- [ ] Add caching of compiled React code ( Speeds up dev )
- [ ] Should be able to ship some reactivity to the browser
- [ ] Support layouts at different levels of the view folder

## Some feature that would be awesome to have
- [ ] Hot reload ( Reload the browser when UI changes)
