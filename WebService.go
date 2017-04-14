package main

import "fmt"
import "io/ioutil"
import "net/http"
import "strconv"


//main service class
type WebService struct {

}

//start and init service
func (service *WebService) Start(listenPort int) error {
	fmt.Println("listening on :" + strconv.Itoa(listenPort))
	http.HandleFunc("/", service.Redirect)
	http.HandleFunc("/index.html", service.ServePage)
	retVal := http.ListenAndServe(":"+strconv.Itoa(listenPort), nil)
	return retVal
}

//redirect all the wrong unplanned queries to index
func (service *WebService) Redirect(responseWriter http.ResponseWriter, request *http.Request) {
	http.Redirect(responseWriter, request, "/index.html", 301)
}

//serve main page request
func (service *WebService) ServePage(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type: text/html", "*")
	content, _ := ioutil.ReadFile("index.html")
	responseWriter.Write(content)
}

