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

func GetRelationshipByUserIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		relationships, err := models.GetRelationshipByUserIdControllerQuery(id)
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Success retrieved all relationships request for user id = %s", idStr), relationships)
	}
}

func CreateRelationshipController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formData models.RelationshipDTO

		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err := validate.Struct(formData)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'SubmittedTo' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Invalid request: Please ensure all required fields are filled out correctly.", "Submitted to can't be empty")
			} else if strings.Contains(errors.Error(), "Error:Field validation for 'SubmittedBy' failed on the 'required' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Invalid request: Please ensure all required fields are filled out correctly.", "Submitted by can't be empty")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Invalid request: Please ensure all required fields are filled out correctly.", "")
			}
			return
		}

		isUserAlreadyHaveRelationship, err := models.CheckUserAlreadyHaveApproveRelationshipQuery(formData.SubmittedTo, formData.SubmittedBy)
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if isUserAlreadyHaveRelationship {
			responses.ErrorResponse(w, http.StatusConflict, "The user has already registered a relationship with another person")
			return
		}

		err = models.CreateRelationshipQuery(&formData)
		if err != nil {
			if strings.Contains(err.Error(), "The relationship application has been registered") {
				responses.ErrorResponse(w, http.StatusConflict, err.Error())
			} else if strings.Contains(err.Error(), "Not found") {
				responses.ErrorResponse(w, http.StatusConflict, err.Error())
			} else {
				log.Println(err)
				responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		responses.SuccessResponse(w, http.StatusOK, "Relationship successfully registered", "")
	}
}

func UpdateStatusRelationshipController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid relationship ID")
			return
		}

		var data models.RelationshipUpdateStatusDTO

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validate := validator.New()
		err = validate.Struct(data)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			if strings.Contains(errors.Error(), "Error:Field validation for 'Status' failed on the 'oneof' tag") {
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Invalid request: Please ensure all required fields are filled out correctly.", "Status value is invalid")
			} else {
				log.Println(errors.Error())
				responses.ErrorValidationResponse(w, http.StatusUnprocessableEntity, "Invalid request: Please ensure all required fields are filled out correctly.", "")
			}
			return
		}

		relationship, err := models.GetRelationshipByIdQuery(id)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Relationship with id = %d not found", id))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		isUserAlreadyHaveRelationship, err := models.CheckUserAlreadyHaveApproveRelationshipQuery(relationship.SubmittedTo, relationship.SubmittedBy)
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if isUserAlreadyHaveRelationship && data.Status == "approve" {
			responses.ErrorResponse(w, http.StatusConflict, "The user has already registered a relationship with another person")
			return
		}

		err = models.UpdateRelationshipStatusQuery(&data, id)
		if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Success update relationship"), "")
	}
}

func GetRelationshipByIdController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			responses.ErrorResponse(w, http.StatusBadRequest, "Invalid relationship ID")
			return
		}

		relationship, err := models.GetRelationshipByIdQuery(id)
		if err == sql.ErrNoRows {
			responses.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Relationship with id = %d not found", id))
			return
		} else if err != nil {
			log.Println(err)
			responses.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		responses.SuccessResponse(w, http.StatusOK, fmt.Sprintf("Success get relationship data with id = %d", id), relationship)
	}
}
