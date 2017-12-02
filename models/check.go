package models

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/kataras/go-errors"
	"log"
)

const CheckTable = "check"

type Check struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	ConfessionName string `json:"confession_name,omitempty"`
	JournalDate    MyTime `json:"journal_date,omitempty"`
	Status         string `json:"status,omitempty"`
}

func (c Check) FindChecksByConfessionNameAndJournalDate(confessionName string, journalDate MyTime) (*[]Check, error) {
	log.Printf("Fetching checks by confession name %s and journal date %s\f", confessionName, journalDate)
	if confessionName == "" || journalDate.IsZero() {
		return nil, errors.New("confession name and journal date must be provided")
	}
	sql := fmt.Sprintf("SELECT `id`, `name`, `confession_name`, `journal_date`, `status` FROM %s WHERE `confession_name`=? and `journal_date`=?", CheckTable)
	rows, err := db.Query(sql, confessionName, journalDate)
	if err != nil {
		log.Printf("Problem querying: %s\n", spew.Sdump(err))
		return nil, err
	}
	defer rows.Close()

	var checks []Check
	for rows.Next() {
		check := Check{}
		err := rows.Scan(&check.Id, &check.Name, &check.ConfessionName, &check.JournalDate, &check.Status)
		if err != nil {
			return nil, err
		}
		checks = append(checks, check)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &checks, nil
}

func (c Check) Save() (*Check, error) {
	log.Printf("Saving Check: %s\n", spew.Sdump(c))
	var sql string
	if c.Id == 0 {
		sql = fmt.Sprintf("INSERT INTO %s (`name`, `confession_name`, `journal_date`, `status` VALUES (?,?,?,?)", CheckTable)
	} else {
		sql = fmt.Sprintf("UPDATE %s SET `name`, `confession_name`, `journal_date`, `status` WHERE `id`=%d", CheckTable, c.Id)
	}
	log.Printf("SQL: %s", sql)
	res, err := db.Exec(sql, c.Name, c.ConfessionName, c.JournalDate, c.Status)
	if err != nil {
		log.Printf("Problem saving check: %s\n", spew.Sdump(err))
		return nil, err
	}
	if c.Id == 0 {
		c.Id, err = res.LastInsertId()
		if err != nil {
			log.Printf("Problem fetching LastInsertId: %s\n", spew.Sdump(err))
			return nil, err
		}
	}
	return &c, nil
}
