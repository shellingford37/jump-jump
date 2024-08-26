package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shellingford37/jump-jump/internal/app/config"
	_ "github.com/shellingford37/jump-jump/internal/app/config"
	"github.com/shellingford37/jump-jump/internal/app/db"
	"github.com/shellingford37/jump-jump/internal/app/routers"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"time"
)

var group errgroup.Group

const (
	flushToDbInterval = 1 * time.Minute
)

func setupDB() error {
	c := db.GetRedisClient()
	pong := c.Ping()
	return pong.Err()
}

// 检查 ALLOWED_HOSTS 设置正确设置
func allowHostsChecking() error {
	if gin.Mode() == gin.ReleaseMode {

		if funk.ContainsString([]string{"", "*"}, os.Getenv("ALLOWED_HOSTS")) {
			return fmt.Errorf("please set ALLOWED_HOSTS environment variable when GIN_MODE=release.\n")
		}
	}

	return nil
}

func Run(addr ...string) error {
	// security checking
	err := allowHostsChecking()

	if err != nil {
		return err
	}

	err = db.OpenMysql()

	if err != nil {
		return err
	}

	err = config.SetupConfig(db.GetRedisClient())

	if err != nil {
		return err
	}

	clearRedis()

	flushToRedisFlag := os.Getenv("MYSQL_TO_REDIS")
	if flushToRedisFlag == "true" {
		err = flushToRedis()
		if err != nil {
			return err
		}
	}

	flushToDbFlag := os.Getenv("REDIS_TO_MYSQL")
	if flushToDbFlag == "true" {
		err = flushToDb()
		if err != nil {
			return err
		}
	}

	group.Go(func() error {
		log.Println("[flushDiffToDb] ticker starts to serve")
		startFlushDiffToDbTicker()
		return nil
	})

	router := routers.SetupRouter()
	err = router.Run(addr...)
	return err
}

func RunLanding(addr ...string) error {
	err := setupDB()

	if err != nil {
		return err
	}

	err = db.OpenMysql()

	if err != nil {
		return err
	}

	err = config.SetupConfig(db.GetRedisClient())

	if err != nil {
		return err
	}

	router := routers.SetupLandingRouter()
	err = router.Run(addr...)
	return err
}

func startFlushDiffToDbTicker() {
	flushToDbTicker := time.NewTicker(flushToDbInterval)
	for range flushToDbTicker.C {
		log.Println("[flushDiffToDb] Start.")
		err := flushDiffToDb()
		if err != nil {
			log.Printf("FlushToDbLog error %s", err)
		}
		log.Println("[flushDiffToDb] Finish.")
	}
}
