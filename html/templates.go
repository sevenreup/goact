package html

const BaseReactRenderWithLayout = `
import { renderToString } from "react-dom/server.browser";
import Page from "{{ .FilePath }}";
import Layout from "{{ .LayoutPath }}"
import React from "react";
{{ .Props }}
{{ .Content }}
renderToString(<Layout><Page  {...props}/></Layout>);
`

const BaseReactRenderNoLayout = `
import { renderToString } from "react-dom/server.browser";
import Page from "{{ .FilePath }}";
import React from "react";
{{ .Props }}
{{ .Content }}
renderToString(<Page  {...props}/>);
`

const BaseHtmlLayout = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Title }}</title>
	{{range $k, $v := .Meta}} <meta name="{{$k}}" content="{{$v}}" /> {{end}}
	{{ .Head }}
	<style>
	  {{ .CSS }}
	</style>
</head>
<body>
    {{ .Body }}
	<script type="module">
		{{ .JS }}
	</script>
</body>
</html>
`
