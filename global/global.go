package global

import (
	"go-admin/config"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	GVA_DB *gorm.DB
	//	GVA_DBList map[string]*gorm.DB
	GVA_REDIS  *redis.Client
	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	// GVA_LOG    *oplogging.Logger
	// GVA_LOG                 *zap.Logger
	// GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control = &singleflight.Group{}

	BlackCache local_cache.Cache
	//lock       sync.RWMutex
)

type GVA_MODEL struct {
	ID        uint           `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

// // GetGlobalDBByDBName 通过名称获取db list中的db
// func GetGlobalDBByDBName(dbname string) *gorm.DB {
// 	lock.RLock()
// 	defer lock.RUnlock()
// 	return GVA_DBList[dbname]
// }

// // MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
// func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
// 	lock.RLock()
// 	defer lock.RUnlock()
// 	db, ok := GVA_DBList[dbname]
// 	if !ok || db == nil {
// 		panic("db no init")
// 	}
// 	return db
// }
