package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ontitansshoulder/urlshortener/handler"
	"github.com/pkg/errors"
)

func main() {
	err := mainLogic()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func mainLogic() error {
	var redirectsConfigFile string
	var redisServerURL string
	flag.StringVar(&redirectsConfigFile, "config", "", "file configures the redirects")
	flag.StringVar(&redisServerURL, "redis-server", "", "redis server for redirects configs")
	flag.Parse()

	mux := defaultMux()
	mapHandler := defaultHandlerFunc(mux)
	var handlerFunc http.Handler
	var err error
	if redisServerURL != "" {
		handlerFunc, err = handler.RedisHandler(redisServerURL, mapHandler)
		if err != nil {
			return err
		}
	} else {
		handlerFunc, err = configBasedHandlerFunc(redirectsConfigFile, mapHandler)
		if err != nil {
			return err
		}
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handlerFunc)
	return nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func defaultHandlerFunc(mux *http.ServeMux) http.HandlerFunc {
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	return handler.MapHandler(pathsToUrls, mux)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func configBasedHandlerFunc(configFile string, mapHandler http.Handler) (http.HandlerFunc, error) {
	// Build the YAMLHandler using the mapHandler as the
	// fallback
	configType, config := readConfig(configFile)
	switch strings.ToLower(configType) {
	case "yaml":
		fallthrough
	case "yml":
		yamlHandlerFunc, err := handler.YAMLHandler(config, mapHandler)
		if err != nil {
			return nil, errors.Wrap(err, "create YAMLHandler")
		}
		return yamlHandlerFunc, nil
	case "json":
		jsonHandlerFunc, err := handler.JSONHandler(config, mapHandler)
		if err != nil {
			return nil, errors.Wrap(err, "create JSONHandler")
		}
		return jsonHandlerFunc, nil
	default:
		return nil, fmt.Errorf("config type didn't match any hander %s\n", configFile)
	}
}

func readConfig(file string) (string, []byte) {
	if file == "" {
		return "none", nil
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return "none", nil
	}

	parts := strings.Split(file, ".")
	if len(parts) == 1 {
		return "none", nil
	}
	return parts[len(parts)-1], bytes
}
