package cache

import (
	"errors"
	"fmt"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/system"
	"github.com/msean/botmanager/server/utils"
	"go.uber.org/zap"
)

type (
	SysCnfCache struct {
		Key   string `json:"key"`   //参数键
		Value string `json:"value"` //参数值
	}
)

var _ CacheItem = SysCnfCache{}

func NewSysCnfCache(key string) *SysCnfCache {
	return &SysCnfCache{
		Key: key,
	}
}

func (SysCnfCache) TableName() string { return system.SysParams{}.TableName() }
func (c SysCnfCache) Release() error  { return CacheDelete(c) }

func (c SysCnfCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"`key`", c.Key},
	}
}

func (SysCnfCache) LoadType() LoadType { return LoadFromDBGet }

func LoadSyscnf(key string, createIfNotExist bool, defaultVal string) (sysCache *SysCnfCache, err error) {
	sysCache = NewSysCnfCache(key)
	var has bool
	if has, err = CacheGetItem(sysCache); err != nil {
		return
	}
	if !has {
		if createIfNotExist {
			param := system.SysParams{
				Key:   key,
				Value: defaultVal,
				Name:  key,
				Desc:  "系统自动初始化默认参数",
			}
			if err = global.GVA_DB.Create(&param).Error; err != nil {
				return
			}
			sysCache.Key = key
			sysCache.Value = defaultVal
		} else {
			err = errors.New("Record Not Found")
			return
		}
	}
	return
}

func ReleaseSysCnf(modelID int) (err error) {
	var object system.SysParams
	var has bool
	if has, err = utils.Get(global.GVA_DB, &object, utils.IDCond(modelID)); !has || err != nil {
		if !has {
			err = fmt.Errorf("record not found")
		}
		global.GVA_LOG.Error("ReleaseSysCnf", zap.Any("id", modelID), zap.Error(err))
		return
	}

	if err = NewSysCnfCache(object.Key).Release(); err != nil {
		global.GVA_LOG.Error("ReleaseSysCnf", zap.Any("id", modelID), zap.Error(err))
	}
	return
}
