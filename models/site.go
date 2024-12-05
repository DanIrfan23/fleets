package models

import (
	"fleets/configs"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Site struct {
	ID          string     `json:"id"`
	SiteName    string     `json:"siteName"`
	SiteAddress *string    `json:"siteAddress"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type SiteDTO struct {
	ID          string  `json:"id" validate:"required"`
	SiteName    string  `json:"siteName" validate:"required"`
	SiteAddress *string `json:"siteAddress"`
}

func GetAllSiteQuery() ([]Site, error) {
	db := configs.GetDB()

	query := "SELECT * from sites"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []Site

	for rows.Next() {
		var i Site
		if err := rows.Scan(&i.ID, &i.SiteName, &i.SiteAddress, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}

		sites = append(sites, i)
	}

	return sites, nil
}

func GetSiteByIdQuery(id string) (Site, error) {
	db := configs.GetDB()
	var data Site

	query := "SELECT * from sites WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.SiteName, &data.SiteAddress, &data.CreatedAt, &data.UpdatedAt,
	)

	return data, err
}

func CreateNewSiteQuery(formData *SiteDTO) error {
	db := configs.GetDB()

	query := `
		INSERT INTO sites (
			id, siteName, siteAddress
		) VALUES (?, ?, ?)`

	_, err := db.Exec(query, formData.ID, formData.SiteName, formData.SiteAddress)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("Kode site: %s sudah terdaftar", formData.ID)
		}

		return err
	}

	return nil
}

func UpdateSiteByIdQuery(data *SiteDTO, id string) error {
	db := configs.GetDB()
	query := `
			UPDATE invents
			SET 
				siteName = ?,
				siteAddress = ?
			WHERE id = ?`

	result, err := db.Exec(query, &data.SiteName, &data.SiteAddress, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("Kode site: %s sudah terdaftar", data.ID)
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Site dengan kode %s tidak ditemukan", id)
	}

	return nil
}

func DeleteSiteByIdQuery(id string) error {
	db := configs.GetDB()
	query := "DELETE FROM sites where id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 && strings.Contains(mysqlErr.Error(), "Cannot delete or update a parent row: a foreign key constraint fails") {
			return fmt.Errorf("Site tidak dapat dihapus")
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Site dengan kode %s tidak ditemukan", id)
	}

	return nil
}
