package main

import "fmt"
import "io/ioutil"
import "net/http"
import (
	"strconv"
	"encoding/json"
	"log"
)

const MAX_IMAGE_FILE_SIZE = 10 * 1024 * 1024



//main service class
type WebService struct {
	ImageTransformer ImageTransformer
}

type errorResponse struct {
	Text string		`json:"text"`
}

type successResponse struct {
	Message string
	SourceHistogram Histogram
	TransformedHistogram Histogram
	FilteredHistogram Histogram
}

//start and init service
func (service *WebService) Start(listenPort int) error {
	fmt.Println("listening on :" + strconv.Itoa(listenPort))
	http.HandleFunc("/", service.Redirect)
	http.HandleFunc("/index.html", service.ServeInterface)
	http.HandleFunc("/upload", service.UploadImage)
	http.HandleFunc("/sourceImage", service.ServeSourceImage)
	http.HandleFunc("/transformedImage", service.ServeTransformedImage)
	http.HandleFunc("/filteredImage", service.ServeFilteredImage)
	retVal := http.ListenAndServe(":"+strconv.Itoa(listenPort), nil)
	return retVal
}

//redirect all the wrong unplanned queries to index
func (service *WebService) Redirect(responseWriter http.ResponseWriter, request *http.Request) {
	http.Redirect(responseWriter, request, "/index.html", 301)
}

//serve main page request
func (service *WebService) ServeInterface(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/html")
	content, _ := ioutil.ReadFile("index.html")
	responseWriter.Write(content)
}

//serve main page request
func (service *WebService) ServeSourceImage(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type","image/png")
	err := service.ImageTransformer.DumpSourceImage(responseWriter)

	if err != nil {
		service.writeErrorResponse(responseWriter, err)
	}
}

//serve main page request
func (service *WebService) ServeTransformedImage(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type","image/png")
	err := service.ImageTransformer.DumpTransformedImage(responseWriter)

	if err != nil {
		service.writeErrorResponse(responseWriter, err)
	}
}

//serve main page request
func (service *WebService) ServeFilteredImage(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type","image/png")
	err := service.ImageTransformer.DumpFilteredImage(responseWriter)

	if err != nil {
		service.writeErrorResponse(responseWriter, err)
	}
}

//upload processing image
func (service *WebService) UploadImage(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	request.ParseMultipartForm(MAX_IMAGE_FILE_SIZE)

	maskSizeString := request.FormValue("mask_size")
	maskSize, _ := strconv.Atoi(maskSizeString)

	imageFile, _, err := request.FormFile("source_image")

	if err != nil {
		service.writeErrorResponse(responseWriter, err)
		return
	}

	defer imageFile.Close()

	if err = service.ImageTransformer.LoadSourceImage(imageFile); err != nil {
		service.writeErrorResponse(responseWriter, err)
		return
	}

	if err = service.ImageTransformer.TransformImage(); err != nil {
		service.writeErrorResponse(responseWriter, err)
		return
	}

	if err = service.ImageTransformer.FilterImage(maskSize); err != nil {
		service.writeErrorResponse(responseWriter, err)
		return
	}

	var hist, transformedHist, filteredHist Histogram

	if hist, err = service.ImageTransformer.GetSourceHistogram(); err != nil {
		service.writeErrorResponse(responseWriter, err)
		return
	}
	if transformedHist, err = service.ImageTransformer.GetTransformedHistogram(); err != nil {
		service.writeErrorResponse(responseWriter, err)
		return
	}
	if filteredHist, err = service.ImageTransformer.GetFilteredHistogram(); err != nil {
		service.writeErrorResponse(responseWriter, err)
		return
	}

	response := successResponse{
		Message: "OK",
		SourceHistogram: hist,
		TransformedHistogram: transformedHist,
		FilteredHistogram: filteredHist,
	}

	service.writeSuccessResponse(responseWriter, response)
}

func (service *WebService) writeSuccessResponse(responseWriter http.ResponseWriter, data interface{}) {
	responseJson, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Erorr during writing success respose:", err.Error())
	}
	responseWriter.Header().Set("Content-Type","application/json")
	responseWriter.Write(responseJson)
}

func (service *WebService) writeErrorResponse(responseWriter http.ResponseWriter, err error) {
	response := errorResponse{err.Error()}
	responseWriter.Header().Set("Content-Type","application/json")
	responseWriter.WriteHeader(http.StatusInternalServerError)
	responseJson, _ := json.Marshal(response)
	responseWriter.Write(responseJson)
}
