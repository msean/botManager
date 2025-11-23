package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	LoadFromDBGet  LoadType = "get"  // 单条
	LoadFromDBList LoadType = "list" // 列表
	DefaultTimeout          = 24 * time.Hour
)

type (
	KvPkPair struct {
		PKCol string
		PKVal any
	}

	KvCacheObject struct {
		table string     // 表名
		pairs []KvPkPair // 主键值
	}

	LoadType string

	// CacheItem 可扩展，目前工具函数中不强制使用
	CacheItem interface {
		TableName() string
		Pairs() []KvPkPair
		LoadType() LoadType
	}
)

var cachePrefix = "sf_backend"

// ---------------------------------------------------------------------------
// KvCacheObject 工具方法
// ---------------------------------------------------------------------------

// 创建 CacheObject
func NewKvCacheObject(table string, pairs []KvPkPair) *KvCacheObject {
	return &KvCacheObject{table: table, pairs: pairs}
}

// 类别前缀 key:  zk:{table}
func (c *KvCacheObject) categoryKey() string {
	return fmt.Sprintf("%s:%s", cachePrefix, c.table)
}

// individualKey 对应每条缓存的唯一 key，例如 "123", "1_2024_abc"
func (c *KvCacheObject) individualKey() string {
	if len(c.pairs) == 0 {
		return ""
	}
	key := fmt.Sprintf("%v", c.pairs[0].PKVal)
	for i := 1; i < len(c.pairs); i++ {
		key += "_" + fmt.Sprintf("%v", c.pairs[i].PKVal)
	}
	return key
}

// 完整 Redis key
func (c *KvCacheObject) Key() string {
	return fmt.Sprintf("%s:%s", c.categoryKey(), c.individualKey())
}

func (c *KvCacheObject) Set(obj any, ttl ...time.Duration) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	expire := DefaultTimeout
	if len(ttl) > 0 {
		expire = ttl[0]
	}

	return global.GVA_REDIS.Set(
		context.Background(),
		c.Key(),
		data,
		expire,
	).Err()
}

func (c *KvCacheObject) Release() error {
	return global.GVA_REDIS.Del(context.Background(), c.Key()).Err()
}

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

func (c *KvCacheObject) fromDB(obj any) (bool, error) {
	conds := make([]utils.Cond, 0, len(c.pairs))
	for _, p := range c.pairs {
		conds = append(conds, utils.NewWhereCond(p.PKCol, p.PKVal))
	}

	has, err := utils.Get(global.GVA_DB.Table(c.table), obj, conds...)
	return has, err
}

func (c *KvCacheObject) listFromDB(obj any) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("listFromDB: obj must be a non-nil pointer")
	}
	v = v.Elem()

	// 找到唯一的 slice 字段
	var sliceField reflect.Value
	var sliceFieldType reflect.Type

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Slice {
			sliceField = f
			sliceFieldType = f.Type()
			break
		}
	}

	if !sliceField.IsValid() {
		return fmt.Errorf("listFromDB: no slice field found in %T", obj)
	}

	// 创建 slice 用于查询 DB
	slicePtr := reflect.New(sliceFieldType)
	sliceVal := slicePtr.Interface()

	// 构造查询条件
	conds := make([]utils.Cond, 0, len(c.pairs))
	for _, p := range c.pairs {
		conds = append(conds, utils.NewWhereCond(p.PKCol, p.PKVal))
	}

	// 查询数据库
	err := utils.Find(global.GVA_DB.Table(c.table), sliceVal, conds...)
	if err != nil {
		return err
	}

	// 写回 obj 的 slice 字段
	sliceField.Set(slicePtr.Elem())
	return nil
}

func CacheGetItem(item CacheItem) (bool, error) {
	if item == nil {
		return false, fmt.Errorf("CacheGetItem: item is nil")
	}

	table := item.TableName()
	pairs := item.Pairs()
	mode := item.LoadType()

	if table == "" {
		return false, fmt.Errorf("CacheGetItem: empty table name")
	}
	if len(pairs) == 0 {
		return false, fmt.Errorf("CacheGetItem: empty primary key pairs")
	}

	cache := NewKvCacheObject(table, pairs)

	// 1. 尝试读取缓存
	has, err := cache.Get(item)
	if err == nil && has {
		return true, nil
	}

	// 2. 读取数据库
	switch mode {
	case LoadFromDBGet:
		has, err = cache.fromDB(item)
	case LoadFromDBList:
		err = cache.listFromDB(item)
		has = err == nil
	default:
		return false, fmt.Errorf("unknown load type: %v", mode)
	}
	if err != nil || !has {
		return has, err
	}

	// 3. 写回缓存
	if setErr := cache.Set(item); setErr != nil {
		global.GVA_LOG.Warn("[cache] set failed", zap.String("key", cache.Key()), zap.Error(setErr))
	}

	return has, nil
}

func CacheDelete(item CacheItem) error {
	if item == nil {
		return fmt.Errorf("CacheDeleteItem: item is nil")
	}
	table := item.TableName()
	pairs := item.Pairs()

	if table == "" {
		return fmt.Errorf("CacheDeleteItem: empty table name")
	}
	if len(pairs) == 0 {
		return fmt.Errorf("CacheDeleteItem: empty primary key pairs")
	}

	cache := NewKvCacheObject(table, pairs)
	err := cache.Release()
	if err != nil {
		global.GVA_LOG.Warn("[cache] delete failed", zap.String("key", cache.Key()), zap.Error(err))
	}
	return err
}
