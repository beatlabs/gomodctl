// +build ignore

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/google/licenseclassifier"
)

const chunkSize = 0x10

const tmpl = `
// Code generated to embed license database. DO NOT EDIT.
package license

var licenseDB = []byte{
	{{ range .Chunks -}}
		{{ range . }} {{ printf "0x%02x" . }}, {{ end }}
	{{ end }}
}
`

type file struct {
	Value []byte
}

func (f *file) Chunks() [][]byte {
	return chunks(f.Value, chunkSize)
}

func chunks(b []byte, n int) [][]byte {
	var c [][]byte

	nChks := len(b) / n

	for i := 0; i < nChks; i++ {
		m := i * n
		c = append(c, b[m:m+n])
	}

	if r := len(b) % n; r > 0 {
		m := n * nChks
		c = append(c, b[m:m+r])
	}

	return c
}

func (f *file) Render(t *template.Template) error {
	outputPath := "license_db.go"

	var buf bytes.Buffer
	if err := t.Execute(&buf, &f); err != nil {
		return err
	}

	b, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(outputPath, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func main() {
	t := template.Must(template.New("constfile").Parse(tmpl))

	licenseFile, err := licenseclassifier.ReadLicenseFile(licenseclassifier.LicenseArchive)
	if err != nil {
		log.Fatal(err)
	}

	f := file{Value: licenseFile}

	if err := f.Render(t); err != nil {
		log.Fatal(err)
	}
}
