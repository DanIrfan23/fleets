package controllers

import (
	"database/sql"
	"fleets/models"
	"fleets/responses"
	"fmt"
	"log"
	"net/http"
)

func GetUserLoginDataController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value("username").(string)

		userAuthorization, err := models.GetUserAuthorizationByUsernameQuery(username)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("User dengan username = %s tidak ditemukan", username))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Berhasil mengambil data authorization user dengan username = %s", username), userAuthorization)
	}
}
