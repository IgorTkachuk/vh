package controllers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
	"vh/internal/models"
	"vh/internal/vh"
)

type Server struct {
	core *vh.Vh
}

func NewServer(core *vh.Vh) *Server {
	return &Server{core: core}
}

func (s *Server) Run(port string) error {
	router := mux.NewRouter()

	router.HandleFunc("/upload", s.UploadVideo)

	http.Handle("/", router)

	return http.ListenAndServe(port, router)
}

func (s *Server) UploadVideo(writer http.ResponseWriter, request *http.Request) {
	src, hdr, err := request.FormFile("video")

	if err != nil {
		http.Error(writer, "Wrong request!", http.StatusBadRequest)
		log.Println(err)
		return
	}

	sourceName := request.PostFormValue("source_name")
	srcDate := request.PostFormValue("src_date")
	srcDateTime, err := time.Parse("2006-01-02", srcDate)
	if err != nil {
		logrus.Errorf("Parse date/time error:%v\n", err)
	}
	customer_pn := request.PostFormValue("customer_pn")
	user := request.PostFormValue("user")
	addition := request.PostFormValue("addition")

	s.core.UploadVideo(
		request.Context(),
		models.ImageUnit{
			PayloadName: hdr.Filename,
			Payload:     src,
			PayloadSize: hdr.Size,
		},
		models.StorageObject{
			SourceName: sourceName,
			SrcDate:    srcDateTime,
			CustomerPN: customer_pn,
			User:       user,
			Addition:   addition,
		},
	)

	// Close the obj file and remove temp file
	src.Close()
	request.MultipartForm.RemoveAll()
}
