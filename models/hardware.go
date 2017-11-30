package models

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"time"
)

const HardwareTable = "hardware"

type Hardware struct {
	Id         int64             `json:"id,omitempty"`
	Host       string            `json:"host,omitempty"`
	CreateDate time.Time         `json:"create_date,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
	Services   []Service         `json:"services,omitempty"`
}

func (hardware Hardware) AllHardware() ([]*Hardware, error) {
	sql := fmt.Sprintf("SELECT * from %s", HardwareTable)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	machines := make([]*Hardware, 0)
	for rows.Next() {
		machine := Hardware{}
		err := rows.Scan(&machine.Id, &machine.Host, &machine.CreateDate)
		if err != nil {
			return nil, err
		}
		t := Tag{}
		tags, err := t.FindTagsByHost(machine.Host)
		if err != nil {
			return nil, err
		}
		machineTags := make(map[string]string)
		for _, tag := range tags {
			machineTags[tag.Key] = tag.Value
		}
		machine.Tags = machineTags

		service := Service{}
		allServices, err := service.FindServicesByHost(machine.Host)
		if err != nil {
			log.Printf("Problem finding services by host: %s\n", spew.Sdump(err))
			return nil, err
		}
		machine.Services = allServices

		machines = append(machines, &machine)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return machines, nil
}

func (hardware Hardware) Save() (*Hardware, error) {
	var sql string
	if hardware.Id == 0 {
		hw, err := hardware.FindByHostName(hardware.Host)
		if err != nil {
			return nil, err
		}
		if hw == nil {
			hardware.CreateDate = time.Now()
			hardware.CreateDate.Format(time.RFC3339)
			sql = fmt.Sprintf("INSERT INTO %s (`host`, `create_date`) VALUES (?,?)", HardwareTable)
		} else {
			return nil, errors.New("host already exists")
		}
	} else {
		sql = fmt.Sprintf("UPDATE %s SET `host`=?, `create_date`=? WHERE `id`=%d", HardwareTable, hardware.Id)
	}

	res, err := db.Exec(sql, hardware.Host, hardware.CreateDate)
	if err != nil {
		log.Printf("Problem executing sql cmd: %s", spew.Sdump(err))
		return nil, err
	}

	if hardware.Id == 0 {
		hardware.Id, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}

	if hardware.Tags != nil {
		// Clean slate - delete previous tags and start anew
		t := Tag{}
		err := t.DeleteTagsByHost(hardware.Host)
		if err != nil {
			return nil, err
		}
		for key, value := range hardware.Tags {
			tag := Tag{
				Key:   key,
				Value: value,
				Host:  hardware.Host,
			}
			_, err := tag.Save()
			if err != nil {
				log.Printf("Problem saving tags: %s", spew.Sdump(err))
				return nil, err
			}
		}
	}

	if hardware.Services != nil {
		// Clean slate - delete previous services and start anew
		s := Service{}
		err := s.DeleteServicesByHost(hardware.Host)
		if err != nil {
			log.Printf("Problem deleting services: %s", spew.Sdump(err))
			return nil, err
		}
		for _, svc := range hardware.Services {
			svc.Host = hardware.Host
			_, err := svc.Save()
			if err != nil {
				log.Printf("Problem saving service: %s", spew.Sdump(err))
				return nil, err
			}
		}
	}

	return &hardware, nil
}

func (hardware Hardware) FindByHostName(host string) (*Hardware, error) {
	log.Printf("Finding host by hostname: %s\n", host)
	if host == "" {
		return nil, errors.New("please provide a host")
	}
	sql := fmt.Sprintf("select `id`, `host`, `create_date` from %s where `host`=?", HardwareTable)
	rows, err := db.Query(sql, host)
	if err != nil {
		log.Printf("Problem finding host by name: %s\n", spew.Sdump(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		h := new(Hardware)
		err := rows.Scan(&h.Id, &h.Host, &h.CreateDate)
		if err != nil {
			log.Printf("Problem mapping result set to struct: %s\n", spew.Sdump(err))
			return nil, err
		}

		tag := Tag{}
		allHostTags, err := tag.FindTagsByHost(h.Host)
		if err != nil {
			log.Printf("Problem finding tags by host: %s\n", spew.Sdump(err))
			return nil, err
		}
		tg := make(map[string]string)
		for _, ts := range allHostTags {
			tg[ts.Key] = ts.Value
		}
		h.Tags = tg

		service := Service{}
		allServices, err := service.FindServicesByHost(h.Host)
		if err != nil {
			log.Printf("Problem finding services by host: %s\n", spew.Sdump(err))
			return nil, err
		}
		h.Services = allServices
		return h, nil
	}
	return nil, nil
}

func (hardware Hardware) FindByHostsFromArray(hosts []string) (*[]Hardware, error) {
	log.Printf("Hosts to search for: %s\n", spew.Sdump(hosts))
	if hosts == nil || len(hosts) == 0 {
		return nil, errors.New("please provide a list of hosts")
	}
	var hardwares []Hardware
	for _, host := range hosts {
		hw, err := hardware.FindByHostName(host)
		if err != nil {
			return nil, err
		}
		hardwares = append(hardwares, *hw)
	}
	return &hardwares, nil
}

func (hardware Hardware) FindById(id int64) (*Hardware, error) {
	if id == 0 {
		return nil, errors.New("please provide an id")
	}
	sql := fmt.Sprintf("select `id`, `host`, `create_date` from %s where `id`=?", HardwareTable)
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		h := new(Hardware)
		err := rows.Scan(&h.Id, &h.Host, &h.CreateDate)
		if err != nil {
			return nil, err
		}
		return h, nil
	}
	return nil, nil
}
