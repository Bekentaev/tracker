package home

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	artistsbook "art/ArtistsBook"
)

type ErrStatus struct {
	StatusCode   int
	StatusString string
}

var artists []artistsbook.ArtistsStr

func Parser() []artistsbook.ArtistsStr {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &artists)
	return artists
}

func Artists(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errPage(w, r, ErrStatus{http.StatusNotFound, http.StatusText(http.StatusNotFound)})
	}
	t, err := template.ParseFiles("ui/html/index.html")
	if err != nil {
		return
	}

	x := Parser()
	err = t.Execute(w, x)
	if err != nil {
		log.Println("1")
	}
}

func ArtistsPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("ui/html/artists.html")
	if err != nil {
		log.Println("Internal Server Error")
		errPage(w, r, ErrStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		return
	}
	if r.URL.Path != "/artists/" {
		errPage(w, r, ErrStatus{http.StatusNotFound, http.StatusText(http.StatusNotFound)})
		return
	}
	if r.Method != http.MethodGet {
		log.Println("Incotect Method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", http.MethodGet)
		errPage(w, r, ErrStatus{http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id <= 0 || id >= 53 {
		errPage(w, r, ErrStatus{http.StatusNotFound, http.StatusText(http.StatusNotFound)})
		return
	}

	var f artistsbook.General = artistsbook.General{
		ArtistsStr:     idParse(id),
		RelationStruct: RelationParser(id),
	}
	t.Execute(w, f)

	if err != nil {
		errPage(w, r, ErrStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		return
	}
}

func idParse(id int) artistsbook.ArtistsStr {
	res, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%d", id))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var a artistsbook.ArtistsStr
	json.Unmarshal(body, &a)

	return a
}

func RelationParser(id int) artistsbook.RelationStruct {
	res, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", id))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var r artistsbook.RelationStruct
	json.Unmarshal(body, &r)
	return r
}

func errPage(w http.ResponseWriter, r *http.Request, status ErrStatus) {
	w.WriteHeader(status.StatusCode)
	t, err := template.ParseFiles("ui/html/error.html")
	if err != nil {
		log.Println("Iternal error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "errors", status)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
