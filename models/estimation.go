package models

import (
	"fleets/configs"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Estimation struct {
	ID             int        `json:"id"`
	CarTypeID      string     `json:"carTypeId"`
	BBMID          string     `json:"bbmId"`
	FuelEstimation int        `json:"fuelEstimation"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type EstimationDTO struct {
	CarTypeID      string `json:"carTypeId" validate:"required"`
	BBMID          string `json:"bbmId" validate:"required"`
	FuelEstimation int    `json:"fuelEstimation" validate:"required"`
}

func GetAllEstimationQuery() ([]Estimation, error) {
	db := configs.GetDB()

	query := "SELECT id, ty_car_id, bbm_id, estisi, created_at, updated_at from estimasis"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estimations []Estimation

	for rows.Next() {
		var i Estimation
		if err := rows.Scan(&i.ID, &i.CarTypeID, &i.BBMID, &i.FuelEstimation, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}

		estimations = append(estimations, i)
	}

	return estimations, nil
}

func GetEstimationByIdQuery(id int) (Estimation, error) {
	db := configs.GetDB()
	var data Estimation

	query := "SELECT id, ty_car_id, bbm_id, estisi, created_at, updated_at from estimasis WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.CarTypeID, &data.BBMID, &data.FuelEstimation, &data.CreatedAt, &data.UpdatedAt,
	)

	return data, err
}

func CreateNewEstimationQuery(formData *EstimationDTO) error {
	db := configs.GetDB()

	query := `
		INSERT INTO estimasis (
			ty_car_id, bbm_id, estisi
		) VALUES (?, ?, ?)`

	_, err := db.Exec(query, formData.CarTypeID, formData.BBMID, formData.FuelEstimation)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "bbm_id") {
			return fmt.Errorf("Data BBM: %s belum terdaftar pada master", formData.BBMID)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "ty_car_id") {
			return fmt.Errorf("Tipe kendaraan: %s belum terdaftar pada master", formData.CarTypeID)
		}

		return err
	}

	return nil
}

func UpdateEstimationByIdQuery(data *EstimationDTO, id int) error {
	db := configs.GetDB()
	query := `
			UPDATE estimasis
			SET 
				ty_car_id = ?,
				bbm_id = ?,
				esitis = ?
			WHERE id = ?`

	result, err := db.Exec(query, &data.CarTypeID, &data.BBMID, &data.FuelEstimation, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "bbm_id") {
			return fmt.Errorf("Data BBM: %s belum terdaftar pada master", data.BBMID)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "ty_car_id") {
			return fmt.Errorf("Tipe kendaraan: %s belum terdaftar pada master", data.CarTypeID)
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Estimasi dengan id %d tidak ditemukan", id)
	}

	return nil
}

func DeleteEstimationByIdQuery(id int) error {
	db := configs.GetDB()
	query := "DELETE FROM estimasis where id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 && strings.Contains(mysqlErr.Error(), "Cannot delete or update a parent row: a foreign key constraint fails") {
			return fmt.Errorf("Data estimasi tidak dapat dihapus")
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Estimasi dengan id %d tidak ditemukan", id)
	}

	return nil
}
