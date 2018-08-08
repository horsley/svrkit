package svrkit

import (
	"html"
	"regexp"
)

var htmlTagRe = regexp.MustCompile("<[^>]*>")

//StripHTMLTags 去除文本中的 html 标记，会做 html 实体转义
func StripHTMLTags(src string) string {
	rst := htmlTagRe.ReplaceAllString(src, "")
	return html.UnescapeString(rst)
}
