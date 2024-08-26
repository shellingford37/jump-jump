package server

import (
	"encoding/json"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
	"log"
)

func flushDiffToDb() error {
	flushDiffUserToDb()
	flushDiffSortLinksToDb()
	return nil
}

func flushToDb() error {
	flushSortLinksKeyToDb()
	flushUserToDb()
	flushRequestHistoryToDb()
	return nil
}

func flushDiffSortLinksToDb() {
	diffkey := utils.GetDiffShortLinkKey()
	redis := db.GetRedisClient()
	diffRhash, err := redis.HGetAll(diffkey).Result()
	if err != nil {
		log.Printf("[FlushDiffToDbLog] 查找短连接失败，error: %v\n", err)
		return
	}
	shortLinkRepo := repository.GetShortLinkMySqRepo()
	for linkId, _ := range diffRhash {
		key := utils.GetShortLinkKey(linkId)
		existFlag, err := redis.Exists(key).Result()
		if err != nil {
			log.Printf("[FlushDiffToDbLog] 查找短链接失败，error: %v\n", err)
			continue
		}
		if existFlag == 0 {
			err = shortLinkRepo.DeleteByLinkId(linkId)
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 短链%s写入数据库失败，error: %v\n", linkId, err)
				continue
			}
		} else {
			s := &models.ShortLink{}
			rs, err := redis.Get(key).Result()
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 短链%s解析失败，error: %v\n", linkId, err)
				continue
			}
			err = json.Unmarshal([]byte(rs), s)
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 短链%s解析JSON失败，error: %v\n", linkId, err)
				continue
			}
			err = shortLinkRepo.SaveOrUpdate(s)
			if err != nil {
				log.Printf("[FlushToDbLog] 短链%s写入数据库失败，error: %v\n", linkId, err)
			}
		}
		_, err = redis.HDel(diffkey, linkId).Result()
		if err != nil {
			log.Printf("[FlushDiffToDbLog] 短链%s删除diff失败，error: %v\n", linkId, err)
		} else {
			log.Printf("[FlushDiffToDbLog] 短链%s写入数据库成功\n", linkId)
		}
	}

}

func flushDiffUserToDb() {
	diffkey := utils.GetDiffUsersKey()
	redis := db.GetRedisClient()
	diffRhash, err := redis.HGetAll(diffkey).Result()
	if err != nil {
		log.Printf("[FlushDiffToDbLog] 查找用户失败，error: %v\n", err)
		return
	}
	key := utils.GetUserKey()
	userRepo := repository.GetUserMySqRepo()
	for username, _ := range diffRhash {
		exist, err := redis.HExists(key, username).Result()
		if err != nil {
			log.Printf("[FlushDiffToDbLog] 查找用户失败，error: %v\n", err)
			continue
		}
		if exist {
			value, err := redis.HGet(key, username).Result()
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 查找用户失败，error: %v\n", err)
				continue
			}
			u := &models.User2{}
			err = json.Unmarshal([]byte(value), u)
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 用户%s解析失败，error: %v\n", key, err)
				continue
			}
			err = userRepo.SaveOrUpdate(u)
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 用户%s写入数据库失败，error: %v\n", key, err)
			}
		} else {
			err = userRepo.DeleteByUsername(username)
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 用户%s写入数据库失败，error: %v\n", key, err)
			}
		}
		_, err = redis.HDel(diffkey, username).Result()
		if err != nil {
			log.Printf("[FlushDiffToDbLog] 用户%s删除diff失败，error: %v\n", key, err)
		} else {
			log.Printf("[FlushDiffToDbLog] 用户%s写入数据库成功\n", key)
		}
	}
}

func flushRequestHistoryToDb() {
	key := utils.GetShortLinksKey()
	redis := db.GetRedisClient()
	linkIds, err := redis.ZRange(key, 0, -1).Result()
	if err != nil {
		log.Printf("[FlushToDbLog] 查找短链接历史失败，error: %v\n", err)
		return
	}
	for _, id := range linkIds {
		rhKey := utils.GetRequestHistoryKey(id)
		rs, err := redis.ZRangeWithScores(rhKey, 0, -1).Result()
		if err != nil {
			log.Printf("[FlushToDbLog] 查找短链接历史失败，error: %v\n", err)
			continue
		}
		for _, one := range rs {
			rh := &models.RequestHistory{}
			err = json.Unmarshal([]byte(one.Member.(string)), rh)
			if err != nil {
				log.Printf("[FlushToDbLog] 短链接历史%s解析失败，error: %v\n", key, err)
				continue
			}

		}
	}
}

func flushUserToDb() {
	key := utils.GetUserKey()
	redis := db.GetRedisClient()
	rhash, err := redis.HGetAll(key).Result()
	if err != nil {
		log.Printf("[FlushToDbLog] 查找用户失败，error: %v\n", err)
		return
	}
	for key, value := range rhash {
		//log.Printf("Key: %s, Value: %s\n", key, value)
		u := &models.User2{}
		err = json.Unmarshal([]byte(value), u)
		if err != nil {
			log.Printf("[FlushToDbLog] 用户%s解析失败，error: %v\n", key, err)
			continue
		}
		userRepo := repository.GetUserMySqRepo()
		err = userRepo.SaveOrUpdate(u)
		if err != nil {
			log.Printf("[FlushToDbLog] 用户%s写入数据库失败，error: %v\n", key, err)
		}
	}

}

func flushSortLinksKeyToDb() {
	key := utils.GetShortLinksKey()
	redis := db.GetRedisClient()
	linkIds, err := redis.ZRange(key, 0, -1).Result()
	if err != nil {
		log.Printf("[FlushToDbLog] 查找短链接失败，error: %v\n", err)
		return
	}
	for _, id := range linkIds {
		slRepo := repository.GetShortLinkRepo(redis)
		s, err := slRepo.Get(id)
		if err != nil {
			log.Printf("[FlushToDbLog] 查找短链接失败，error: %v\n", err)
			continue
		}
		shortLinkRepo := repository.GetShortLinkMySqRepo()
		err = shortLinkRepo.SaveOrUpdate(s)
		if err != nil {
			log.Printf("[FlushToDbLog] 短链%s写入数据库失败，error: %v\n", s.Id, err)
		}
	}
}
