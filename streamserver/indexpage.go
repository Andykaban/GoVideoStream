package streamserver

import (
	"net/http"
	"log"
	"html/template"
)

var indexPage = `
<html>
  <head>
    <title>Video Streaming Demonstration</title>
  </head>
  <body>
    <h1>Video Streaming Demonstration</h1>
    <img src="http://127.0.0.1:1488/stream">
  </body>
</html>
`
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	indexTemplate, err := template.New("index").Parse(indexPage)
	if (err != nil) {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
	indexTemplate.Execute(w, nil)
}
