package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func (api *API) ImgProfileView(w http.ResponseWriter, r *http.Request) {
	// fileBytes, err := ioutil.ReadFile("./assets/images/img-avatar.png")
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte("File not found"))
	// 	return
	// }
	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/octet-stream")
	// w.Write(fileBytes)
	_, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			error := model.ErrorResponse{
				Error: "http: named cookie not present",
			}

			JsonData, err := json.Marshal(error)
			if err != nil {
				panic(err)
			}
			w.Write(JsonData)
			return
		}
		// Untuk jenis error lainnya, return bad request status
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// View with response image `img-avatar.png` from path `assets/images`

	fileBytes, err := ioutil.ReadFile("./assets/images/img-avatar.png") // membaca file image menjadi bytes
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(fileBytes) // menampilkan image sebagai response
	// TODO: answer here
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
