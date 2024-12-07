package models

import (
	"database/sql"
	"fleets/configs"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Driver struct {
	ID           string  `json:"id"`
	DriverName   string  `json:"driverName"`
	DriverTelp   *string `json:"driverTelp"`
	Username     *string `json:"username"`
	DriverActive string  `json:"driverActive"`
	SiteId       *string `json:"site_id"`
	SiteName     *string `json:"siteName"`
	Nik          *string `json:"nik"`
	Slip         string  `json:"slip"`
	DriverType   string  `json:"drvtype"`
	DriverCode   *string `json:"drvcode"`
	Delete       string  `json:"delete"`
	CreatedAt    *string `json:"created_at"`
	UpdatedAt    *string `json:"updated_at"`
	DeletedAt    *string `json:"deleted_at"`
}

type DriverDTO struct {
	ID           string  `json:"id" validate:"required"`
	DriverName   string  `json:"driverName" validate:"required"`
	DriverTelp   *string `json:"driverTelp"`
	Username     *string `json:"username"`
	Password     *string `json:"password"`
	DriverActive string  `json:"driverActive" validate:"required,oneof=Y N"`
	SiteId       string  `json:"site_id" validate:"required"`
	Nik          *string `json:"nik"`
	Slip         *string `json:"slip" validate:"oneof=Y N"`
	DriverType   string  `json:"drvtype" validate:"required,oneof=Driver Helper"`
	DriverCode   *string `json:"drvcode"`
}

func GetDriverNewIdQuery() (string, error) {
	db := configs.GetDB()

	prefix := "DRV"
	currentTime := time.Now()
	currentDate := currentTime.Format("060102")

	search := prefix + currentDate + "%"
	query := "SELECT id FROM drivers WHERE id LIKE ? ORDER BY id DESC LIMIT 1"

	var lastId string
	err := db.QueryRow(query, search).Scan(&lastId)

	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	var sequence int
	if lastId != "" {
		sequence, err = parseSequenceNumberForDriver(lastId)
		if err != nil {
			return "", err
		}
	}

	sequence++

	newId := fmt.Sprintf("%s%s%03d", prefix, currentDate, sequence)

	return newId, nil
}

func parseSequenceNumberForDriver(lastId string) (int, error) {
	if len(lastId) < 10 {
		return 0, fmt.Errorf("invalid ID format: %s", lastId)
	}

	seqPart := lastId[len(lastId)-3:]
	sequence, err := strconv.Atoi(seqPart)
	if err != nil {
		return 0, fmt.Errorf("failed to convert sequence part: %v", err)
	}

	return sequence, nil
}

func GetAllDriverQuery() ([]Driver, error) {
	db := configs.GetDB()

	query := "SELECT a.id, a.driverName, a.driverTelp, a.username, a.driverActive, a.site_id, b.siteName, a.nik, a.needslip, a.drvtype, a.drvcode, a.delete, a.created_at, a.updated_at, a.deleted_at from drivers a left join sites b on a.site_id = b.id where a.delete = '0' order by a.id desc"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []Driver

	for rows.Next() {
		var i Driver
		if err := rows.Scan(&i.ID, &i.DriverName, &i.DriverTelp, &i.Username, &i.DriverActive, &i.SiteId, &i.SiteName, &i.Nik, &i.Slip, &i.DriverType, &i.DriverCode, &i.Delete, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt); err != nil {
			return nil, err
		}

		drivers = append(drivers, i)
	}

	return drivers, nil
}

func GetDriverByIdQuery(id string) (Driver, error) {
	db := configs.GetDB()
	var data Driver

	query := "SELECT a.id, a.driverName, a.driverTelp, a.username, a.driverActive, a.site_id, b.siteName, a.nik, a.needslip, a.drvtype, a.drvcode, a.delete, a.created_at, a.updated_at, a.deleted_at from drivers a left join sites b on a.site_id = b.id where a.delete = '0' and a.id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.DriverName, &data.DriverTelp, &data.Username, &data.DriverActive, &data.SiteId, &data.SiteName, &data.Nik, &data.Slip, &data.DriverType, &data.DriverCode, &data.Delete, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
	)

	return data, err
}

func CreateNewDriverQuery(formData *DriverDTO) error {
	db := configs.GetDB()

	query := `
		INSERT INTO itemserv (
			id, driverName, driverTelp, username, password, driverActive, site_id, nik, needslip, drvtype, drvcode
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, formData.ID, formData.DriverName, formData.DriverTelp, formData.Username, formData.Password, formData.DriverActive, formData.SiteId, formData.Nik, formData.Slip, formData.DriverType, formData.DriverCode)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("ID item service: %s sudah terdaftar", formData.ID)
		}

		return err
	}

	return nil
}

func UpdateDriverByIdQuery(data *DriverDTO, id string) error {
	db := configs.GetDB()

	query := `
			UPDATE drivers
			SET 
				driverName = ?,
				driverTelp = ?,
				username = ?,
				password = ?,
				driverActive = ?,
				site_id = ?,
				nik = ?,
				needslip = ?,
				drvtype = ?,
				drvcode = ?
			WHERE id = ?`

	result, err := db.Exec(query, &data.DriverName, &data.DriverTelp, &data.Username, &data.Password, &data.DriverActive, &data.SiteId, &data.Nik, &data.Slip, &data.DriverType, &data.DriverCode, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Driver dengan id %d tidak ditemukan", id)
	}

	return nil
}

func DeleteDriverByIdQuery(id string) error {
	db := configs.GetDB()

	now := time.Now()

	deletedTime := now.Format("2006-01-02 15:04:05")

	query := `
			UPDATE drivers
			SET 
				delete = '1',
				deleted_at = ?,
			WHERE id = ?`

	result, err := db.Exec(query, deletedTime, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Driver dengan id %s tidak ditemukan", id)
	}

	return nil
}
