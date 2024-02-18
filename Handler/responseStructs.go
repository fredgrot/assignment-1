package handler

type GutendexResponse struct {
	Count   int    `json:"count"`
	Next    string `json:"next"`
	Results []struct {
		Authors []struct {
			Name string `json:"name"`
		} `json:"authors"`
		Languages []string `json:"languages"`
	} `json:"results"`
}

type GutendexResponseCount struct {
	Count int `json:"count"`
}

type Language2CountriesResponse struct {
	Isocode string `json:"ISO3166_1_Alpha_2"`
	Country string `json:"Official_Name"`
}

type CountryPopulationResponse struct {
	Population int `json:"population"`
}
