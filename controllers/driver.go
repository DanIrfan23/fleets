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

func GetDriverLastId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := models.GetDriverNewIdQuery()
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil mengambil id driver", id)
	}
}

func GetAllDriverController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		drivers, err := models.GetAllDriverQuery()
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil mengambil keseluruhan data driver", drivers)
	}
}

func GetDriverByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		driver, err := models.GetDriverByIdQuery(idStr)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Driver dengan id = %s tidak ditemukan", idStr))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil mengambil data driver dengan id = %s", idStr), driver)
	}
}

func CreateNewDriverController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formData models.DriverDTO

		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(formData)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'ID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "ID driver tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverName' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Nama driver tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverActive' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status driver tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverActive' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status driver tidak terdaftar")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SiteId' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Lokasi site tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverType' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tipe driver tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverType' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tipe driver tidak terdaftar")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Slip' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status kebutuhan slip gaji tidak terdaftar")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.CreateNewDriverQuery(&formData)
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

		responses.SuccessResponse(w, http.StatusOK, "Data driver berhasil di daftarkan", "")
	}
}

func UpdateDriverByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		var data models.DriverDTO

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(data)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'ID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "ID driver tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverName' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Nama driver tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverActive' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status driver tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverActive' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status driver tidak terdaftar")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SiteId' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Lokasi site tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverType' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tipe driver tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'DriverType' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tipe driver tidak terdaftar")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Slip' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status kebutuhan slip gaji tidak terdaftar")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.UpdateDriverByIdQuery(&data, idStr)
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

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil memperbaharui data driver dengan id = %s", idStr), "")
	}
}

func DeleteDriverByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		err := models.DeleteDriverByIdQuery(idStr)
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

		responses.SuccessResponse(w, http.StatusOK, "Berhasil menghapus data driver", "")
	}
}
