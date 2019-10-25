package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Choose Your Own Adventure</title>
</head>
<body>
  <h1>{{.Title}}</h1>
  {{range .Paragraphs}}
    <p>{{.}}</p>
  {{end}}

  <ul>
    {{range .Options}}
      <li>
        <a href="/{{.Chapter}}">{{.Text}}</a>
      </li>
    {{end}}
  </ul>
</body>
</html>
`

// NewHandler handler for story template generation via HTTP
func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:] // "/intro" => "intro"

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)

		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusBadRequest)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

// JSONStory - method for reading a story contained in JSON
func JSONStory(file io.Reader) (Story, error) {
	d := json.NewDecoder(file)
	var story Story
	if error := d.Decode(&story); error != nil {
		return nil, error
	}

	return story, nil
}

// Story the structure for a certain story
type Story map[string]Chapter

// Chapter chapter for the story
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option probable results of a user action
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
