package main

import (
	"bytes"
	"fmt"
	"strings"

	"baliance.com/gooxml/document"
	bf "gopkg.in/russross/blackfriday.v2"
)

type DocxRenderer struct {
	doc           *document.Document
	w             bytes.Buffer
	inHeading     bool
	inParagraph   bool
	inCode        bool
	para          document.Paragraph
	headingLevel  int
	headingText   string
	paragraphText string
}

func (d *DocxRenderer) RenderNode(node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	case bf.Heading:
		if entering {
			d.inHeading = true
			d.headingLevel = node.Level
			d.para = d.doc.AddParagraph()
			var style string
			switch d.headingLevel {
			case 1:
				style = "Head1"
			case 2:
				style = "Head2"
			case 3:
				style = "Head3b"
			}
			d.para.SetStyle(style)

		}
	case bf.Text:
		if d.inHeading {
			d.para.AddRun().AddText(string(node.Literal))
		} else if d.inParagraph {
			d.para.AddRun().AddText(string(node.Literal))
		}
	case bf.Paragraph:
		if entering {
			d.inParagraph = true
			d.para = d.doc.AddParagraph()
			d.para.SetStyle("It")
		} else {
			d.inParagraph = false
		}
	case bf.CodeBlock:
		codeSplit := strings.Split(string(node.Literal), "\n")
		for _, line := range codeSplit {
			para := d.doc.AddParagraph()
			para.SetStyle("Code")
			para.AddRun().AddText(line)
		}
	}
	return bf.GoToNext
}

func (d *DocxRenderer) Render(ast *bf.Node) []byte {
	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		fmt.Println(node.Type)
		return d.RenderNode(node, entering)
	})
	return d.w.Bytes()
}
