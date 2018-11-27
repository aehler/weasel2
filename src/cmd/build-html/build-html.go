package main

import (
	//	"bytes"
	"flag"
	"fmt"
	"github.com/flosch/pongo2"
	//"github.com/tdewolff/minify"
	//"github.com/tdewolff/minify/html"
	//"github.com/akdcode/auth/rights"
	//"github.com/akdcode/modules/Go/services/access_logs"
	//"github.com/akdcode/modules/Go/services/classifiers/static"
	//"github.com/akdcode/modules/Go2/procedure/form/lib"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"math/rand"
)

var templateSrc = ""
var templateDst = ""

func init() {
	flag.StringVar(&templateSrc, "src", "", "")
	flag.StringVar(&templateDst, "dst", "", "")
}

func main() {

	flag.Parse()

	if _, err := os.Stat(templateSrc); os.IsNotExist(err) {

		panic(fmt.Sprintf("Нет %s", templateSrc))
	}

	filepath.Walk(templateSrc, func(path string, fi os.FileInfo, _ error) error {

		if fi.IsDir() {

			if strings.Index(path, "internal") == -1 {

				if err := os.MkdirAll(templateDst+strings.TrimPrefix(path, templateSrc), 0755); err != nil {

					panic(err)
				}
			}

		} else if strings.HasSuffix(path, ".html") {

			tpl := pongo2.Must(pongo2.FromFile(path))

			if content, err := tpl.Execute(context()); err == nil {

				tplName := templateDst + "/" + strings.TrimPrefix(strings.TrimPrefix(path, templateSrc), "/")
				/*
					min := minify.NewMinifier()
					min.Add("text/html", html.Minify)
					buffer := &bytes.Buffer{}
					if err := min.Minify("text/html", buffer, bytes.NewBufferString(content)); err != nil {
						panic(err)
					}
				*/
				ioutil.WriteFile(tplName, replaceMultipleWhitespace(content), 0644)

			} else {

				panic(err.Error())
			}

		}

		return nil
	})
}

func context() pongo2.Context {

	return pongo2.Context{

		//"globals": map[string]interface{}{
		//	"access_log_types": access_logs.Types,
		//
		//	/*"form": map[string]interface{}{
		//		"elements":       lib.Elements,
		//		"option_loaders": lib.OptionLoaders,
		//		"formats":        lib.Formats,
		//	},*/
		//	"classifiers": map[string]interface{}{
		//		"timezones": static.TimeZones(),
		//		"smb":       static.SmallMediumBusinesses(),
		//		"configurations_property": static.ConfigurationsProperty(),
		//		"legal_forms":             static.LegalForms(),
		//	},
		//	"rights": rights.Rights,
		//},
		"randJS": rand.Int31(),
	}
}

func replaceMultipleWhitespace(html string) []byte {

	return []byte(html)

	s := []byte(html)
	i := 0
	t := make([]byte, len(s))
	previousSpace := false
	for _, x := range s {
		if isWhitespace(x) {
			if !previousSpace {
				previousSpace = true
				t[i] = ' '
				i++
			}
		} else {
			previousSpace = false
			t[i] = x
			i++
		}
	}
	return t[:i]
}
func isWhitespace(x byte) bool {
	return x == ' ' || x == '\t' || x == '\n' || x == '\r' || x == '\f'
}