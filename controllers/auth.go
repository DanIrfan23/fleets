package controllers

import (
	"encoding/json"
	"fleets/models"
	"fleets/responses"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func LoginController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var acc models.LoginDTO

		if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(acc)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'Username' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Username tidak boleh kosong")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'Password' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Password tidak boleh kosong")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Validation Error", "Periksa kembali keseluruhan input")
			}
			return
		}

		username, err := models.CheckAccountQuery(&acc)
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		tokenReturn, err := models.GenerateTokenQuery(username)
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusConflict, "Failed to generate token")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Login berhasil", tokenReturn)
	}
}
