package models

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/kataras/go-errors"
	"log"
)

const JournalTable = "journal"

type Journal struct {
	Id             int64   `json:"id,omitempty"`
	By             string  `json:"by,omitempty"`
	ConfessionName string  `json:"confession_name,omitempty"`
	Checks         []Check `json:"checks,omitempty"`
	Status         string  `json:"status,omitempty"`
	JournalDate    MyTime  `json:"journal_date,omitempty"`
}

func (j Journal) FindByConfessionName(confessionName string) (*[]Journal, error) {
	log.Printf("Finding Journal for Confession name: %s\n", confessionName)
	if confessionName == "" {
		return nil, errors.New("Confession name cannot be empty")
	}
	sql := fmt.Sprintf("SELECT `id`, `by`, `confession_name`, `status`, `journal_date` FROM %s WHERE `confession_name`=?", JournalTable)
	rows, err := db.Query(sql, confessionName)
	if err != nil {
		log.Printf("Problem querying journal: %s\n", spew.Sdump(err))
		return nil, err
	}
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		journal := Journal{}
		err := rows.Scan(&journal.Id, &journal.By, &journal.ConfessionName, &journal.Status, &journal.JournalDate)
		if err != nil {
			log.Printf("Problem putting result set into struct")
			return nil, err
		}
		check := Check{}
		checks, err := check.FindChecksByConfessionNameAndJournalDate(confessionName, journal.JournalDate)
		if err != nil {
			log.Printf("Problem fetching checks for journal: %s\n", spew.Sdump(err))
			return nil, err
		}
		journal.Checks = *checks
		journals = append(journals, journal)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &journals, nil
}

func (j Journal) Save() (*Journal, error) {
	log.Printf("Saving journal: %s\n", spew.Sdump(j))
	var sql string
	if j.Id == 0 {
		sql = fmt.Sprintf("INSERT INTO %s (`by`, `confession_name`, `status`, `journal_date`) VALUES (?,?,?,?)", JournalTable)
	} else {
		sql = fmt.Sprintf("UPDATE %s SET `by`=?, `confession_name`=?, `status`=?, `journal_date`=? WHERE `id`=?", JournalTable)
	}
	log.Printf("sql: %s", sql)
	res, err := db.Exec(sql, j.By, j.ConfessionName, j.Status, j.JournalDate, j.Id)
	if err != nil {
		log.Printf("Problem saving journal: %s\n", spew.Sdump(err))
		return nil, err
	}
	if j.Id == 0 {
		j.Id, err = res.LastInsertId()
		if err != nil {
			log.Printf("Problem getting the auto-generated id: %s\n", spew.Sdump(err))
			return nil, err
		}
	}
	return &j, nil
}
