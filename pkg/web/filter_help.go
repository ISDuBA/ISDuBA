// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package web

// Copy the source into this folder because we cannot embed
// files outside the package dir.
//go:generate cp ../../docs/search.md ./search.md

import (
	"net/http"
	"sync"

	"bytes"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"

	_ "embed"
)

//go:embed search.md
var filterHelpRaw []byte

const startTxt = `# Filter expression`
const helpTmplTxt = `<html lang="en">
<head>
   <meta charset="UTF-8">
   <title>Filter expressions</title>
</head>
<body>
<div>{{ . }}</div>
</body>
</html>`

var filterHelpTmpl = sync.OnceValue(func() *template.Template {
	return template.Must(template.New("help").Parse(helpTmplTxt))
})

var filterHelpTxt = sync.OnceValue(func() []byte {
	start := bytes.Index(filterHelpRaw, []byte(startTxt))
	if start == -1 {
		return nil
	}
	return filterHelpRaw[start:]
})

type htmlStream func(http.ResponseWriter) error

func (hs htmlStream) Render(w http.ResponseWriter) error {
	return hs(w)
}

func (htmlStream) WriteContentType(w http.ResponseWriter) {
	render.HTML{}.WriteContentType(w)
}

func (c *Controller) filterHelp(ctx *gin.Context) {

	help := filterHelpTxt()

	opts := html.RendererOptions{
		Flags: html.FlagsNone,
	}
	renderer := html.NewRenderer(opts)
	output := string(markdown.ToHTML(help, nil, renderer))

	tmpl := filterHelpTmpl()

	ctx.Render(http.StatusOK,
		htmlStream(func(w http.ResponseWriter) error {
			return tmpl.Execute(w, output)
		}))
}
