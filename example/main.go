package main

import (
	"log"
	"net/http"

	gosimplessg "github.com/K-Sato1995/go-simple-ssg"
	"github.com/radovskyb/watcher"

	"github.com/K-Sato1995/go-simple-ssg/config"
)

func main() {
	baseConfig := config.NewConfig(config.Config{
		SiteInfo: config.SiteInfo{
			Title:       "My custom Blog",
			Description: "This is my custom blog",
		},
	})
	engine := gosimplessg.New(baseConfig)
	go startHMR(baseConfig.TemplatePath, engine)
	serveFiles()
}

func serveFiles() {
	fs := http.FileServer(http.Dir(config.GENERATED_HTML_DIR))
	http.Handle("/", fs)
	log.Println("Serving files on http://localhost:8080...")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func startHMR(templatePath string, engine *gosimplessg.Engine) {
	w := watcher.New()
	w.FilterOps(watcher.Write)
	if err := w.AddRecursive(templatePath); err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case <-w.Event:
				log.Println("Change detected. Rebuilding...")
				engine.Build()
				log.Println("Rebuild completed.")
			case err := <-w.Error:
				log.Println("Watcher error:", err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Start(250); err != nil {
		log.Fatal(err)
	}
}
