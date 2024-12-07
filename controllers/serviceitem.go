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

func GetServiceItemLastIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := models.GetServiceItemNewIdQuery()
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil mengambil id item service", id)
	}
}

func GetAllServiceItemController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceitems, err := models.GetAllServiceItemQuery()
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil mengambil keseluruhan data item service", serviceitems)
	}
}

func GetServiceItemByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		site, err := models.GetServiceItemByIdQuery(idStr)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Item service dengan id = %s tidak ditemukan", idStr))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil mengambil data item service dengan id = %s", idStr), site)
	}
}

func CreateNewServiceItemController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value("username").(string)

		var formData models.ServiceItemDTO

		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(formData)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'ID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "ID item service tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'ItemName' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Nama item service tidak boleh kosong")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.CreateNewServiceItemQuery(&formData, username)
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

		responses.SuccessResponse(w, http.StatusOK, "Data item service berhasil di daftarkan", "")
	}
}

func UpdateServiceItemByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value("username").(string)

		vars := mux.Vars(r)
		idStr := vars["id"]

		var data models.ServiceItemDTO

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(data)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'ID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "ID item service tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'ItemName' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Nama item service tidak boleh kosong")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.UpdateServiceItemByIdQuery(&data, idStr, username)
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

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil memperbaharui data item service dengan id = %s", idStr), "")
	}
}

func DeleteServiceItemByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		err := models.DeleteServiceItemByIdQuery(idStr)
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

		responses.SuccessResponse(w, http.StatusOK, "Berhasil menghapus data item service", "")
	}
}
