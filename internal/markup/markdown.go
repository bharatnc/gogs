// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package markup

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	ext "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/lazyregexp"
)

// IsMarkdownFile reports whether name looks like a Markdown file based on its extension.
func IsMarkdownFile(name string) bool {
	extension := strings.ToLower(filepath.Ext(name))
	for _, ext := range conf.Markdown.FileExtensions {
		if strings.ToLower(ext) == extension {
			return true
		}
	}
	return false
}

var validLinksPattern = lazyregexp.New(`^[a-z][\w-]+://|^mailto:`)

// isLink reports whether link fits valid format.
func isLink(link []byte) bool {
	return validLinksPattern.Match(link)
}

func RawMarkdown(body []byte) []byte {
	re := regexp.MustCompile(`((https?|ftp):\/\/|\/)[-A-Za-z0-9+&@#\/%?=~_|!:,.;\(\)]+`)
	md := goldmark.New(
		goldmark.WithExtensions(ext.GFM,
			ext.NewLinkify(
				ext.WithLinkifyURLRegexp(
					re,
				))),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(body, &buf); err != nil {
		return []byte(err.Error())
	}
	return buf.Bytes()
}

// Markdown takes a string or []byte and renders to HTML in Markdown syntax with special links.
func Markdown(input interface{}, urlPrefix string, metas map[string]string) []byte {
	return Render(MARKDOWN, input, urlPrefix, metas)
}
