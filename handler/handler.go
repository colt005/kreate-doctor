package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kreate_doctor/dbcontroller"
	"kreate_doctor/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetDoctors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var documents model.Doctor
	rows, err := dbcontroller.DB.Query("select * from doctor")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r model.DoctorElement
		err = rows.Scan(&r.ID, &r.Name, &r.Speciality, &r.Age, &r.Phone, &r.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}

		documents = append(documents, r)
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(documents) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "No Doctors Found"}`))
	} else {

		w.WriteHeader(http.StatusOK)
		//write as json
		bytes, err := documents.Marshal()
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(bytes))
	}

}

func GetPatientsByDoctor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	fmt.Println(vars)
	doctor_id := vars["id"]
	var documents model.Patient
	rows, err := dbcontroller.DB.Query("select * from patient where doctor_id = $1", doctor_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r model.PatientElement
		err = rows.Scan(&r.ID, &r.Name, &r.Age, &r.Phone, &r.DoctorID, &r.Disease, &r.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}

		documents = append(documents, r)
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(documents) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "No Patients Found"}`))
	} else {

		w.WriteHeader(http.StatusOK)
		//write as json
		bytes, err := documents.Marshal()
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(bytes))
	}

}

func GetPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	patient_id := vars["id"]
	var patient model.PatientElement
	err := dbcontroller.DB.QueryRow("select * from patient where id = $1", patient_id).Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Phone, &patient.DoctorID, &patient.Disease, &patient.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {

			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "No Patient found with ID  : ` + patient_id + `"}`))
			return
		}
		log.Fatal(err)
	} else {

		w.WriteHeader(http.StatusOK)
		//write as json
		bytes, err := json.Marshal(patient)
		if err != nil {

			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Not Found"}`))
			return
		}
		w.Write([]byte(bytes))
	}

}

func UpdatePatient(w http.ResponseWriter, r *http.Request) {
	//PUT
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)
	patient_id := vars["id"]
	if err != nil {
		log.Fatal(err)
	}

	var patient model.PatientElement
	json.Unmarshal(body, &patient)
	fmt.Println(patient_id)
	fmt.Println(strconv.FormatInt(*patient.ID, 10))

	if patient_id != strconv.FormatInt(*patient.ID, 10) {

		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(`{"message": "Patient ID mismatch"}`))
		return
	}

	_, err = dbcontroller.DB.Exec("UPDATE patient SET name=$1, age=$2, phone=$3, disease=$4,doctor_id=$5 where id=$6", patient.Name, patient.Age, patient.Phone, patient.Disease, patient.DoctorID, patient.ID)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(`{"message": "Failed to Update Patient", "reason" : "` + err.Error() + `"}`))
	} else {

		w.WriteHeader(http.StatusOK)
		//write as json
		w.Write([]byte(`{"message":"updated"}`))
	}

}

func DeletePatient(w http.ResponseWriter, r *http.Request) {
	//DELETE
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	patient_id := vars["id"]

	_, err := dbcontroller.DB.Exec("DELETE FROM patient where id = $1", patient_id)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(`{"message": "Failed to Delete Patient", "reason" : "` + err.Error() + `"}`))
	} else {

		w.WriteHeader(http.StatusOK)
		//write as json
		w.Write([]byte(`{"message":"Deleted Patient"}`))
	}

}

func AddPatient(w http.ResponseWriter, r *http.Request) {
	//POST

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	doctor_id := vars["id"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var patient model.PatientElement
	json.Unmarshal(body, &patient)

	var addedPatient model.PatientElement

	row, err := dbcontroller.DB.Query("INSERT INTO patient(name,age,phone,doctor_id,disease) values($1,$2,$3,$4,$5) returning *", patient.Name, patient.Age, patient.Phone, doctor_id, patient.Disease)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(`{"message": "Failed to Add Patient", "reason" : "` + err.Error() + `"}`))
		return
	} else {

		for row.Next() {
			err = row.Scan(&addedPatient.ID, &addedPatient.Name, &addedPatient.Age, &addedPatient.Phone, &addedPatient.DoctorID, &addedPatient.Disease, &addedPatient.CreatedAt)
			if err != nil {
				log.Fatal(err)
			}
		}

		w.WriteHeader(http.StatusOK)
		//write as json
		bytes, err := json.Marshal(addedPatient)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte(`{"message": "Failed to Add Patient", "reason" : "` + err.Error() + `"}`))
		}
		w.Write([]byte(bytes))
	}

}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}
