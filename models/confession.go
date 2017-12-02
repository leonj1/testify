package models

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/kataras/go-errors"
	"log"
)

const ConfessionTable = "confession"

type Confession struct {
	Id         int64              `json:"id,omitempty"`
	Name       string             `json:"name,omitempty"`
	EntityType string             `json:"entity_type,omitempty"`
	LastUpdate LastUpdate         `json:"last_update,omitempty"`
	Journal    map[MyTime]Journal `json:"journal,omitempty"`
}

func (c Confession) Save() (*Confession, error) {
	log.Printf("Saving confession: %s\n", spew.Sdump())
	var sql string
	if c.Id == 0 {
		sql = fmt.Sprintf("INSERT INTO %s (`name`, `entity_type`, `last_update`) VALUES (?,?,?)", ConfessionTable)
	} else {
		sql = fmt.Sprintf("UPDATE %s SET `name`=?, `entity_type`=?, `last_update`=? WHERE `id`=%d", ConfessionTable, c.Id)
	}
	log.Printf("sql: %s", sql)
	res, err := db.Exec(sql, c.Name, c.EntityType, c.LastUpdate)
	if err != nil {
		log.Printf("Problem saving service: %s\n", spew.Sdump(err))
		return nil, err
	}
	if c.Id == 0 {
		c.Id, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func (c Confession) FindAll() (*[]Confession, error) {
	sql := fmt.Sprintf("SELECT `id`, `name`, `entity_type`, `last_update` FROM %s", ConfessionTable)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("Problem getting all records: %s\n", spew.Sdump(err))
		return nil, err
	}
	defer rows.Close()
	var confessions []Confession
	for rows.Next() {
		confession := Confession{}
		err := rows.Scan(&confession.Id, &confession.Name, &confession.EntityType, &confession.LastUpdate)
		if err != nil {
			log.Printf("Problem putting result set into struct: %s\n", spew.Sdump(err))
			return nil, err
		}
		j := Journal{}
		journals, err := j.FindByConfessionName(confession.Name)
		if err != nil {
			log.Printf("Problem fetching journam for confession: %s\n", spew.Sdump(err))
			return nil, err
		}
		confessionJournals := make(map[MyTime]Journal)
		for _, journal := range *journals {
			confessionJournals[journal.JournalDate] = journal
		}
		confession.Journal = confessionJournals
		confessions = append(confessions, confession)
	}
	return &confessions, nil
}

func (c Confession) FindByName(confessionName string) (*Confession, error) {
	log.Printf("Find by name: %s\n", confessionName)
	if confessionName == "" {
		log.Printf("Confession name cannot be blank")
		return nil, errors.New("Confession name cannot be blank")
	}

	sql := fmt.Sprintf("SELECT `id`, `name`, `entity_type`, `last_update` FROM %s WHERE `name`=?", ConfessionTable)
	log.Printf("sql: %s\n", sql)
	rows, err := db.Query(sql, confessionName)
	if err != nil {
		log.Printf("Problem fetchin confession by name: %s\n", spew.Sdump(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		confession := new(Confession)
		err := rows.Scan(&confession.Id, &confession.Name, &confession.EntityType, &confession.LastUpdate)
		if err != nil {
			log.Printf("Problem putting result set into struct: %s\n", spew.Sdump(err))
			return nil, err
		}
		j := Journal{}
		journals, err := j.FindByConfessionName(confession.Name)
		if err != nil {
			log.Printf("Problem fetching journam for confession: %s\n", spew.Sdump(err))
			return nil, err
		}
		confessionJournals := make(map[MyTime]Journal)
		for _, journal := range *journals {
			confessionJournals[journal.JournalDate] = journal
		}
		confession.Journal = confessionJournals
		return confession, nil
	}
	return nil, nil
}
