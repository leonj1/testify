package models

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"time"
)

const ServicesTable = "services"

type Service struct {
	Id          int64     `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	ShortName   string    `json:"short_name,omitempty"`
	Repo        string    `json:"repo,omitempty"`
	Version     string    `json:"version,omitempty"`
	InstallDate time.Time `json:"install_date,omitempty"`
	IsDocker    bool      `json:"is_docker,omitempty"`
	Host        string    `json:"host,omitempty"`
}

func (s Service) FindByShortNameAndHost(name, host string) (*Service, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if host == "" {
		return nil, errors.New("host cannot be empty")
	}
	sql := fmt.Sprintf("SELECT `id`, `name`, `short_name`, `repo`, `version`, `install_date`, `is_docker`, `host` FROM %s where `host`=? and `short_name`=?", ServicesTable)
	rows, err := db.Query(sql, host, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		service := new(Service)
		err := rows.Scan(&service.Id, &service.Name, &service.ShortName, &service.Repo, &service.Version, &service.InstallDate, &service.IsDocker, &service.Host)
		if err != nil {
			return nil, err
		}
		return service, nil
	}
	return nil, nil
}

func (s Service) Save() (*Service, error) {
	log.Printf("Saving service: %s\n", spew.Sdump(s))
	if s.Host == "" {
		return nil, errors.New("host name cannot be empty")
	}
	var sql string
	if s.Id == 0 {
		service, err := s.FindByShortNameAndHost(s.Host, s.Name)
		if err != nil {
			log.Printf("Problem finding service by name: %s\n", spew.Sdump(err))
			return nil, err
		}
		if service == nil {
			sql = fmt.Sprintf("INSERT INTO %s (`name`, `short_name`, `repo`, `version`, `install_date`, `is_docker`, `host`) VALUES (?,?,?,?,?,?,?)", ServicesTable)
		} else {
			log.Printf("Service already found, therefore updating by that Id")
			sql = fmt.Sprintf("UPDATE %s SET `name`=?, `short_name`=?,`repo`=?, `version`=?,`install_date`=?, `is_docker`=?, `host`=? where id=%d", ServicesTable, service.Id)
		}
	} else {
		log.Printf("Service has an Id, therefore updating by that Id")
		sql = fmt.Sprintf("UPDATE %s SET `name`=?, `short_name`=?,`repo`=?, `version`=?,`install_date`=?, `is_docker`=?, `host`=? where id=%d", ServicesTable, s.Id)
	}
	log.Printf("Executing cmd: %s\n", sql)
	res, err := db.Exec(sql, s.Name, s.ShortName, s.Repo, s.Version, s.InstallDate, s.IsDocker, s.Host)
	if err != nil {
		log.Printf("Problem saving service: %s\n", spew.Sdump(err))
		return nil, err
	}
	if s.Id == 0 {
		s.Id, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func (s Service) FindServicesByHost(host string) ([]Service, error) {
	if host == "" {
		return nil, errors.New("host cannot be empty")
	}
	sql := fmt.Sprintf("SELECT `id`, `name`, `short_name`, `repo`, `version`, `install_date`, `is_docker`, `host` from %s where `host`=?", ServicesTable)
	rows, err := db.Query(sql, host)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	services := make([]Service, 0)
	for rows.Next() {
		service := Service{}
		err := rows.Scan(&service.Id, &service.Name, &service.ShortName, &service.Repo, &service.Version, &service.InstallDate, &service.IsDocker, &service.Host)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func (s Service) FindHostsByServiceShortName(name string) (*[]Hardware, error) {
	log.Printf("Searching for host by service name: %s\n", name)
	if name == "" {
		return nil, errors.New("service name cannot be empty")
	}
	sql := fmt.Sprintf("SELECT `host` from %s where `short_name`=?", ServicesTable)
	rows, err := db.Query(sql, name)
	if err != nil {
		log.Printf("Problem querying to find hosts by service name: %s\n", spew.Sdump(err))
		return nil, err
	}
	defer rows.Close()
	hosts := make([]string, 0)
	for rows.Next() {
		service := Service{}
		err := rows.Scan(&service.Host)
		if err != nil {
			return nil, err
		}
		if service.Host == "" {
			log.Printf("Warning: Service %s has no host set\n", name)
		} else {
			hosts = append(hosts, service.Host)
		}
	}
	if err = rows.Err(); err != nil {
		log.Printf("Problem with service result set: %s\n", spew.Sdump(err))
		return nil, err
	}
	tmpHost := Hardware{}
	allHosts, err := tmpHost.FindByHostsFromArray(hosts)
	if err != nil {
		log.Printf("Problem fetching all hosts when searching by service name: %s\n", spew.Sdump(err))
		return nil, err
	}
	return allHosts, nil
}

func (s Service) DeleteServicesByHost(host string) error {
	log.Printf("Deleting services for host %s\n", host)
	if host == "" {
		return errors.New("host cannot be empty")
	}
	sql := fmt.Sprintf("DELETE from %s where `host`=?", ServicesTable)
	log.Printf("Executing sql: %s\n", sql)
	_, err := db.Exec(sql, host)
	if err != nil {
		return err
	}
	return nil
}
