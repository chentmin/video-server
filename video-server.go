/*
	Hosts videos and audios on your computer so that others can view without downloading them first
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

var (
	dir           = http.Dir("./")
	fileHandler   = http.StripPrefix("/files/", logPanic(http.FileServer(dir)))
	videoTemplate = template.Must(template.New("video").Parse(videoTemplateString))
	audioTemplate = template.Must(template.New("audio").Parse(audioTemplateString))
	indexTemplate = template.Must(template.New("index").Parse(indexTemplateString))
	hostAddr      = flag.String("host", ":8080", "Address to listen")
)

type videoInfo struct {
	Name     string
	BaseName string
}

type audioInfo struct {
	Name     string
	BaseName string
}

type indexInfo struct {
	Videos []*videoInfo
	Audios []*audioInfo
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Handle("/files/", http.HandlerFunc(logPanic(redirectHandler)))
	http.Handle("/", http.HandlerFunc(logPanic(indexHandler)))
	http.Handle("/watch/", http.StripPrefix("/watch/", http.HandlerFunc(logPanic(watchHandler))))
	http.Handle("/listen/", http.StripPrefix("/listen/", http.HandlerFunc(logPanic(listenHandler))))
	err := http.ListenAndServe(*hostAddr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

type handlerFunc func(http.ResponseWriter, *http.Request)

func logPanic(function handlerFunc) handlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v] caught panic: %v\n", request.RemoteAddr, x)
			}
		}()
		function(writer, request)
	}
}

/*
	Prevents others from downloading arbitary files. Only media files are allowed.
*/
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	toWatch, err := url.QueryUnescape(r.URL.RequestURI())
	if err != nil || len(toWatch) <= 1 {
		http.NotFound(w, r)
		return
	}

	suffix := filepath.Ext(toWatch)

	if suffix == ".mp4" || suffix == ".m4v" || suffix == ".mp3" || suffix == ".m4a" {
		fileHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
		return
	}
}

func watchHandler(w http.ResponseWriter, r *http.Request) {
	// get the actual watch page
	toWatch, err := url.QueryUnescape(r.URL.RequestURI())
	if err != nil || len(toWatch) <= 1 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// make sure the file really exists
	f, err := dir.Open(toWatch)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	f.Close()

	info := &videoInfo{Name: toWatch, BaseName: filepath.Base(toWatch)}

	videoTemplate.Execute(w, info)
}

func listenHandler(w http.ResponseWriter, r *http.Request) {
	// get the actual listen page
	toListen, err := url.QueryUnescape(r.URL.RequestURI())
	if err != nil || len(toListen) <= 1 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// make sure the file really exists
	f, err := dir.Open(toListen)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	f.Close()

	info := &audioInfo{Name: toListen, BaseName: filepath.Base(toListen)}

	audioTemplate.Execute(w, info)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	videos := make([]*videoInfo, 0, 16)
	musics := make([]*audioInfo, 0, 16)
	walkFn := func(path string, info os.FileInfo, err error) error {
		suffix := filepath.Ext(info.Name())
		if suffix == ".mp4" || suffix == ".m4v" {
			videos = append(videos, &videoInfo{Name: &path, BaseName: filepath.Base(path)})
		} else if suffix == ".mp3" || suffix == ".m4a" {
			musics = append(musics, &audioInfo{Name: &path, BaseName: filepath.Base(path)})
		}
		return nil
	}

	filepath.Walk("./", walkFn)

	totalInfo := &indexInfo{Videos: videos, Audios: musics}

	indexTemplate.Execute(w, totalInfo)
}
