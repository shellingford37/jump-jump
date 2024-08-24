package server

import (
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
	"log"
)

func flushToDb() error {
	flushSortLinksKeyToDb()
	flushUserToDb()
	return nil
}

func flushUserToDb() {

}

func flushSortLinksKeyToDb() {
	key := utils.GetShortLinksKey()
	var redis = db.GetRedisClient()
	linkIds, _ := redis.ZRange(key, 0, -1).Result()
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
