package groupietracker

// We can create a struct with everything including artisti inside then create var artists []Artist

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Relations struct {
	ID            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type LocationsStruct struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type DatesStruct struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

func FetchData(url string) ([]Artists, error) {
	// Get eevrything from json
	// Get from the server - err is if theres an error with the server

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// close to avoid memory leak
	defer resp.Body.Close()

	// response body will be decoded it has to be to the ADRESS of data variable
	// if this error is triggered this means the mapping is incorrect
	var data []Artists
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func FetchLocations(url string) (LocationsStruct, error) {
	resp, err := http.Get(url)
	if err != nil {
		return LocationsStruct{}, err // Return an empty instance of LocationsStruct
	}
	defer resp.Body.Close()

	var data LocationsStruct
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return LocationsStruct{}, err // Return an empty instance of LocationsStruct
	}

	return data, nil // Return the fetched data
}

func FetchDates(url string) (DatesStruct, error) {
	resp, err := http.Get(url)
	if err != nil {
		return DatesStruct{}, err
	}
	defer resp.Body.Close()

	var data DatesStruct

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return DatesStruct{}, err
	}

	return data, nil
}

func FetchArtistData(url string) (Artists, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Artists{}, err
	}
	defer resp.Body.Close()

	var data Artists
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Artists{}, err
	}

	return data, nil
}

func FetchRelationData(url string) (Relations, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Relations{}, err
	}
	defer resp.Body.Close()

	var data Relations
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Relations{}, err
	}

	return data, nil
}

func RenderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		// Handle template parsing error
		handleInternalServerError(w, err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		// Handle template execution error
		handleInternalServerError(w, err)
		return
	}
}

func handleInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl, err := template.ParseFiles("../Static/500.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing 500.html template:", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
		log.Println("Error executing 500.html template:", err)
	}
}

func HandlePageNotFoundError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFiles("../Static/404.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing 500.html template:", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
		log.Println("Error executing 500.html template:", err)
	}
}
