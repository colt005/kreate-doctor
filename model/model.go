package model

import "encoding/json"

type Doctor []DoctorElement

func UnmarshalDoctor(data []byte) (Doctor, error) {
	var r Doctor
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Doctor) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type DoctorElement struct {
	ID         int64  `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Speciality string `json:"speciality,omitempty"`
	Age        int64  `json:"age,omitempty"`
	Phone      int64  `json:"phone,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

type Patient []PatientElement

func UnmarshalPatient(data []byte) (Patient, error) {
	var r Patient
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Patient) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PatientElement struct {
	ID        *int64  `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	Age       *int64  `json:"age,omitempty"`
	Phone     *int64  `json:"phone,omitempty"`
	DoctorID  *int64  `json:"doctor_id,omitempty"`
	Disease   *string `json:"disease,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
}
