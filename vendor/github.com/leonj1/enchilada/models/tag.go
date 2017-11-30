package models

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
)

const TagTable = "tags"

type Tag struct {
	Id    int64  `json:"id,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	Host  string `json:"host,omit"`
}

func (t Tag) FindByNameAndHost(host, tagName string) (*Tag, error) {
	if tagName == "" {
		return nil, errors.New("tag name cannot be empty")
	}
	sql := fmt.Sprintf("SELECT `id`, `key`, `value`, `host` FROM %s where `host`=? and `key`=?", TagTable)
	rows, err := db.Query(sql, t.Host, t.Key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tag := new(Tag)
		err := rows.Scan(&tag.Id, &tag.Key, &tag.Value, &tag.Host)
		if err != nil {
			return nil, err
		}
		return tag, nil
	}
	return nil, nil
}

func (t Tag) Save() (*Tag, error) {
	log.Printf("Saving tag: %s\n", spew.Sdump(t))
	if t.Host == "" {
		return nil, errors.New("host cannot be empty")
	}
	var sql string
	if t.Id == 0 {
		tag, err := t.FindByNameAndHost(t.Host, t.Key)
		if err != nil {
			log.Printf("Problem finding tag by name: %s\n", spew.Sdump(err))
			return nil, err
		}
		if tag == nil {
			sql = fmt.Sprintf("INSERT INTO %s (`key`, `value`, `host`) VALUES (?,?,?)", TagTable)
		} else {
			sql = fmt.Sprintf("UPDATE %s SET `key`=?, `value`=?, `host`=? where id=%d", TagTable, t.Id)
		}
	} else {
		sql = fmt.Sprintf("UPDATE %s SET `key`=?, `value`=?, `host`=? where id=%d", TagTable, t.Id)
	}
	log.Printf("Executing sql: %s\n", sql)
	res, err := db.Exec(sql, t.Key, t.Value, t.Host)
	if err != nil {
		return nil, err
	}
	if t.Id == 0 {
		t.Id, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func (t Tag) FindTagsByHost(host string) ([]*Tag, error) {
	if host == "" {
		return nil, errors.New("host cannot be empty")
	}
	sql := fmt.Sprintf("SELECT `id`, `key`, `value`, `host` from %s where `host`=?", TagTable)
	rows, err := db.Query(sql, host)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]*Tag, 0)
	for rows.Next() {
		tag := new(Tag)
		err := rows.Scan(&tag.Id, &tag.Key, &tag.Value, &tag.Host)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (t Tag) FindHostsByTag(tag string) (*[]Hardware, error) {
	log.Printf("Searching for host by tag: %s\n", tag)
	if tag == "" {
		return nil, errors.New("tag cannot be empty")
	}
	sql := fmt.Sprintf("SELECT `host` from %s where `key`=?", TagTable)
	rows, err := db.Query(sql, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	hosts := make([]string, 0)
	for rows.Next() {
		tag := new(Tag)
		err := rows.Scan(&tag.Host)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, tag.Host)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	tmpHost := Hardware{}
	allHosts, err := tmpHost.FindByHostsFromArray(hosts)
	if err != nil {
		return nil, err
	}
	return allHosts, nil
}

func (t Tag) DeleteTagsByHost(host string) error {
	log.Printf("Deleting tags for host %s\n", host)
	if host == "" {
		return errors.New("host cannot be empty")
	}
	sql := fmt.Sprintf("DELETE from %s where `host`=?", TagTable)
	log.Printf("Executing sql: %s\n", sql)
	_, err := db.Exec(sql, host)
	if err != nil {
		return err
	}
	return nil
}
