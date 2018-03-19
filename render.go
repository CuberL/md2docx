package main

import (
	"bytes"
	"fmt"
	"strings"

	"baliance.com/gooxml/common"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	//	"baliance.com/gooxml/schema/soo/ofc/sharedTypes"
	"baliance.com/gooxml/schema/soo/wml"
	bf "gopkg.in/russross/blackfriday.v2"
)

type DocxRenderer struct {
	doc           *document.Document
	w             bytes.Buffer
	inHeading     bool
	inParagraph   bool
	inCode        bool
	inItem        bool
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
		} else if d.inItem {
			d.para.AddRun().AddText(string(node.Literal))
		}
	case bf.Paragraph:
		if entering && !d.inItem {
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
	case bf.Item:
		if entering {
			d.inItem = true
			d.para = d.doc.AddParagraph()
			d.para.SetStyle("Item")
		} else {
			d.inItem = false
		}
	case bf.Image:
		if entering {
			dest := string(node.LinkData.Destination)
			fmt.Println(dest)
			img, err := common.ImageFromFile(dest)
			if err != nil {
				fmt.Println(err)
				break
			}
			iref, err := d.doc.AddImage(img)
			if err != nil {
				fmt.Println(err)
				break
			}
			para := d.doc.AddParagraph()
			run := para.AddRun()
			anchored, err := run.AddDrawingInline(iref)
			if err != nil {
				fmt.Println(err)
			}
			para.Properties().SetAlignment(wml.ST_JcCenter)
			//			run.Properties().SetVerticalAlignment(sharedTypes.ST_VerticalAlignRunSubscript)
			anchored.SetSize(iref.RelativeWidth(3*measurement.Inch), 3*measurement.Inch)
			//			anchored.SetName(string(node.LinkData.Title))
			//			anchored.SetHAlignment(wml.WdST_AlignHCenter)
			//			anchored.SetYOffset(3 * measurement.Inch)
			//			anchored.SetTextWrapSquare(wml.WdST_WrapTextLeft)
			//			anchored.SetSize(iref.RelativeWidth(3*measurement.Inch), 3*measurement.Inch)
		}
	}
	return bf.GoToNext
}

func (d *DocxRenderer) Render(ast *bf.Node) []byte {
	d.inCode = false
	d.inHeading = false
	d.inItem = false
	d.inParagraph = false

	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		return d.RenderNode(node, entering)
	})
	return d.w.Bytes()
}
