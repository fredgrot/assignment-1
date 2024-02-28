# Library stats

This API offers 3 services: Bookcount per language, potential readership per language, and diagnostic check.


## Instructions


### Bookcount

URL: https://assignment1-5jce.onrender.com/librarystats/v1/bookcount

This endpoint will return the amount of books written in one or more specified languages.
Here you must include a language parameter. Languages are written in their 2-letter ISO codes, and seperated
by commas. Include at least one language.
Example: https://assignment1-5jce.onrender.com/librarystats/v1/bookcount?language=no,fi


### Potential readership

URL: https://assignment1-5jce.onrender.com/librarystats/v1/readership

This endpoint will return the potential readership for each country that speaks this language.
Here you must specify the language, written in its 2-letter ISO code. Optionally, you can include a limit parameter, that will limit the amount of countries returned.
Example: https://assignment1-5jce.onrender.com/librarystats/v1/bookcount/en?limit=4


### Diagnostic check

URL: https://assignment1-5jce.onrender.com/librarystats/v1/status

This endpoint will return the status of all third party services.