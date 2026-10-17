package libxslt_test

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/llarhub/libxml2"
	"github.com/llarhub/libxslt"
)

// Adapted from https://github.com/goplus/llcppg/blob/v0.8.0/_llcppgtest/libxslt/demo/withdeplibxml/demo.go.
func TestTransformXML(t *testing.T) {
	libxml2.InitParser()
	defer libxml2.CleanupParser()
	libxslt.Init()
	defer libxslt.CleanupGlobals()

	xml := `<?xml version="1.0"?>
<root>
  <person>
    <name>Alice</name>
    <age>25</age>
  </person>
</root>`
	xslt := `<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:template match="/">
    <html>
      <body>
        <h1>个人信息</h1>
        <p>姓名: <xsl:value-of select="root/person/name"/></p>
        <p>年龄: <xsl:value-of select="root/person/age"/></p>
      </body>
    </html>
  </xsl:template>
</xsl:stylesheet>`

	xmlDoc := libxml2.ReadMemory(c.AllocaCStr(xml), c.Int(len(xml)), nil, nil, 0)
	xsltDoc := libxml2.ReadMemory(c.AllocaCStr(xslt), c.Int(len(xslt)), nil, nil, 0)
	if xmlDoc == nil || xsltDoc == nil {
		t.Fatal("failed to read XML or XSLT document")
	}
	defer xmlDoc.Free()

	stylesheet := libxslt.ParseStylesheetDoc(xsltDoc)
	if stylesheet == nil {
		xsltDoc.Free()
		t.Fatal("failed to parse XSLT stylesheet")
	}
	defer stylesheet.Free()

	result := stylesheet.Apply(xmlDoc, nil)
	if result == nil {
		t.Fatal("failed to apply XSLT stylesheet")
	}
	defer result.Free()

	root := result.RootElement()
	if root == nil {
		t.Fatal("transformation produced no root element")
	}
	content := root.NodeGetContent()
	if content == nil {
		t.Fatal("transformation produced no content")
	}
	got := c.GoString((*c.Char)(unsafe.Pointer(content)))
	for _, want := range []string{"Alice", "25", "个人信息"} {
		if !strings.Contains(got, want) {
			t.Errorf("result does not contain %q: %s", want, got)
		}
	}
	t.Log(got)
}
