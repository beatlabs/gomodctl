// +build ignore

package main

//We need to embed the license database to detect licenses.
//There is a lot of file embedding libraries mentioned in https://tech.townsourced.com/post/embedding-static-files-in-go/ and https://github.com/avelino/awesome-go#resource-embedding.
//If static files you want to embed are in your project, it is easy. But they aren't easy to use for our case because the file we want to embed is inside a module.
//Also, the license checker has a function that returns file as byte slice.
//So we took some inspiration from https://github.com/flazz/togo and modified a little for our use case.

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
