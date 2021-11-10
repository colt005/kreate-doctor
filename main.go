package main

import (
	"kreate_doctor/apiauth"
	"kreate_doctor/dbcontroller"
	"kreate_doctor/handler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	_, err := dbcontroller.PreparePSQL()
	if err != nil {
		panic(err)
	}

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/doctors", apiauth.Middleware(http.HandlerFunc(handler.GetDoctors), dbcontroller.DB)).Methods(http.MethodGet)
	api.HandleFunc("/doctor/{id:[0-9]+}/patients", apiauth.Middleware(http.HandlerFunc(handler.GetPatientsByDoctor), dbcontroller.DB)).Methods(http.MethodGet)
	api.HandleFunc("/patient/{id:[0-9]+}", apiauth.Middleware(http.HandlerFunc(handler.GetPatient), dbcontroller.DB)).Methods(http.MethodGet)
	api.HandleFunc("/patient/{id:[0-9]+}", apiauth.Middleware(http.HandlerFunc(handler.UpdatePatient), dbcontroller.DB)).Methods(http.MethodPut)
	api.HandleFunc("/patient/{id:[0-9]+}", apiauth.Middleware(http.HandlerFunc(handler.DeletePatient), dbcontroller.DB)).Methods(http.MethodDelete)
	api.HandleFunc("/doctor/{id:[0-9]+}/patient", apiauth.Middleware(http.HandlerFunc(handler.AddPatient), dbcontroller.DB)).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
