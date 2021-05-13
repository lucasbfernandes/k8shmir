package models

type Request struct {
	Id		   string				`json:"id"`
	RequestURI string				`json:"requestURI"`
	Host       string				`json:"host"`
	Method 	   string				`json:"method"`
	Url 	   string				`json:"url"`
	Headers    map[string][]string	`json:"headers"`
	Body 	   []byte				`json:"body"`
}