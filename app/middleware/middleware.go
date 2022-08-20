package middleware

import (
	"context"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Image struct {
	UserName    string
	Filename    string
	Size        int64
	ContentType string
	FilePath    string
}

func Middleware(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.Write([]byte("ok"))
		return
	}

	// auth handler

	// router handler

	if r.Method == http.MethodGet {
		getHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		postHandler(w, r)
		return
	}

	w.Write([]byte("method not allowed"))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context = r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1000000)
	file, fileHeader, err := r.FormFile("image")
	select {
	case <-time.After(1 * time.Millisecond):
		if err != nil {
			log.Print(err.Error())
			w.Write([]byte("something wrong with your input. err : " + err.Error()))
			return
		}

		var image *Image = &Image{
			UserName:    r.FormValue("user-name"),
			Filename:    fileHeader.Filename,
			Size:        fileHeader.Size,
			ContentType: fileHeader.Header.Get("content-type"),
		}

		if image.Size > 100000 {
			w.Write([]byte("file to large"))
			return
		}

		filePath := "assets/images/" + "abcde."

		switch strings.Split(image.ContentType, "/")[1] {
		case "jpeg":
			filePath += "jpeg"
		case "jpg":
			filePath += "jpg"
		case "png":
			filePath += "png"
		default:
			w.Write([]byte("format file not allowed"))
			return
		}

		locationFileName := "./" + filePath

		out, err := os.Create(locationFileName)

		if err != nil {
			log.Print(err.Error())
			w.Write([]byte("filed save fole to the system"))
			return
		}

		defer out.Close()

		_, err = io.Copy(out, file)

		if err != nil {
			log.Print(err.Error())
			w.Write([]byte("filed copy file"))
			return
		}

		Profile, err := template.ParseFiles("views/profile.html")

		if err != nil {
			log.Print(err.Error())
			w.Write([]byte("not found"))
			return
		}

		image.FilePath = "http://localhost:9091/" + filePath

		w.Header().Add("content-type", "text/html")

		Profile.Execute(w, image)
		return
	case <-ctx.Done():
		log.Printf("%v Request success canceled ... ", r.Method)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context = r.Context()
	_, err := io.Copy(ioutil.Discard, r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "body not allowed", http.StatusBadRequest)
		return
	}
	select {
	case <-time.After(1 * time.Millisecond):
		if r.URL.Path == "/" {
			file, err := template.ParseFiles("views/index.html")
			if err != nil {
				w.Write([]byte("not found"))
				return
			}

			file.Execute(w, nil)
			return
		}

		if r.URL.Path == "/index.css" {
			file, err := template.ParseFiles("views/index.css")
			if err != nil {
				w.Write([]byte("not found"))
				return
			}

			w.Header().Add("content-type", "text/css")
			file.Execute(w, nil)
			return
		}

		if r.URL.Path == "/index.js" {
			file, err := template.ParseFiles("views/index.js")
			if err != nil {
				w.Write([]byte("not found"))
				return
			}

			w.Header().Add("content-type", "text/javascript")
			file.Execute(w, nil)
			return
		}

		assets := strings.Split(r.URL.Path, "/")

		if assets[1] == "assets" {
			file, err := os.Open("assets/images/" + assets[len(assets)-1])
			if err != nil {
				log.Print(err.Error())
				w.Write([]byte("not found"))
				return
			}
			defer file.Close()
			w.Header().Add("Content-Type", "image/"+strings.Split(assets[len(assets)-1], ".")[1])
			io.Copy(w, file)
			return
		}

		w.Write([]byte("path not found"))
	case <-ctx.Done():
		log.Printf("%v Request success canceled ... ", r.Method)
	}
}
