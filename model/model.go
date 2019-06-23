package model

import (
	"database/sql"
)

type Tomanrate struct {
	Date     string
	Currency int
	Rate     float64
}

func InsertTomanrate(db *sql.DB, tRate Tomanrate) int64 {
	insertTR, _ := db.Exec("INSERT INTO `tomanrate` (`date`,`id_currency`, `rate`) VALUES ( ?,?,?);", tRate.Date, tRate.Currency, tRate.Rate)
	id, _ := insertTR.LastInsertId()
	return id
}
