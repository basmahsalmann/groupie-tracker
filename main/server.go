package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	f "groupietracker"

	// "golang.org/x/text/search"
)

func main() {
	http.HandleFunc("/", generateHandler)
	// if were counting from the second slash we need to add it
	http.HandleFunc("/artistDetails.html", artistDetailsHandler)

	// serve := http.FileServer(http.Dir("../Static"))
	// http.Handle("/Static/", http.StripPrefix("/Static/", serve))
	fmt.Println("Server running at http://localhost:8080/")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// index page always be restrcictive only allow get we are not sending anything back or updating
// the first thing we need to check is if the methof is not get return an error
//  writer wries for user r reads to user
// the error message is the header (w.WriteHeader(http.StatusMethodNotAllowed))
// we can tets it with curl -x post localhost.....
// to parse -> tmpl := template.Must(template.ParseFles([[insert path]]))
// then tmpl.Execute(w,nil)

// the second thing to check is unallowed routes
// r.url.path != bla bla bla

// finally we can create the index or homepage in html and serve it

// execute the home page with execute(w, Artists) <-  the variable

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		f.HandlePageNotFoundError(w, errors.New("Invalid Path"))
		return
	}

	data, err := f.FetchData("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Error fetching data from API", http.StatusInternalServerError)
		return
	}

	_, err = f.FetchLocations("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		http.Error(w, "Error fetching data from API", http.StatusInternalServerError)
		return
	}

	_, err = f.FetchDates("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		http.Error(w, "Error fetching data from API", http.StatusInternalServerError)
		return
	}

	input := r.FormValue("Name")
	filtered := make([]f.Artists, 0)
	for _,obj := range data {		

		if strings.Contains(strings.ToLower(obj.Name), strings.ToLower(input)){
			filtered = append(filtered, obj)
		}


	}

	if len(filtered) == 0 {
		input = ""
	}

	if input != ""{
		f.RenderTemplate(w, "../Static/web.html", filtered)
		return
	}

	f.RenderTemplate(w, "../Static/web.html", data)



	//fix:
	//valid search
	//there's an issue when restarting the program it still displays the previous search
}

// new dynamic route
// can't restrict it to a path because its dynamic
// method still get

func artistDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// localhost/artist/4
	// /artist/ is 8 charecters
	// r.url.Path[:8]
	id := r.URL.Query().Get("id")

	// here is checking no one enters a bigger number

	// convet the string id to a int ID
	idnum, err := strconv.Atoi(id)
	if err != nil || idnum <= 0 || idnum >= 53 {
		f.HandlePageNotFoundError(w, err)
		return
	}

	artistData, err := f.FetchArtistData("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		http.Error(w, "Error fetching artist data from API", http.StatusInternalServerError)
		return
	}

	relationData, err := f.FetchRelationData("https://groupietrackers.herokuapp.com/api/relation/" + id)
	if err != nil {
		http.Error(w, "Error fetching relation data from API", http.StatusInternalServerError)
		return
	}

	// tmpl....parse
	// tmpl.Execute(w, Artists[id-1])

	f.RenderTemplate(w, "../Static/artistDetails.html", struct {
		Artists  f.Artists
		Relation f.Relations
	}{
		Artists:  artistData,
		Relation: relationData,
	})

	//get input from search

	
}

// func search(allArtists string, input string){
// 	// artists := r.FormValue("Name")
// 	// allArtists := f.Artists.Name

// 	// id := r.URL.Query().Get("id")

// 	// var data []f.Artists
// 	// err = json.NewDecoder(resp.Body).Decode(&data)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }


// 	if strings.Contains(allArtists, input){
// 		//filter function
// 		fmt.Print(input)
// 	} 

// 	//check for each attribute then include all of them in one page and also in search bar
// }

