package frontend

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"encore.app/insurance"
)

var (
	//go:embed dist
	dist embed.FS

	assets, _ = fs.Sub(dist, "dist")
)


type entry struct {
	Name string
	Done bool
	Content string
}

type Infos struct {
	Ts string
	List []*insurance.Insurance
}

//encore:api public raw path=/frontend
func Serve(w http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFS(assets, "index.html", "insurances.html"))
	t.Execute(w, Infos{
		List: insurance.Insurances,
	})
}

//encore:api public raw path=/insurances
func Insurances(w http.ResponseWriter, req *http.Request) {
	insurance.Validate()
	t := template.Must(template.ParseFS(assets, "insurances.html"))
	t.Execute(w, Infos{
		List: insurance.Insurances,
	})
}

//encore:api public raw path=/purchase
func Purchase(w http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFS(assets, "d/index.html"))
	t.Execute(w, nil)
}


//encore:api public raw path=/htmx/time
func Time(w http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFS(assets, "htmx_time.html"))
	t.Execute(w, Infos{
		Ts: time.Now().Format(time.RFC1123Z),
	})
}
