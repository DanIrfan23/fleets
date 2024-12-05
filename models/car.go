package models

import (
	"fleets/configs"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Car struct {
	ID              int        `json:"id"`
	Ipolisi         string     `json:"Ipolisi"`
	Ipemilik        string     `json:"Ipemilik"`
	Imesin          *string    `json:"Imesin"`
	Irangka         *string    `json:"Irangka"`
	IBpkb           *string    `json:"IBpkb"`
	IType           *string    `json:"IType"`
	IMerk           *string    `json:"IMerk"`
	TyCarID         string     `json:"ty_car_id"`
	TyCar           string     `json:"tyCar"`
	BbmID           string     `json:"bbm_id"`
	BbmDesc         string     `json:"bbmDesc"`
	SiteID          string     `json:"site_id"`
	SiteName        string     `json:"siteName"`
	Iactive         *string    `json:"Iactive"`
	Itglpembelian   *string    `json:"Itglpembelian"`
	Itglexpkir      *time.Time `json:"Itglexpkir"`
	Itglexpajak     *time.Time `json:"Itglexpajak"`
	Itglexpstnk     *time.Time `json:"Itglexpstnk"`
	Iposisi         *string    `json:"Iposisi"`
	Igpsno          *string    `json:"Igpsno"`
	Icardno         *string    `json:"Icardno"`
	GpsExpDate      *time.Time `json:"gpsExpDate"`
	GpsCardExpDate  *time.Time `json:"gpsCardExpDate"`
	NoFlazz         *string    `json:"noFlazz"`
	CostProjectCpid *string    `json:"costproject_cpid"`
	Cpdesc          *string    `json:"cpdesc"`
	IStatus         *string    `json:"IStatus"`
	ItglStatus      *time.Time `json:"ItglStatus"`
	VehicleId       *string    `json:"vehicle_id"`
	StatusRemark    *string    `json:"statusRemark"`
	Delete          string     `json:"delete"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

type CarDTO struct {
	Ipolisi         string  `json:"Ipolisi" validate:"required"`
	Ipemilik        string  `json:"Ipemilik" validate:"required"`
	Imesin          *string `json:"Imesin"`
	Irangka         *string `json:"Irangka"`
	IBpkb           *string `json:"IBpkb"`
	IType           string  `json:"IType" validate:"required"`
	IMerk           *string `json:"IMerk"`
	TyCarID         string  `json:"ty_car_id" validate:"required"`
	BbmID           string  `json:"bbm_id" validate:"required"`
	SiteID          string  `json:"site_id" validate:"required"`
	Iactive         string  `json:"Iactive" validate:"required,oneof=Y N"`
	Itglpembelian   string  `json:"Itglpembelian" validate:"required,datetime=2006"`
	Itglexpkir      string  `json:"Itglexpkir" validate:"required,datetime=2006-01-02"`
	Itglexpajak     string  `json:"Itglexpajak" validate:"required,datetime=2006-01-02"`
	Itglexpstnk     string  `json:"Itglexpstnk" validate:"required,datetime=2006-01-02"`
	Iposisi         string  `json:"Iposisi" validate:"required"`
	Igpsno          *string `json:"Igpsno"`
	Icardno         *string `json:"Icardno"`
	GpsExpDate      *string `json:"gpsExpDate" validate:"datetime=2006-01-02"`
	GpsCardExpDate  *string `json:"gpsCardExpDate" validate:"datetime=2006-01-02"`
	NoFlazz         *string `json:"noFlazz"`
	CostProjectCpid string  `json:"costproject_cpid" validate:"required"`
	Cpdesc          *string `json:"cpdesc"`
	IStatus         string  `json:"IStatus" validate:"required,oneof=Active Asuransi Bengkel Dijual Perpanjangan Rusak"`
	ItglStatus      *string `json:"ItglStatus" validate:"datetime=2006-01-02"`
	VehicleId       *string `json:"vehicle_id"`
	StatusRemark    *string `json:"statusRemark"`
}

func GetAllCarQuery() ([]Car, error) {
	db := configs.GetDB()

	query := "SELECT a.id, a.Ipolisi, a.Ipemilik, a.Imesin, a.Irangka, a.Ibpkb, a.Itype, a.Imerk, a.ty_car_id, b.tyCar, a.bbm_id, c.bbmDesc, a.Iactive, a.site_id, d.siteName, a.Itglpembelian, a.Itglexpkir, a.Itglexpajak, a.Itglexpstnk, a.Iposisi, a.Igpsno, a.Icardno, a.gpsExpDate, a.gpsCardExpDate, a.noFlazz, a.costproject_cpid, a.cpdesc, a.IStatus, a.ItglStatus, a.statusRemark, a.delete, a.created_at, a.updated_at, a.deleted_at from invents a left join ty_cars b on a.ty_car_id = b.id left join bbms c on a.bbm_id = c.id left join sites d on a.site_id = d.id where a.delete = '0'"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []Car

	for rows.Next() {
		var i Car
		if err := rows.Scan(&i.ID, &i.Ipolisi, &i.Ipemilik, &i.Imesin, &i.Irangka, &i.IBpkb, &i.IType, &i.IMerk, &i.TyCarID, &i.TyCar, &i.BbmID, &i.BbmDesc, &i.Iactive, &i.SiteID, &i.SiteName, &i.Itglpembelian, &i.Itglexpkir, &i.Itglexpajak, &i.Itglexpstnk, &i.Iposisi, &i.Igpsno, &i.Icardno, &i.GpsExpDate, &i.GpsCardExpDate, &i.NoFlazz, &i.CostProjectCpid, &i.Cpdesc, &i.IStatus, &i.ItglStatus, &i.StatusRemark, &i.Delete, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt); err != nil {
			return nil, err
		}

		cars = append(cars, i)
	}

	return cars, nil
}

func GetCarByIdQuery(id int) (Car, error) {
	db := configs.GetDB()
	var data Car

	query := "SELECT a.id, a.Ipolisi, a.Ipemilik, a.Imesin, a.Irangka, a.Ibpkb, a.Itype, a.Imerk, a.ty_car_id, b.tyCar, a.bbm_id, c.bbmDesc, a.Iactive, a.site_id, d.siteName, a.Itglpembelian, a.Itglexpkir, a.Itglexpajak, a.Itglexpstnk, a.Iposisi, a.Igpsno, a.Icardno, a.gpsExpDate, a.gpsCardExpDate, a.noFlazz, a.costproject_cpid, a.cpdesc, a.IStatus, a.ItglStatus, a.statusRemark, a.delete, a.created_at, a.updated_at, a.deleted_at from invents a left join ty_cars b on a.ty_car_id = b.id left join bbms c on a.bbm_id = c.id left join sites d on a.site_id = d.id WHERE a.id = ?"
	err := db.QueryRow(query, id).Scan(
		&data.ID, &data.Ipolisi, &data.Ipemilik, &data.Imesin, &data.Irangka, &data.IBpkb, &data.IType, &data.IMerk, &data.TyCarID, &data.TyCar, &data.BbmID, &data.BbmDesc, &data.Iactive, &data.SiteID, &data.SiteName, &data.Itglpembelian, &data.Itglexpkir, &data.Itglexpajak, &data.Itglexpstnk, &data.Iposisi, &data.Igpsno, &data.Icardno, &data.GpsExpDate, &data.GpsCardExpDate, &data.NoFlazz, &data.CostProjectCpid, &data.Cpdesc, &data.IStatus, &data.ItglStatus, &data.StatusRemark, &data.Delete, &data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
	)

	return data, err
}

func CreateNewCarQuery(formData *CarDTO) error {
	db := configs.GetDB()

	query := `
		INSERT INTO car_table (
			Ipolisi, Ipemilik, Imesin, Irangka, IBpkb, IType, IMerk, TyCarID, BbmID, SiteID, Iactive, 
			Itglpembelian, Itglexpkir, Itglexpajak, Itglexpstnk, Iposisi, Igpsno, Icardno, GpsExpDate, 
			GpsCardExpDate, NoFlazz, CostProjectCpid, Cpdesc, IStatus, ItglStatus, VehicleId, StatusRemark
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, formData.Ipolisi, formData.Ipemilik, formData.Imesin, formData.Irangka, formData.IBpkb, formData.IType, formData.IMerk, formData.TyCarID, formData.BbmID, formData.Iactive, formData.SiteID, formData.Itglpembelian, formData.Itglexpkir, formData.Itglexpajak, formData.Itglexpstnk, formData.Iposisi, formData.Igpsno, formData.Icardno, formData.GpsExpDate, formData.GpsCardExpDate, formData.NoFlazz, formData.CostProjectCpid, formData.Cpdesc, formData.IStatus, formData.ItglStatus, formData.StatusRemark)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "Ipolisi") {
			return fmt.Errorf("No. Polisi: %s sudah terdaftar", formData.Ipolisi)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "bbm_id") {
			return fmt.Errorf("Data BBM: %s belum terdaftar pada master", formData.BbmID)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "ty_car_id") {
			return fmt.Errorf("Data tipe kendaraan: %s belum terdaftar pada master", formData.TyCarID)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "site_id") {
			return fmt.Errorf("Data site: %s belum terdaftar pada master", formData.SiteID)
		}
		return err
	}

	return nil
}

func UpdateCarByIdQuery(data *CarDTO, id int) error {
	db := configs.GetDB()
	query := `
			UPDATE invents
			SET 
				Ipolisi = ?,
				Ipemilik = ?,
				Imesin = ?,
				Irangka = ?,
				IBpkb = ?,
				IType = ?,
				IMerk = ?,
				ty_car_id = ?,
				bbm_id = ?,
				site_id = ?,
				Iactive = ?,
				Itglpembelian = ?,
				Itglexpkir = ?,
				Itglexpajak = ?,
				Itglexpstnk = ?,
				Iposisi = ?,
				Igpsno = ?,
				Icardno = ?,
				gpsExpDate = ?,
				gpsCardExpDate = ?,
				noFlazz = ?,
				costproject_cpid = ?,
				cpdesc = ?,
				IStatus = ?,
				ItglStatus = ?,
				vehicle_id = ?,
				statusRemark = ?
			WHERE id = ?`

	result, err := db.Exec(query, &data.Ipolisi, &data.Ipemilik, &data.Imesin, &data.Irangka, &data.IBpkb, &data.IType, &data.IMerk, &data.TyCarID, &data.BbmID, &data.Iactive, &data.SiteID, &data.Itglpembelian, &data.Itglexpkir, &data.Itglexpajak, &data.Itglexpstnk, &data.Iposisi, &data.Igpsno, &data.Icardno, &data.GpsExpDate, &data.GpsCardExpDate, &data.NoFlazz, &data.CostProjectCpid, &data.Cpdesc, &data.IStatus, &data.ItglStatus, &data.StatusRemark, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Error(), "Ipolisi") {
			return fmt.Errorf("No. Polisi: %s sudah terdaftar", data.Ipolisi)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "bbm_id") {
			return fmt.Errorf("Data BBM: %s belum terdaftar pada master", data.BbmID)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "ty_car_id") {
			return fmt.Errorf("Data tipe kendaraan: %s belum terdaftar pada master", data.TyCarID)
		} else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1452 && strings.Contains(mysqlErr.Error(), "site_id") {
			return fmt.Errorf("Data site: %s belum terdaftar pada master", data.SiteID)
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Mobil dengan id %d tidak ditemukan", id)
	}

	return nil
}

func DeleteCarByIdQuery(id int) error {
	db := configs.GetDB()

	now := time.Now()

	deletedTime := now.Format("2006-01-02 15:04:05")

	query := `
			UPDATE invents
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
		return fmt.Errorf("Mobil dengan id %d tidak ditemukan", id)
	}

	return nil
}
