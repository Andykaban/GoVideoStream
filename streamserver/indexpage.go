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
    <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js"></script>
		<script>
            function isMobileMac() {
                if ( navigator.userAgent.match(/iPhone/i)
                     || navigator.userAgent.match(/iPad/i)
                     || navigator.userAgent.match(/iPod/i))
                   {
                     return true;
                   } 
                   else {
                     return false;
                   }
            }
			$(document).ready(function() {
				// Active camera will refresh every 2 seconds
				var TIMEOUT = 1000;
                if (isMobileMac())
                {
				    var refreshInterval = setInterval(function() {
					    $('img#camera').attr('src', '/stream');
				    }, TIMEOUT);
                }
                else {
                    $("img#camera").attr("src", "/stream");
                }
			});
		</script>
  </head>
  <body>
    <h1>Video Streaming Demonstration</h1>
    <img src="" id="camera">
  </body>
</html>
`

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	indexTemplate, err := template.New("index").Parse(indexPage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	indexTemplate.Execute(w, nil)
}
