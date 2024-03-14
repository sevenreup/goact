package html

const BaseHtmlWithLayout = `
import { renderToString } from "react-dom/server.browser";
import Page from "{{ .FilePath }}";
import Layout from "{{ .LayoutPath }}"
import React from "react";
{{ .Props }}
{{ .Content }}
renderToString(<Layout><Page  {...props}/></Layout>);
`

const BaseHtmlNoLayout = `
import { renderToString } from "react-dom/server.browser";
import Page from "{{ .FilePath }}";
import React from "react";
{{ .Props }}
{{ .Content }}
renderToString(<Page  {...props}/>);
`
