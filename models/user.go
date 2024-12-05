package models

import "fleets/configs"

type UserAuthorization struct {
	Username    string  `json:"username"`
	Level       string  `json:"level"`
	Status      string  `json:"status"`
	Insert      *string `json:"insert"`
	Delete      *string `json:"delete"`
	Project     *string `json:"project"`
	Master      *string `json:"master"`
	Trans       *string `json:"trans"`
	Report      *string `json:"report"`
	User        *string `json:"user"`
	Monitoring  *string `json:"monitoring"`
	Approval    *string `json:"approval"`
	Maintenance *string `json:"maintenance"`
}

func GetUserAuthorizationByUsernameQuery(username string) (UserAuthorization, error) {
	db := configs.GetDB()
	var data UserAuthorization

	query := "SELECT a.usrid, a.usrlevel, a.usrstatus, a.usrinsert, a.usrdelete, a.project_id, b.m_master, b.m_trans, b.m_report, b.m_user, b.m_monitoring, b.m_approval, b.m_maintenance FROM cltbuser a LEFT JOIN cltbduser b ON a.usrid = b.cltbuser_usrid WHERE usrid =  ?"
	err := db.QueryRow(query, username).Scan(
		&data.Username, &data.Level, &data.Status, &data.Insert, &data.Delete, &data.Project, &data.Master, &data.Trans, &data.Report, &data.User, &data.Monitoring, &data.Approval, &data.Maintenance,
	)

	return data, err
}
