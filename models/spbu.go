package models

import (
	"fleets/configs"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Spbu struct {
	ID          string     `json:"id"`
	SpbuName    string     `json:"spbuName"`
	SpbuAddress *string    `json:"spbuAddress"`
	SpbuCity    *string    `json:"spbuCity"`
	SpbuPhone   *string    `json:"spbuPhone"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type SpbuDTO struct {
	ID          string  `json:"id" validate:"required"`
	SpbuName    string  `json:"spbuName" validate:"required"`
	SpbuAddress *string `json:"spbuAddress"`
	SpbuCity    string  `json:"spbuCity" validate:"required"`
	SpbuPhone   *string `json:"spbuPhone" validate:"omitempty,numeric"`
}

func GetAllSpbuQuery() ([]Spbu, error) {
	db := configs.GetDB()

	query := "SELECT * from spbus"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spbus []Spbu

	for rows.Next() {
		var i Spbu
		if err := rows.Scan(&i.ID, &i.SpbuName, &i.SpbuAddress, &i.SpbuCity, &i.SpbuPhone, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}

		spbus = append(spbus, i)
	}

	return spbus, nil
}

func GetSpbuByIdQuery(id string) (Spbu, error) {
	db := configs.GetDB()
	var data Spbu

	query := "SELECT * from spbus WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.SpbuName, &data.SpbuAddress, &data.SpbuCity, &data.SpbuPhone, &data.CreatedAt, &data.UpdatedAt,
	)

	return data, err
}

func CreateNewSpbuQuery(formData *SpbuDTO) error {
	db := configs.GetDB()

	query := `
		INSERT INTO spbus (
			id, spbuName, spbuAddress, spbuCity, spbuPhone
		) VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(query, formData.ID, formData.SpbuName, formData.SpbuAddress, formData.SpbuCity, formData.SpbuPhone)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("Kode SPBU: %s sudah terdaftar", formData.ID)
		}

		return err
	}

	return nil
}

func UpdateSpbuByIdQuery(data *SpbuDTO, id string) error {
	db := configs.GetDB()
	query := `
			UPDATE spbus
			SET 
				spbuName = ?,
				spbuAddress = ?,
				spbuCity = ?,
				spbuPhone = ?
			WHERE id = ?`

	result, err := db.Exec(query, &data.SpbuName, &data.SpbuAddress, &data.SpbuCity, &data.SpbuPhone, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("Kode SPBU: %s sudah terdaftar", data.ID)
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("SPBU dengan kode %s tidak ditemukan", id)
	}

	return nil
}

func DeleteSpbuByIdQuery(id string) error {
	db := configs.GetDB()
	query := "DELETE FROM spbus where id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 && strings.Contains(mysqlErr.Error(), "Cannot delete or update a parent row: a foreign key constraint fails") {
			return fmt.Errorf("SPBU tidak dapat dihapus")
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("SPBU dengan kode %s tidak ditemukan", id)
	}

	return nil
}
