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

func GetAllSiteController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sites, err := models.GetAllSiteQuery()
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil mengambil keseluruhan data site", sites)
	}
}

func GetSiteByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		site, err := models.GetSiteByIdQuery(idStr)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Site dengan kode = %s tidak ditemukan", idStr))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil mengambil data site dengan kode = %s", idStr), site)
	}
}

func CreateNewSiteController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formData models.SiteDTO

		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(formData)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'ID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Kode site tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SiteName' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Nama site tidak boleh kosong")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.CreateNewSiteQuery(&formData)
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

		responses.SuccessResponse(w, http.StatusOK, "Data mobil berhasil di daftarkan", "")
	}
}

func UpdateSiteByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		var data models.SiteDTO

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(data)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'ID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Kode site tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SiteName' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Nama site tidak boleh kosong")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.UpdateSiteByIdQuery(&data, idStr)
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

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil memperbaharui data site dengan kode = %s", idStr), "")
	}
}

func DeleteSiteByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		err := models.DeleteSiteByIdQuery(idStr)
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

		responses.SuccessResponse(w, http.StatusOK, "Berhasil menghapus data site", "")
	}
}
