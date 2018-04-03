package app

import (
	"app/bindata/templates"
	"github.com/flosch/pongo2"
	"fmt"
	"io/ioutil"
	"log"
)

var Templates = map[string]*pongo2.Template{}

var dir = "/srv/src/weasel/templates/pages"

func InitBinaryTemplates(d string) {

	filters()

	for _, name := range templates.AssetNames() {

		if asset, err := templates.Asset(name); err == nil {

			log.Println("Reading", name)

			tmpl := pongo2.Must(pongo2.FromString(string(asset)))

			Templates[fmt.Sprintf("/%s",name)] = tmpl

		} else {

			log.Fatal("Template not found", name, err.Error())

		}

	}

}

func InitTemplates(d string) {

	filters()

	dir = d

	parseDir("")

}

func parseDir(dirname string) {

	cdir := fmt.Sprintf("%s%s", dir, dirname)

	fi, err := ioutil.ReadDir(fmt.Sprintf("%s%s", dir, dirname))
	if err != nil {
		log.Fatal("Cannot access template dir", cdir)
	}

	for _, file := range fi {

		if file.IsDir() {

			parseDir(fmt.Sprintf("%s/%s", dirname, file.Name()))

		} else {

			tmpl := pongo2.Must(pongo2.FromFile(fmt.Sprintf("%s/%s", cdir, file.Name())))

			Templates[fmt.Sprintf("%s/%s", dirname, file.Name())] = tmpl

			fmt.Println("Added template", fmt.Sprintf("%s/%s", dirname, file.Name()))
		}
	}
}
