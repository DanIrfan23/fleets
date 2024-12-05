package models

import (
	"fleets/configs"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type VehicleType struct {
	ID       string `json:"id"`
	TypeName string `json:"typeName"`
}

type VehicleTypeDTO struct {
	ID       string `json:"id" validate:"required"`
	TypeName string `json:"typeName" validate:"required"`
}

func GetAllVehicleTypeQuery() ([]VehicleType, error) {
	db := configs.GetDB()

	query := "SELECT * from ty_cars"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []VehicleType

	for rows.Next() {
		var i VehicleType
		if err := rows.Scan(&i.ID, &i.TypeName); err != nil {
			return nil, err
		}

		types = append(types, i)
	}

	return types, nil
}

func GetVehicleTypeByIdQuery(id string) (VehicleType, error) {
	db := configs.GetDB()
	var data VehicleType

	query := "SELECT * from ty_cars WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.TypeName,
	)

	return data, err
}

func CreateNewVehicleTypeQuery(formData *VehicleTypeDTO) error {
	db := configs.GetDB()

	query := `
		INSERT INTO ty_cars (
			id, tyCar
		) VALUES (?, ?)`

	_, err := db.Exec(query, formData.ID, formData.TypeName)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("Kode tipe kendaraan: %s sudah terdaftar", formData.ID)
		}

		return err
	}

	return nil
}

func UpdateVehicleTypeByIdQuery(data *VehicleTypeDTO, id string) error {
	db := configs.GetDB()
	query := `
			UPDATE ty_cars
			SET 
				tyCar = ?
			WHERE id = ?`

	result, err := db.Exec(query, &data.TypeName, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("Kode tipe kendaraan: %s sudah terdaftar", data.ID)
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Tipe kendaraan dengan kode %s tidak ditemukan", id)
	}

	return nil
}

func DeleteVehicleTypeByIdQuery(id string) error {
	db := configs.GetDB()
	query := "DELETE FROM ty_cars where id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 && strings.Contains(mysqlErr.Error(), "Cannot delete or update a parent row: a foreign key constraint fails") {
			return fmt.Errorf("Tipe kendaraan tidak dapat dihapus")
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Tipe kendaraan dengan kode %s tidak ditemukan", id)
	}

	return nil
}
