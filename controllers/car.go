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

func GetAllCarController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cars, err := models.GetAllCarQuery()
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Berhasil mengambil keseluruhan data mobil", cars)
	}
}

func GetCarByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid car ID")
			return
		}

		car, err := models.GetCarByIdQuery(id)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Mobil dengan id = %d tidak ditemukan", id))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil mengambil data mobil dengan id = %d", id), car)
	}
}

func CreateNewCarController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formData models.CarDTO

		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(formData)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'Ipolisi' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "No. Polisi tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Ipemilik' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Data pemilik tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'IType' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Jenis mobil tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'TyCarID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tipe mobil tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'BbmID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Jenis BBM tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SiteID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Data site tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Iactive' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status keaktifan kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglpembelian' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tahun pembelian kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpkir' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tanggal kadaluwarsa KIR tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpajak' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tanggal kadaluwarsa pajak tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpstnk' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tanggal kadaluwarsa STNK tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Iposisi' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Posisi penggunaan kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'CostProjectCpid' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Data project kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'IStatus' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglpembelian' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tahun pembelian mobil salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpkir' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa KIR salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpajak' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa pajak salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpstnk' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa SNTK salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'GpsExpDate' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa GPS salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'GpsCardExpDate' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa kartu GPS salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'ItglStatus' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal perubahan status salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Iactive' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format status keaktifan tidak ditemukan")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'IStatus' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format status tidak ditemukan")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.CreateNewCarQuery(&formData)
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

func UpdateCarByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid car ID")
			return
		}

		var data models.CarDTO

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err = validate.Struct(data)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'Ipolisi' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "No. Polisi tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Ipemilik' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Data pemilik tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'IType' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Jenis mobil tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'TyCarID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tipe mobil tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'BbmID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Jenis BBM tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SiteID' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Data site tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Iactive' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status keaktifan kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglpembelian' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tahun pembelian kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpkir' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tanggal kadaluwarsa KIR tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpajak' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tanggal kadaluwarsa pajak tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpstnk' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Tanggal kadaluwarsa STNK tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Iposisi' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Posisi penggunaan kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'CostProjectCpid' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Data project kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'IStatus' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Status kendaraan tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglpembelian' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tahun pembelian mobil salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpkir' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa KIR salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpajak' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa pajak salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Itglexpstnk' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa SNTK salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'GpsExpDate' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa GPS salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'GpsCardExpDate' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal kadaluwarsa kartu GPS salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'ItglStatus' failed on the 'datetime' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format tanggal perubahan status salah")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Iactive' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format status keaktifan tidak ditemukan")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'IStatus' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Format status tidak ditemukan")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali semua isian")
			}
			return
		}

		err = models.UpdateCarByIdQuery(&data, id)
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

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil memperbaharui data mobil dengan id = %d", id), "")
	}
}

func DeleteCarByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid car ID")
			return
		}

		err = models.DeleteCarByIdQuery(id)
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
