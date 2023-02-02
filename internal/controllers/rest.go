package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
	"vh/internal/models"
	"vh/internal/vh"
	"vh/package/jwt"
)

type Server struct {
	core      *vh.Vh
	jwtHelper jwt.Helper
}

func NewServer(core *vh.Vh, jwtHelper jwt.Helper) *Server {
	return &Server{core: core, jwtHelper: jwtHelper}
}

func (s *Server) Run(port string) error {
	router := mux.NewRouter()

	router.HandleFunc("/new", s.newObject)
	router.HandleFunc("/getobjbypn", jwt.Middleware(s.getObjectByBillingPn)).Methods("GET")
	router.HandleFunc("/content/{billing_pn}/{id}", s.getContent).Methods("GET")
	router.HandleFunc("/auth", s.auth).Methods("POST", "PUT")

	http.Handle("/", router)

	return http.ListenAndServe(port, router)
}

func (s *Server) newObject(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	src, hdr, err := request.FormFile("file")

	if err != nil {
		http.Error(writer, "Wrong request!", http.StatusBadRequest)
		log.Println(err)
		return
	}

	origDate := request.PostFormValue("orig_date")
	origDateTime, err := time.Parse("2006-01-02", origDate)
	if err != nil {
		logrus.Errorf("Parse date/time error:%v\n", err)
	}

	billingPn := request.PostFormValue("billing_pn")
	userName := request.PostFormValue("user_name")
	notes := request.PostFormValue("notes")

	err = s.core.UploadObject(
		request.Context(),
		models.StorageObjectUnit{
			PayloadName: hdr.Filename,
			Payload:     src,
			PayloadSize: hdr.Size,
		},
		models.StorageObjectMeta{
			OrigName:  hdr.Filename,
			OrigDate:  origDateTime,
			BillingPn: billingPn,
			UserName:  userName,
			Notes:     notes,
		},
	)

	if err != nil {
		http.Error(writer, "Create object error", http.StatusInternalServerError)
		logrus.Error("Error occurred when create object")
	}

	// Close the obj file and remove temp file
	src.Close()
	request.MultipartForm.RemoveAll()
}

func (s *Server) getObjectByBillingPn(w http.ResponseWriter, r *http.Request) {
	billingPn := r.URL.Query().Get("billing_pn")

	objList, err := s.core.GetObjectByBillingPn(billingPn)
	if err != nil {
		http.Error(w, "Can't retrieve object list for given billing personal number", http.StatusInternalServerError)
		return
	}

	objListByte, err := json.Marshal(objList)
	if err != nil {
		http.Error(w, "Can't serialize object list for given billing personal number", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(objListByte)
}

func (s *Server) getContent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	billingPn := params["billing_pn"]
	id := params["id"]

	objContent, err := s.core.GetContent(
		r.Context(),
		fmt.Sprintf("%s/%s", billingPn, id),
	)
	if err != nil {
		http.Error(w, "Error when retrieve object content from storage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, objContent.PayloadName, time.Now(), objContent.Payload)
	//if _, err := io.Copy(w, objContent.Payload); err != nil {
	//	logrus.Errorf("Cant send object: %v\n", err)
	//	http.Error(w, "Can`t download object!", http.StatusInternalServerError)
	//}
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var tokens []byte
	var err error

	switch r.Method {
	case http.MethodPost:
		var userDto models.UserDto
		if err = json.NewDecoder(r.Body).Decode(&userDto); err != nil {
			http.Error(w, "Authorization has been refused for those credentials", http.StatusUnauthorized)
		}

		tokens, err = s.jwtHelper.GenerateAccessToken(userDto)
		if err != nil {
			http.Error(w, "Error due generate tokens", http.StatusUnauthorized)
		}
	case http.MethodPut:
		var rt jwt.RT
		if err = json.NewDecoder(r.Body).Decode(&rt); err != nil {
			http.Error(w, "Error during decoding received refresh token", http.StatusUnauthorized)
		}
		tokens, err = s.jwtHelper.UpdateRefreshToken(rt)
		if err != nil {
			http.Error(w, "Error due generate tokens", http.StatusUnauthorized)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tokens)
}
