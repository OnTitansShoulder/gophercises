package game

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/ontitansshoulder/textgame/game/html"
	"github.com/ontitansshoulder/textgame/story"
	"github.com/pkg/errors"
)

const (
	PageTitle = "TextGame"
)

type WebGame struct {
	InitArc    string
	CurrentArc string
	ArcsHtml   map[string]string
}

func NewWebGame(initArc string, arcs map[string]story.Arc) (*WebGame, error) {
	if _, ok := arcs[initArc]; !ok {
		return nil, fmt.Errorf("initial story arc %s is not defined", initArc)
	}

	arcTemplate, err := template.New("arcpage").Parse(html.ArcTemplateStr)
	if err != nil {
		return nil, errors.Wrap(err, "parse arcTemplate")
	}

	var writer bytes.Buffer
	arcsHtml := make(map[string]string, len(arcs)+1)
	for name, arc := range arcs {
		err = arcTemplate.Execute(&writer, arc)
		if err != nil {
			return nil, errors.Wrap(err, "execute template on "+name)
		}
		arcsHtml[name] = writer.String()
		writer.Reset()
	}
	arcsHtml["intro"] = arcsHtml[initArc]

	webGame := WebGame{
		InitArc:    initArc,
		CurrentArc: initArc,
		ArcsHtml:   arcsHtml,
	}
	return &webGame, nil
}

func (game *WebGame) Start() error {
	mux := defaultMux(game.InitArc)
	fmt.Println("Starting game at localhost:8080")
	http.ListenAndServe(":8080", arcHandler(game.ArcsHtml, mux))
	return nil
}

func defaultMux(initArc string) *http.ServeMux {
	initArcPath := "/" + initArc

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Host
		http.Redirect(w, r, path+initArcPath, 302)
	})
	return mux
}

func arcHandler(arcsHtml map[string]string, fallback http.Handler) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		path = strings.Trim(path, "/")
		if val, ok := arcsHtml[path]; ok {
			fmt.Fprintln(w, val)
			return
		}
		fallback.ServeHTTP(w, r)
	}
	return handler
}
