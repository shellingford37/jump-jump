package server

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/shellingford37/jump-jump/internal/app/db"
	"github.com/shellingford37/jump-jump/internal/app/models"
	"github.com/shellingford37/jump-jump/internal/app/repository"
	"github.com/shellingford37/jump-jump/internal/app/utils"
	"log"
	"time"
)

func clearRedis() {
	redisClient := db.GetRedisClient()
	key := utils.GetActiveLinkKey()
	redisClient.Del(key)
}

func flushToRedis() error {
	log.Printf("[FlushToRedisLog] 从数据库写入redis开始...")
	err := flushUserToRedis()
	if err != nil {
		return err
	}
	err = flushShortLinkToRedis()
	if err != nil {
		return err
	}
	log.Printf("[FlushToRedisLog] 从数据库写入redis结束...")
	return nil
}

func flushDiffToDb() error {
	flushDiffUserToDb()
	flushDiffSortLinksToDb()
	flushDiffSortLinksKeyToDb()
	return nil
}

func flushToDb() error {
	log.Printf("[FlushToDbLog] 从redis写入数据库开始...")
	flushSortLinksKeyToDb()
	flushUserToDb()
	flushRequestHistoryToDb()
	log.Printf("[FlushToDbLog] 从redis写入数据库结束...")
	return nil
}

func flushShortLinkToRedis() error {
	redisClient := db.GetRedisClient()
	shortLinkMySqlRepo := repository.GetShortLinkMySqRepo()
	linkList, err := shortLinkMySqlRepo.ReadAll()
	if err != nil {
		log.Printf("[FlushToRedisLog] 查找短连接失败，error: %v\n", err)
		return err
	}
	for _, link := range linkList {
		jsonStr, err := json.Marshal(link)
		if err != nil {
			log.Printf("[FlushToRedisLog] 短连接序列化失败，error: %v\n", err)
			return err
		}
		redisClient.Set(utils.GetShortLinkKey(link.Id), jsonStr, 0)
		// 保存用户的短链接记录，保存到创建者及全局
		record := redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: link.Id,
		}
		redisClient.ZAdd(utils.GetUserShortLinksKey(link.CreatedBy), record)
		redisClient.ZAdd(utils.GetShortLinksKey(), record)
	}
	return nil
}

func flushUserToRedis() error {
	redisClient := db.GetRedisClient()
	userMySqlRepo := repository.GetUserMySqRepo()
	userList, err := userMySqlRepo.ReadAll()
	if err != nil {
		log.Printf("[FlushToRedisLog] 查找用户失败，error: %v\n", err)
		return err
	}
	key := utils.GetUserKey()
	for _, user := range userList {
		jsonStr, err := json.Marshal(user)
		if err != nil {
			log.Printf("[FlushToRedisLog] 用户序列化失败，error: %v\n", err)
			return err
		}
		redisClient.HSet(key, user.Username, jsonStr)
	}
	return nil
}

func flushDiffSortLinksToDb() {
	diffkey := utils.GetDiffShortLinkKey()
	redisClient := db.GetRedisClient()
	diffRhash, err := redisClient.HGetAll(diffkey).Result()
	if err != nil {
		log.Printf("[FlushDiffToDbLog] 查找短连接失败，error: %v\n", err)
		return
	}
	shortLinkRepo := repository.GetShortLinkMySqRepo()
	for linkId, _ := range diffRhash {
		key := utils.GetShortLinkKey(linkId)
		existFlag, err := redisClient.Exists(key).Result()
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
			rs, err := redisClient.Get(key).Result()
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
		_, err = redisClient.HDel(diffkey, linkId).Result()
		if err != nil {
			log.Printf("[FlushDiffToDbLog] 短链%s删除diff失败，error: %v\n", linkId, err)
		} else {
			log.Printf("[FlushDiffToDbLog] 短链%s写入数据库成功\n", linkId)
		}
	}

}

func flushDiffUserToDb() {
	diffkey := utils.GetDiffUsersKey()
	redisClient := db.GetRedisClient()
	diffRhash, err := redisClient.HGetAll(diffkey).Result()
	if err != nil {
		log.Printf("[FlushDiffToDbLog] 查找用户失败，error: %v\n", err)
		return
	}
	key := utils.GetUserKey()
	userRepo := repository.GetUserMySqRepo()
	for username, _ := range diffRhash {
		exist, err := redisClient.HExists(key, username).Result()
		if err != nil {
			log.Printf("[FlushDiffToDbLog] 查找用户失败，error: %v\n", err)
			continue
		}
		if exist {
			value, err := redisClient.HGet(key, username).Result()
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
		_, err = redisClient.HDel(diffkey, username).Result()
		if err != nil {
			log.Printf("[FlushDiffToDbLog] 用户%s删除diff失败，error: %v\n", key, err)
		} else {
			log.Printf("[FlushDiffToDbLog] 用户%s写入数据库成功\n", key)
		}
	}
}

func flushDiffSortLinksKeyToDb() {
	key := utils.GetShortLinksKey()
	redisClient := db.GetRedisClient()
	linkIds, err := redisClient.ZRange(key, 0, -1).Result()
	if err != nil {
		log.Printf("[FlushDiffToDbLog] 查找短链接历史失败，error: %v\n", err)
		return
	}
	shortLinkHistoryMySqlRepo := repository.GetShortLinkHistoryMySqlRepo()
	for _, id := range linkIds {
		rhKey := utils.GetRequestHistoryKey(id)
		rs, err := redisClient.ZRangeWithScores(rhKey, 0, 100).Result()
		if err != nil {
			log.Printf("[FlushDiffToDbLog] 查找短链接历史失败，error: %v\n", err)
			continue
		}
		for _, one := range rs {
			rh := &models.RequestHistory{}
			err = json.Unmarshal([]byte(one.Member.(string)), rh)
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 短链接历史%s解析失败，error: %v\n", key, err)
				continue
			}
			err = shortLinkHistoryMySqlRepo.SaveOrUpdate(id, rh)
			if err != nil {
				log.Printf("[FlushDiffToDbLog] 短链接历史%s写入数据库失败，error: %v\n", key, err)
				continue
			}
			redisClient.ZRem(rhKey, one.Member)
		}
	}
}

func flushRequestHistoryToDb() {
	key := utils.GetShortLinksKey()
	redisClient := db.GetRedisClient()
	linkIds, err := redisClient.ZRange(key, 0, -1).Result()
	if err != nil {
		log.Printf("[FlushToDbLog] 查找短链接历史失败，error: %v\n", err)
		return
	}
	shortLinkHistoryMySqlRepo := repository.GetShortLinkHistoryMySqlRepo()
	for _, id := range linkIds {
		rhKey := utils.GetRequestHistoryKey(id)
		rs, err := redisClient.ZRangeWithScores(rhKey, 0, -1).Result()
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
			err = shortLinkHistoryMySqlRepo.SaveOrUpdate(id, rh)
			if err != nil {
				log.Printf("[FlushToDbLog] 短链接历史%s写入数据库失败，error: %v\n", key, err)
				continue
			}
		}
	}
}

func flushUserToDb() {
	key := utils.GetUserKey()
	redisClient := db.GetRedisClient()
	rhash, err := redisClient.HGetAll(key).Result()
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
	redisClient := db.GetRedisClient()
	linkIds, err := redisClient.ZRange(key, 0, -1).Result()
	if err != nil {
		log.Printf("[FlushToDbLog] 查找短链接失败，error: %v\n", err)
		return
	}
	for _, id := range linkIds {
		slRepo := repository.GetShortLinkRepo(redisClient)
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
