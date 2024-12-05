package controllers

import (
	"database/sql"
	"encoding/json"
	"fleets/models"
	"fleets/responses"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func GetAllEstimationController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		estimations, err := models.GetAllEstimationQuery()
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil mengambil keseluruhan data estimasi", estimations)
	}
}

func GetEstimationByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid estimasi ID")
			return
		}

		estimasion, err := models.GetEstimationByIdQuery(id)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Estimasi dengan id = %d tidak ditemukan", id))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil mengambil data estimasi dengan id = %d", id), estimasion)
	}
}

func CreateNewEstimationController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formData models.EstimationDTO

		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(formData)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'CarTypeID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tipe kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'BBMID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Jenis BBM tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'FuelEstimation' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Rasio penggunaan bahan bakar tidak boleh kosong")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.CreateNewEstimationQuery(&formData)
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

		responses.SuccessResponse(w, http.StatusOK, "Data estimasi berhasil di daftarkan", "")
	}
}

func UpdateEstimationByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid estimasi ID")
			return
		}

		var data models.EstimationDTO

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err = validate.Struct(data)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'CarTypeID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tipe kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'BBMID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Jenis BBM tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'FuelEstimation' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Rasio penggunaan bahan bakar tidak boleh kosong")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.UpdateEstimationByIdQuery(&data, id)
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

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil memperbaharui data estimasi dengan id = %d", id), "")
	}
}

func DeleteEstimationByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid estimasi ID")
			return
		}

		err = models.DeleteEstimationByIdQuery(id)
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

		responses.SuccessResponse(w, http.StatusOK, "Berhasil menghapus data estimasi", "")
	}
}
