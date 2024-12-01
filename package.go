package servant

import (
	"embed"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"text/template"
)

func init() {
	htdocs = template.Must(
		template.New("").Funcs(funcMap).ParseFS(asset, "htdocs/*.html"),
	)
}

var htdocs *template.Template
var funcMap = template.FuncMap{
	//"doX": func() string { return "x" },
}

//go:embed htdocs static
var asset embed.FS

var debug = log.New(ioutil.Discard, "D ", log.LstdFlags|log.Lshortfile)

func init() {
	if yes, _ := strconv.ParseBool(os.Getenv("D")); yes {
		debug.SetOutput(os.Stderr)
	}
}
