package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var StartTime time.Time

func booksBylanguage(language string) (BookCount, string) {
	//Get total amount of books
	res0, err := http.Get("http://129.241.150.113:8000/books/")
	if err != nil {
		return BookCount{}, "Something went wrong."
	}

	var data GutendexResponseCount

	decoder := json.NewDecoder(res0.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return BookCount{}, "Something went wrong during JSON decoding."
	}
	totalBooks := data.Count

	res0.Body.Close()

	// Retrieve book information
	url := "http://129.241.150.113:8000/books?languages=" + language

	// Counting
	bookAmt := 0
	authorAmt := 0

	// Keeping track of authors to avoid duplicates
	countedAuthors := make(map[string]bool)

	for { // Loop through all pages

		res, err := http.Get(url)
		if err != nil {
			return BookCount{}, "Something went wrong."
		}

		// Decode JSON
		var data GutendexResponse
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&data)
		if err != nil {
			return BookCount{}, "Something went wrong during JSON decoding."
		}

		res.Body.Close()

		// Count for current page
		for j := 0; j < len(data.Results); j++ {
			authors := data.Results[j].Authors
			l := len(authors)

			if l > 0 { // If the author is unknown, we don't count the book
				bookAmt++
				for k := 0; k < l; k++ { // Checking all authors
					name := authors[k].Name
					if !countedAuthors[name] { // Checks if author has not been counted yet
						authorAmt++
						countedAuthors[name] = true
					}
				}
			}
		}
		if data.Next == "" { // There are no more pages
			break
		} else {
			url = data.Next // Next page
		}
	}

	fraction := float32(bookAmt) / float32(totalBooks)

	count := BookCount{
		Language: language,
		Books:    bookAmt,
		Authors:  authorAmt,
		Fraction: fraction}

	return count, ""
}

func Bookcount(w http.ResponseWriter, r *http.Request) {
	// Getting parameters (languages)
	languagesString := r.URL.Query().Get("language")

	// A language parameter is required, so if there is none, the user is instructed
	if languagesString == "" {
		instructions := "To use this service, add a language parameter."
		instructions += "\nExample: localhost:8080/librarystats/v1/bookcount/?language=no"
		instructions += "\nYou can write multiple languages as well."
		instructions += "\nExample: localhost:8080/librarystats/v1/bookcount/?language=no,fi"

		_, err := fmt.Fprint(w, instructions)
		if err != nil {
			http.Error(w, "Error during generation of response.", http.StatusInternalServerError)
			return
		}
	}

	languages := strings.Split(languagesString, ",")
	languageAmt := len(languages)

	var count []BookCount

	// Retrieve book information for each language
	for i := 0; i < languageAmt; i++ {

		c, errStr := booksBylanguage(languages[i])
		if errStr != "" {
			http.Error(w, errStr, http.StatusInternalServerError)
			return
		}

		count = append(count, c)
	}

	w.Header().Add("content-type", "application/json")

	// Sending the result as JSON object
	encoder := json.NewEncoder(w)
	err3 := encoder.Encode(count)
	if err3 != nil {
		http.Error(w, "Something went wrong during JSON encoding.", http.StatusInternalServerError)
		return
	}

}

func Readership(w http.ResponseWriter, r *http.Request) {
	pathArr := strings.Split(r.URL.Path, "/")
	subpath := ""

	if len(pathArr) >= 5 {
		subpath = pathArr[4][0:2] //Getting the subpath(language)
	}

	// If language is not specified, the user is instructed
	if subpath == "" {
		instructions := "To use this service, choose a language."
		instructions += "\nExample: localhost:8080/librarystats/v1/readership/no"
		instructions += "\nYou can also add a limit parameter, to limit the amount of countries you want to see."
		instructions += "\nExample: localhost:8080/librarystats/v1/readership/en/?limit=8"

		_, err := fmt.Fprint(w, instructions)
		if err != nil {
			http.Error(w, "Error during generation of response.", http.StatusInternalServerError)
			return
		}

		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, errLmt := strconv.Atoi(limitStr)

	if limitStr == "" {
		limit = 999
	} else if errLmt != nil || limit < 1 {
		http.Error(w, "Invalid parameter.", http.StatusBadRequest)
		return
	}

	l2c_url := "http://129.241.150.113:3000/language2countries/" + subpath

	// Sending request to retrieve countries
	res0, err := http.Get(l2c_url)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Getting the countries
	var data0 []Language2CountriesResponse
	decoder := json.NewDecoder(res0.Body)
	err = decoder.Decode(&data0)
	if err != nil {
		http.Error(w, "Something went wrong during JSON decoding.", http.StatusInternalServerError)
		return
	}

	res0.Body.Close()

	// Retrieving book information of the language
	c, errStr := booksBylanguage(subpath)
	if errStr != "" {
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	cAmt := min(len(data0), limit) // Amount of countries

	var countries []Countries

	// Going through all countries
	for i := 0; i < cAmt; i++ {
		country_url := "http://129.241.150.113:8080/v3.1/name/" + data0[i].Country
		res1, err := http.Get(country_url)
		if err != nil {
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}

		// Getting the population / readership of the country
		var data1 []CountryPopulationResponse
		decoder := json.NewDecoder(res1.Body)
		err = decoder.Decode(&data1)
		if err != nil {
			http.Error(w, "Something went wrong during JSON decoding.", http.StatusInternalServerError)
			return
		}

		res1.Body.Close()

		countries = append(countries, Countries{
			Country:    data0[i].Country,
			Isocode:    data0[i].Isocode,
			Books:      c.Books,
			Authors:    c.Authors,
			Readership: data1[0].Population})
	}

	w.Header().Add("content-type", "application/json")

	// Send result as JSON object
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&countries)
	if err != nil {
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
	}

}

func Status(w http.ResponseWriter, r *http.Request) {
	books_url := "http://129.241.150.113:8000/books/"
	l2c_url := "http://129.241.150.113:3000/language2countries/"
	countries_url := "http://129.241.150.113:8080/"

	res_books, err := http.Get(books_url)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
	defer res_books.Body.Close()

	res_l2c, err := http.Get(l2c_url)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
	defer res_l2c.Body.Close()

	res_countries, err := http.Get(countries_url)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
	}
	defer res_countries.Body.Close()

	currentTime := time.Since(StartTime)
	uptime := currentTime.Seconds()

	status := API_Status{
		GutendexAPI:  res_books.StatusCode,
		LanguageAPI:  res_l2c.StatusCode,
		CountriesAPI: res_countries.StatusCode,
		Version:      "v1",
		Uptime:       uptime}

	w.Header().Add("content-type", "application/json")

	// Send result as JSON object
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&status)
	if err != nil {
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
	}
}
