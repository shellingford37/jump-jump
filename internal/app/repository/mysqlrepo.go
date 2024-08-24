package repository

import (
	"database/sql"
	"errors"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
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
