package models

import (
	"encoding/base64"
	"fleets/configs"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Relationship struct {
	ID          int       `json:"id"`
	SubmittedTo int       `json:"submitted_to"`
	SubmittedBy int       `json:"submitted_by"`
	Fullname    string    `json:"fullname"`
	Username    string    `json:"username"`
	ImageUrl    *string   `json:"imageUrl"`
	ImageBase64 *string   `json:"imageBase64"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RelationshipResponseGetById struct {
	ID          int       `json:"id"`
	SubmittedTo int       `json:"submitted_to"`
	SubmittedBy int       `json:"submitted_by"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RelationshipDTO struct {
	SubmittedTo int `json:"submitted_to" validate:"required"`
	SubmittedBy int `json:"submitted_by" validate:"required"`
}

type RelationshipUpdateStatusDTO struct {
	Status string `json:"status" validate:"required,oneof=reject approve"`
}

func GetRelationshipByUserIdControllerQuery(id int) ([]Relationship, error) {
	db := configs.GetDB()

	query := `SELECT 
				a.id,
				a.submitted_to, 
				a.submitted_by, 
				c.username, 
				c.fullname, 
				c.imageUrl,
				a.status,
				a.created_at,
				a.updated_at
			FROM relationships a
			LEFT JOIN users c
			ON a.submitted_by = c.id
			WHERE a.status = 'pending' AND a.submitted_to = ?`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relationships []Relationship

	for rows.Next() {
		var i Relationship
		if err := rows.Scan(&i.ID, &i.SubmittedTo, &i.SubmittedBy, &i.Username, &i.Fullname, &i.ImageUrl, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}

		if i.ImageUrl != nil {
			parts := strings.Split(*i.ImageUrl, ".")
			extPart := parts[1]
			imagePath := filepath.Join("public", "profilepics", *i.ImageUrl)

			imgData, err := os.ReadFile(imagePath)
			if err == nil {
				base64Str := "data:image/" + extPart + ";base64," + base64.StdEncoding.EncodeToString(imgData)
				i.ImageBase64 = &base64Str
			} else {
				fmt.Println("Failed to read image:", err)
			}
		}

		relationships = append(relationships, i)
	}

	return relationships, nil
}

func CheckUserAlreadyHaveApproveRelationshipQuery(user1 int, user2 int) (bool, error) {
	db := configs.GetDB()

	var count int

	query := "SELECT COUNT(id) FROM relationships WHERE (submitted_to = ? AND relationships.status = 'approve') OR (submitted_by = ? AND relationships.status = 'approve') OR (submitted_to = ? AND relationships.status = 'approve') OR (submitted_by = ? AND relationships.status = 'approve')"
	err := db.QueryRow(query, user1, user1, user2, user2).Scan(
		&count,
	)

	if count > 0 {
		return true, nil
	} else {
		return false, err
	}
}

func CreateRelationshipQuery(formData *RelationshipDTO) error {
	db := configs.GetDB()

	query := "INSERT INTO relationships (submitted_to, submitted_by) VALUES (?, ?)"

	_, err := db.Exec(query, formData.SubmittedTo, formData.SubmittedBy)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "relationships.submitted_to") {
			return fmt.Errorf("The relationship application has been registered")
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "submitted_to") {
			return fmt.Errorf("Not found user with id: %d", formData.SubmittedTo)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "submitted_by") {
			return fmt.Errorf("Not found user with id: %d", formData.SubmittedBy)
		}
		return err
	}

	return nil
}

func UpdateRelationshipStatusQuery(data *RelationshipUpdateStatusDTO, id int) error {
	db := configs.GetDB()
	query := "UPDATE relationships SET status = ? where id = ?"

	result, err := db.Exec(query, data.Status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Relationship with ID %d not found", id)
	}

	return nil
}

func GetRelationshipByIdQuery(id int) (RelationshipResponseGetById, error) {
	db := configs.GetDB()
	var data RelationshipResponseGetById

	query := "SELECT id, submitted_to, submitted_by, status, created_at, updated_at from relationships WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.SubmittedTo, &data.SubmittedBy, &data.Status, &data.CreatedAt, &data.UpdatedAt,
	)

	return data, err
}
