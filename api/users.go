package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	// Read username and password request with FormValue.
	creds := model.Credentials{} // TODO: replace this

	
	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	// TODO: answer here
	err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
        return
    }
    
    if creds.Username =="" || creds.Password =="" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
        return
    }

	err = api.usersRepo.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	filepath := path.Join("views", "status.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}


	var data = map[string]string{"name": creds.Username, "message": "register success!"}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}

    w.WriteHeader(http.StatusOK)

}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	// Read usernmae and password request with FormValue.
	
	
	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	// TODO: answer here
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal server error"})
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username =="" || password =="" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
        return
    }

	creds := model.Credentials{
		Password: password,
		Username: username,
	}
	fmt.Println("credential : " , creds)

	err := api.usersRepo.LoginValid(creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Wrong User or Password!"})
		return
	}

	
	// Generate Cookie with Name "session_token", Path "/", Value "uuid generated with github.com/google/uuid", Expires time to 5 Hour.
	// TODO: answer here
	sessionToken := uuid.NewString()
    expiresAt := time.Now().Add(5 * time.Hour)
    http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
		Path: 	 "/",
        Value:   sessionToken,
        Expires: expiresAt, 
    })

	session := model.Session{
		Token: sessionToken,
		Username: creds.Username,
        Expiry:   expiresAt,
	} // TODO: replace this
	err = api.sessionsRepo.AddSessions(session)
    w.WriteHeader(http.StatusOK)
	api.dashboardView(w, r)
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	//Read session_token and get Value:
	// sessionToken := "" // TODO: replace this
	c, err := r.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            json.NewEncoder(w).Encode(model.ErrorResponse{Error: "http: named cookie not present"})
            return
        }
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    session := c.Value

	api.sessionsRepo.DeleteSessions(session)

	//Set Cookie name session_token value to empty and set expires time to Now:
	// TODO: answer here
	http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
        Value:   "",
        Expires: time.Now(),
    })


	filepath := path.Join("views", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}

	
}
