package repository

import (
	"database/sql"
	"errors"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"time"
)

type shortLinkMySqlRepository struct {
}

var shortLinkMySqlRepo *shortLinkMySqlRepository

func GetShortLinkMySqRepo() *shortLinkMySqlRepository {
	if shortLinkMySqlRepo == nil {
		shortLinkMySqlRepo = &shortLinkMySqlRepository{}
	}
	return shortLinkMySqlRepo
}

func (r *shortLinkMySqlRepository) DeleteByLinkId(linkId string) error {
	mysqlDB, err := db.GetMySql()
	if err != nil {
		return err
	}
	stmt, err := mysqlDB.Prepare("DELETE FROM `short_link` where `link_id` = ? ")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(linkId)
	if err != nil {
		return err
	}
	return nil
}

func (r *shortLinkMySqlRepository) SaveOrUpdate(link *models.ShortLink) error {
	mysqlDB, err := db.GetMySql()
	if err != nil {
		return err
	}
	//stmt, err := mysqlDB.Prepare("SELECT `id`,`link_id`,`url`,`description`,`is_enable`,`created_by`,`create_time`,`update_time` FROM `short_link` where `id` = ? ")
	stmt, err := mysqlDB.Prepare("select `id` FROM `short_link` where `link_id` = ? ")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var dbId int64
	err = stmt.QueryRow(link.Id).Scan(&dbId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		stmt, err = mysqlDB.Prepare("INSERT INTO `short_link` (`link_id`,`url`,`description`,`is_enable`,`created_by`,`create_time`,`update_time`) VALUES (?,?,?,?,?,?,?) ")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(link.Id, link.Url, link.Description, link.IsEnable, link.CreatedBy, link.CreateTime, link.UpdateTime)
		if err != nil {
			return err
		}
	} else {
		stmt, err = mysqlDB.Prepare("UPDATE `short_link` SET `url` = ?, `description` = ?, `is_enable` = ?, `created_by` = ?, `create_time` = ? , `update_time` = ? WHERE `id` = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(link.Url, link.Description, link.IsEnable, link.CreatedBy, link.CreateTime, link.UpdateTime, dbId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *shortLinkMySqlRepository) ReadAll() ([]*models.ShortLink, error) {
	mysqlDB, err := db.GetMySql()
	if err != nil {
		return nil, err
	}
	stmt, err := mysqlDB.Prepare("SELECT link_id, url, description, is_enable, created_by, create_time, update_time FROM short_link ")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	datas := make([]*models.ShortLink, 0)
	for rows.Next() {
		data := &models.ShortLink{}
		var createdTimeStr, updatedTimeStr string
		err = rows.Scan(&data.Id, &data.Url, &data.Description, &data.IsEnable, &data.CreatedBy, &createdTimeStr, &updatedTimeStr)
		if err != nil {
			continue
		}
		data.CreateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", createdTimeStr, time.Local)
		data.UpdateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", updatedTimeStr, time.Local)
		datas = append(datas, data)
	}
	return datas, nil
}

type userMySqlRepository struct {
}

var userMySqlRepo *userMySqlRepository

func GetUserMySqRepo() *userMySqlRepository {
	if userMySqlRepo == nil {
		userMySqlRepo = &userMySqlRepository{}
	}
	return userMySqlRepo
}

func (r *userMySqlRepository) ReadAll() ([]*models.User2, error) {
	mysqlDB, err := db.GetMySql()
	if err != nil {
		return nil, err
	}
	stmt, err := mysqlDB.Prepare("SELECT username, `role`, password, salt, create_time FROM `user` ")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	datas := make([]*models.User2, 0)
	for rows.Next() {
		data := &models.User2{}
		var timeStr string
		err = rows.Scan(&data.Username, &data.Role, &data.Password, &data.Salt, &timeStr)
		if err != nil {
			continue
		}
		data.CreateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
		datas = append(datas, data)
	}
	return datas, nil
}

func (r *userMySqlRepository) SaveOrUpdate(user *models.User2) error {
	mysqlDB, err := db.GetMySql()
	if err != nil {
		return err
	}
	stmt, err := mysqlDB.Prepare("SELECT `id` FROM `user` where `username` = ? ")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var dbId int64
	err = stmt.QueryRow(user.Username).Scan(&dbId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		stmt, err = mysqlDB.Prepare("INSERT INTO `user` (`username`, `role`, `password`, `salt`, `create_time`)  VALUES (?,?,?,?,?) ")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(user.Username, user.Role, user.Password, user.Salt, user.CreateTime)
		if err != nil {
			return err
		}
	} else {
		stmt, err = mysqlDB.Prepare("UPDATE `user` SET `role` = ?, `password` = ? , `salt` = ? , `create_time` = ? WHERE `id` = ? ")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(user.Role, user.Password, user.Salt, user.CreateTime, dbId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *userMySqlRepository) DeleteByUsername(username string) error {
	mysqlDB, err := db.GetMySql()
	if err != nil {
		return err
	}
	stmt, err := mysqlDB.Prepare("DELETE FROM `user` where `username` = ? ")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}
	return nil
}

type shortLinkHistoryMySqlRepository struct {
}

var shortLinkHistoryMySqlRepo *shortLinkHistoryMySqlRepository

func GetShortLinkHistoryMySqlRepo() *shortLinkHistoryMySqlRepository {
	if shortLinkHistoryMySqlRepo == nil {
		shortLinkHistoryMySqlRepo = &shortLinkHistoryMySqlRepository{}
	}
	return shortLinkHistoryMySqlRepo
}

func (r *shortLinkHistoryMySqlRepository) SaveOrUpdate(linkId string, rh *models.RequestHistory) error {
	mysqlDB, err := db.GetMySql()
	if err != nil {
		return err
	}
	stmt, err := mysqlDB.Prepare("SELECT id FROM `request_history` where link_id =? and bid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var dbId int64
	err = stmt.QueryRow(linkId, rh.Id).Scan(&dbId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		stmt, err = mysqlDB.Prepare("INSERT INTO request_history (link_id, bid, url, ip, ua, `time`) VALUES(?, ?, ?, ?, ?, ?) ")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(linkId, rh.Id, rh.Url, rh.IP, rh.UA, rh.Time)
		if err != nil {
			return err
		}
	} else {
		stmt, err = mysqlDB.Prepare("UPDATE request_history SET url=?, ip=?, ua=?, `time`= ? WHERE id = ? ")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(rh.Url, rh.IP, rh.UA, rh.Time, dbId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *shortLinkHistoryMySqlRepository) FindByDateRange(linkId string, startTime, endTime time.Time) []*models.RequestHistory {
	mysqlDB, err := db.GetMySql()
	if err != nil {
		return nil
	}
	stmt, err := mysqlDB.Prepare("SELECT bid, url, ip, ua, `time` FROM `request_history` where link_id =? and `time` >= ? and `time` <= ? order by `time` desc ")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(linkId, startTime, endTime)
	if err != nil {
		return nil
	}
	rhs := make([]*models.RequestHistory, 0)
	for rows.Next() {
		rh := &models.RequestHistory{}
		var timeBytes string
		err = rows.Scan(&rh.Id, &rh.Url, &rh.IP, &rh.UA, &timeBytes)
		if err != nil {
			continue
		}
		rh.Time, _ = time.ParseInLocation("2006-01-02 15:04:05", timeBytes, time.Local)
		rhs = append(rhs, rh)
	}
	return rhs
}
