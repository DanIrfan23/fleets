package models

import (
	"fleets/configs"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Bbm struct {
	ID        string     `json:"id"`
	BbmDesc   string     `json:"bbmDesc"`
	BbmPrice  float32    `json:"bbmPrice"`
	BbmRemark *string    `json:"bbmRemark"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type BbmDTO struct {
	ID        string  `json:"id" validate:"required"`
	BbmDesc   string  `json:"bbmDesc" validate:"required"`
	BbmPrice  float32 `json:"bbmPrice" validate:"required"`
	BbmRemark *string `json:"bbmRemark"`
}

func GetAllBbmQuery() ([]Bbm, error) {
	db := configs.GetDB()

	query := "SELECT * from bbms"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bbms []Bbm

	for rows.Next() {
		var i Bbm
		if err := rows.Scan(&i.ID, &i.BbmDesc, &i.BbmPrice, &i.BbmRemark, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}

		bbms = append(bbms, i)
	}

	return bbms, nil
}

func GetBbmByIdQuery(id string) (Bbm, error) {
	db := configs.GetDB()
	var data Bbm

	query := "SELECT * from bbms WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.BbmDesc, &data.BbmPrice, &data.BbmRemark, &data.CreatedAt, &data.UpdatedAt,
	)

	return data, err
}

func CreateNewBbmQuery(formData *BbmDTO) error {
	db := configs.GetDB()

	query := `
		INSERT INTO bbms (
			id, bbmDesc, bbmPrice, bbmRemark
		) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, formData.ID, formData.BbmDesc, formData.BbmPrice, formData.BbmRemark)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "bbms.PRIMARY") {
			return fmt.Errorf("Kode BBM: %s sudah terdaftar", formData.ID)
		}

		return err
	}

	return nil
}

func UpdateBbmByIdQuery(data *BbmDTO, id string) error {
	db := configs.GetDB()
	query := `
			UPDATE bbms
			SET 
				bbmDesc = ?,
				bbmPrice = ?,
				bbmRemark = ?
			WHERE id = ?`

	result, err := db.Exec(query, &data.BbmDesc, &data.BbmPrice, &data.BbmRemark, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("Kode BBM: %s sudah terdaftar", data.ID)
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("BBM dengan kode %s tidak ditemukan", id)
	}

	return nil
}

func DeleteBbmByIdQuery(id string) error {
	db := configs.GetDB()
	query := "DELETE FROM bbms where id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 && strings.Contains(mysqlErr.Error(), "Cannot delete or update a parent row: a foreign key constraint fails") {
			return fmt.Errorf("BBM tidak dapat dihapus")
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("BBM dengan kode %s tidak ditemukan", id)
	}

	return nil
}
