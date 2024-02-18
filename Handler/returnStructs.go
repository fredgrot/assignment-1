package handler

type BookCount struct {
	Language string
	Books    int
	Authors  int
	Fraction float32
}

type Countries struct {
	Country    string
	Isocode    string
	Books      int
	Authors    int
	Readership int
}

type API_Status struct {
	GutendexAPI  int
	LanguageAPI  int
	CountriesAPI int
	Version      string
	Uptime       float64
}
