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

type ServiceItem struct {
	ID        string  `json:"id"`
	ItemName  string  `json:"itemName"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type ServiceItemDTO struct {
	ID       string `json:"id" validate:"required"`
	ItemName string `json:"itemName" validate:"required"`
}

func GetServiceItemNewIdQuery() (string, error) {
	db := configs.GetDB()

	prefix := "I"

	query := "SELECT id FROM itemserv ORDER BY id DESC LIMIT 1"

	var lastId string
	err := db.QueryRow(query).Scan(&lastId)

	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	var sequence int
	if lastId != "" {
		sequence, err = parseSequenceNumberForItem(lastId)
		if err != nil {
			return "", err
		}
	}

	sequence++

	newId := fmt.Sprintf("%s%04d", prefix, sequence)

	return newId, nil
}

func parseSequenceNumberForItem(lastId string) (int, error) {
	if len(lastId) < 10 {
		return 0, fmt.Errorf("invalid ID format: %s", lastId)
	}

	seqPart := lastId[len(lastId)-4:]
	sequence, err := strconv.Atoi(seqPart)
	if err != nil {
		return 0, fmt.Errorf("failed to convert sequence part: %v", err)
	}

	return sequence, nil
}

func GetAllServiceItemQuery() ([]ServiceItem, error) {
	db := configs.GetDB()

	query := "SELECT * from itemserv"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servitems []ServiceItem

	for rows.Next() {
		var i ServiceItem
		if err := rows.Scan(&i.ID, &i.ItemName, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}

		servitems = append(servitems, i)
	}

	return servitems, nil
}

func GetServiceItemByIdQuery(id string) (ServiceItem, error) {
	db := configs.GetDB()
	var data ServiceItem

	query := "SELECT * from itemserv WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.ItemName, &data.CreatedAt, &data.UpdatedAt,
	)

	return data, err
}

func CreateNewServiceItemQuery(formData *ServiceItemDTO, username string) error {
	db := configs.GetDB()

	now := time.Now()

	createdTime := now.Format("2006-01-02 15:04:05")

	createdAt := username + " " + createdTime

	query := `
		INSERT INTO itemserv (
			id, itemName, created_at
		) VALUES (?, ?, ?)`

	_, err := db.Exec(query, formData.ID, formData.ItemName, createdAt)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "id") {
			return fmt.Errorf("ID item service: %s sudah terdaftar", formData.ID)
		}

		return err
	}

	return nil
}

func UpdateServiceItemByIdQuery(data *ServiceItemDTO, id string, username string) error {
	db := configs.GetDB()

	now := time.Now()

	updatedTime := now.Format("2006-01-02 15:04:05")

	updatedAt := username + " " + updatedTime

	query := `
			UPDATE itemserv
			SET 
				itemName = ?,
				updated_at = ?
			WHERE id = ?`

	result, err := db.Exec(query, &data.ItemName, updatedAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Item service dengan id %d tidak ditemukan", id)
	}

	return nil
}

func DeleteServiceItemByIdQuery(id string) error {
	db := configs.GetDB()
	query := "DELETE FROM itemserv where id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 && strings.Contains(mysqlErr.Error(), "Cannot delete or update a parent row: a foreign key constraint fails") {
			return fmt.Errorf("Item service tidak dapat dihapus")
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Item service dengan id %s tidak ditemukan", id)
	}

	return nil
}
