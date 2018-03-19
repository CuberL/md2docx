// md2docx project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"baliance.com/gooxml/document"
	bf "gopkg.in/russross/blackfriday.v2"
)

func main() {
	mdf, err := os.Open("/home/cuberl/gopath/src/connect-core/docs/论文.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	input, err := ioutil.ReadAll(mdf)

	//	input = []byte(`
	//- ListItem1
	//- ListItem2
	//`)
	doc, err := document.OpenTemplate("/home/cuberl/backup2.docx")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	for _, s := range doc.Styles.Styles() {
		fmt.Printf("%s: (%s)\n", s.Name(), s.StyleID())
	}

	renderer := &DocxRenderer{doc: doc}

	extension := bf.FencedCode

	md := bf.New(bf.WithExtensions(extension))
	ast := md.Parse(input)
	renderer.Render(ast)
	err = renderer.doc.SaveToFile("/home/cuberl/new.docx")
	if err != nil {
		fmt.Println(err)
	}
}
