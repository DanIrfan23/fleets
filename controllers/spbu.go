package controllers

import (
	"database/sql"
	"encoding/json"
	"fleets/models"
	"fleets/responses"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func GetAllSpbuController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spbus, err := models.GetAllSpbuQuery()
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil mengambil keseluruhan data spbu", spbus)
	}
}

func GetSpbuByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		spbu, err := models.GetSpbuByIdQuery(idStr)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("SPBU dengan kode = %s tidak ditemukan", idStr))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil mengambil data spbu dengan kode = %s", idStr), spbu)
	}
}

func CreateNewSpbuController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formData models.SpbuDTO

		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(formData)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'ID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Kode SPBU tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SpbuName' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Nama SPBU tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SpbuCity' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Lokasi kota SPBU tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SpbuPhone' failed on the 'numeric' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "No. Telp tidak sesuai format")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.CreateNewSpbuQuery(&formData)
		if err != nil {
			if strings.Contains(err.Error(), "sudah terdaftar") {
				responses.ErrorResponse(w, http.StatusConflict, err.Error())
			} else if strings.Contains(err.Error(), "belum terdaftar pada master") {
				responses.ErrorResponse(w, http.StatusConflict, err.Error())
			} else {
				log.Println(err)
				responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Data SPBU berhasil di daftarkan", "")
	}
}

func UpdateSpbuByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		var data models.SpbuDTO

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(data)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'ID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Kode SPBU tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SpbuName' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Nama SPBU tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SpbuCity' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Lokasi kota SPBU tidak boleh kosong")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.UpdateSpbuByIdQuery(&data, idStr)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				responses.ErrorResponse(w, http.StatusNotFound, err.Error())
			} else if strings.Contains(err.Error(), "sudah terdaftar") {
				responses.ErrorResponse(w, http.StatusConflict, err.Error())
			} else if strings.Contains(err.Error(), "belum terdaftar pada master") {
				responses.ErrorResponse(w, http.StatusConflict, err.Error())
			} else {
				log.Println(err)
				responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil memperbaharui data SPBU dengan kode = %s", idStr), "")
	}
}

func DeleteSpbuByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		err := models.DeleteSpbuByIdQuery(idStr)
		if err != nil {
			if strings.Contains(err.Error(), "tidak ditemukan") {
				responses.ErrorResponse(w, http.StatusNotFound, err.Error())
			} else if strings.Contains(err.Error(), "tidak dapat dihapus") {
				responses.ErrorResponse(w, http.StatusConflict, err.Error())
			} else {
				log.Println(err)
				responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}

			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil menghapus data SPBU", "")
	}
}
