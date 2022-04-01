package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

type Files struct {
	FileName []string
}

func main() {
	os.Mkdir("files", 0664)
	http.HandleFunc("/upload", IndexHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/upload", http.StatusSeeOther)
	})

	http.ListenAndServe(":80", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir("files")
	var data Files
	for _, f := range files {
		data.FileName = append(data.FileName, f.Name())
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		tmpl, _ := template.ParseFiles("index.html")
		tmpl.Execute(w, data)
	}
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create("files/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, data)
}
