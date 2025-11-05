package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 默认缓存过期时间
const DefaultTimeout = 24 * time.Hour

// 缓存加载方式
type LoadType string

const (
	LoadFromDBGet  LoadType = "get"
	LoadFromDBList LoadType = "list"
)

// KvPkPair 表示主键列与值
type KvPkPair struct {
	PKCol string
	PKVal any
}

// KvCacheObject 封装 Redis KV 缓存
type KvCacheObject struct {
	table string
	pairs []KvPkPair
}

// 创建新的缓存对象
func NewKvCacheObject(table string, pairs []KvPkPair) *KvCacheObject {
	return &KvCacheObject{table: table, pairs: pairs}
}

func (c *KvCacheObject) categoryKey() string {
	return fmt.Sprintf("zk:%s", c.table)
}

func (c *KvCacheObject) individualKey() string {
	if len(c.pairs) == 0 {
		return ""
	}

	key := fmt.Sprintf("%v", c.pairs[0].PKVal) // 转成字符串
	for i := 1; i < len(c.pairs); i++ {
		key += "_" + fmt.Sprintf("%v", c.pairs[i].PKVal)
	}
	return key
}

func (c *KvCacheObject) Key() string {
	return fmt.Sprintf("%s:%s", c.categoryKey(), c.individualKey())
}

// Set 设置缓存
func (c *KvCacheObject) Set(obj any, ttl ...time.Duration) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	expire := DefaultTimeout
	if len(ttl) > 0 {
		expire = ttl[0]
	}
	return global.GVA_REDIS.Set(context.Background(), c.Key(), data, expire).Err()
}

// Release 删除缓存
func (c *KvCacheObject) Release() error {
	return global.GVA_REDIS.Del(context.Background(), c.Key()).Err()
}

// Get 从缓存读取
func (c *KvCacheObject) Get(obj any) (bool, error) {
	data, err := global.GVA_REDIS.Get(context.Background(), c.Key()).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal(data, obj); err != nil {
		return false, err
	}
	return true, nil
}

// 从数据库读取单条
func (c *KvCacheObject) fromDB(obj any) (bool, error) {
	conds := make([]utils.Cond, 0, len(c.pairs))
	for _, pair := range c.pairs {
		conds = append(conds, utils.NewWhereCond(pair.PKCol, pair.PKVal))
	}
	has, err := utils.Get(global.GVA_DB.Table(c.table), obj, conds...)
	return has, err
}

// 从数据库读取多条
func (c *KvCacheObject) listFromDB(obj any) error {
	conds := make([]utils.Cond, 0, len(c.pairs))
	for _, pair := range c.pairs {
		conds = append(conds, utils.NewWhereCond(pair.PKCol, pair.PKVal))
	}
	return utils.Find(global.GVA_DB.Table(c.table), obj, conds...)
}

// 通用缓存加载逻辑
func CacheGet(table string, pairs []KvPkPair, obj any, mode LoadType) (bool, error) {
	cache := NewKvCacheObject(table, pairs)
	var (
		has bool
		err error
	)

	// 优先从缓存读
	has, err = cache.Get(obj)
	if err == nil && has {
		return true, nil
	}

	global.GVA_LOG.Debug("[cache] miss", zap.String("key", cache.Key()))

	// 缓存 miss，从数据库加载
	switch mode {
	case LoadFromDBGet:
		has, err = cache.fromDB(obj)
	case LoadFromDBList:
		err = cache.listFromDB(obj)
		has = err == nil
	default:
		err = fmt.Errorf("unknown load type: %v", mode)
	}

	if err != nil || !has {
		return has, err
	}

	// 写入缓存
	if setErr := cache.Set(obj); setErr != nil {
		global.GVA_LOG.Warn("[cache] write failed", zap.String("key", cache.Key()), zap.Error(setErr))
	}
	return has, nil
}

// CacheDelete 删除缓存
func CacheDelete(table string, pairs []KvPkPair) error {
	cache := NewKvCacheObject(table, pairs)
	err := cache.Release()
	if err != nil {
		global.GVA_LOG.Warn("[cache] delete failed", zap.String("key", cache.Key()), zap.Error(err))
	}
	return err
}

// CacheUpdate 更新缓存（支持强制从DB加载）
// func CacheUpdate(table string, pairs []KvPkPair, obj any, reload bool) error {
// 	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
// 		return fmt.Errorf("obj must be pointer type")

// 	}
// 	cache := NewKvCacheObject(table, pairs)

// 	// 删除旧缓存
// 	if err := cache.Release(); err != nil {
// 		return err
// 	}

// 	if !reload {
// 		return cache.Set(obj)
// 	}

// 	// 从数据库重新加载
// 	has, err := cache.fromDB(obj)
// 	if err != nil || !has {
// 		return err
// 	}

// 	return cache.Set(obj)
// }
