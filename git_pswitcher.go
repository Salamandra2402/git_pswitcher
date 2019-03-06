/*
* Sample API with GET and POST endpoint.
* POST data is converted to string and saved in internal memory.
* GET endpoint returns all strings in an array.
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Salamandra2402/git_pswitcher/git"
	"github.com/Salamandra2402/git_pswitcher/profile"
	"github.com/gorilla/mux"
	"github.com/zserge/webview"
)

var (
	// flagPort is the open port the application listens on
	flagPort = flag.String("port", "9000", "Port to listen on")
	server   *http.Server
)

var results []string

// ListHandler returns json with Git profiles
func ListHandler(w http.ResponseWriter, r *http.Request) {
	db := profile.CreateDefaultJsonFileDb()
	profiles, err := db.GetProfiles()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(&profiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(data))
}

// AddHandler add a new Git user profile
func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := profile.CreateDefaultJsonFileDb()
		name := r.FormValue("name")
		email := r.FormValue("email")
		err := db.AddProfile(profile.Profile{Name: name, Email: email}, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// SwitchHandler switches to an another user
func SwitchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		db := profile.CreateDefaultJsonFileDb()
		profile, err := db.GetProfile(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		} else {
			git.SwitchToProfile(profile)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// UpdateHandler updates email for existing user
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db := profile.CreateDefaultJsonFileDb()
		profiles, err := db.GetProfiles()
		name := r.FormValue("name")
		email := r.FormValue("email")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		nameIsFound := false
		for _, item := range profiles {
			if item.Name == name {
				nameIsFound = true
				break
			}
		}
		if !nameIsFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't find the profile to update"))
			return
		}

		err = db.AddProfile(profile.Profile{Name: name, Email: email}, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

func main() {
	results = append(results, time.Now().Format(time.RFC3339))

	muxAPI := http.NewServeMux()
	muxAPI.HandleFunc("/list", ListHandler)
	muxAPI.HandleFunc("/add", AddHandler)
	muxAPI.HandleFunc("/update", UpdateHandler)
	muxAPI.HandleFunc("/switch", SwitchHandler)
	fmt.Println("Sharing API on localhost:9000")
	server = &http.Server{Addr: ":9000", Handler: muxAPI}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Can't start the server")
		}
	}()

	dir := "/public/"
	port := "8000"
	router := mux.NewRouter()
	router.
		PathPrefix("/").
		Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("."+dir))))
	log.Println("Serve web on localhost:" + port)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, router))
	}()
	time.Sleep(1000)
	webview.Open("Git Profile Switcher",
		"http://localhost:8000/web/index.html", 551, 401, false)

}
