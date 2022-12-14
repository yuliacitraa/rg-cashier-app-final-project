package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func (api *API) ImgProfileView(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := ioutil.ReadFile("./assets/images/img-avatar.png")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func (api *API) ImgProfileUpdate(w http.ResponseWriter, r *http.Request) {
	// Update image `img-avatar.png` from path `assets/images`
	// TODO: answer here
	
	if err := r.ParseMultipartForm(1024); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadedFile, handler, err := r.FormFile("file-avatar")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s%s", "img-avatar", filepath.Ext(handler.Filename))

	fileLocation := filepath.Join(dir, "template/assets/images", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.dashboardView(w, r)

}


func (api *API) HandleImage(w http.ResponseWriter, r *http.Request) {
    imageName := r.URL.Query().Get("image-name") // mengambil nama image dari query url
    dir, err := os.Getwd()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	absPath := filepath.Join(dir, "assets/images", imageName)
	fmt.Println(absPath)
	fileBytes, err := ioutil.ReadFile(absPath) // membaca file image menjadi bytes
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("File not found"))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Write(fileBytes) // menampilkan image sebagai response
    return
}