package service

import (
	storageCore "TTPanel/internal/core/storage"
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
)

type StorageService struct {
}

// AddStorage 添加存储
func (s *StorageService) AddStorage(param *request.AddStorageR) error {
	storage, err := s.NewStorageCore(param.Category, &param.StorageConfigR)
	if err != nil {
		return err
	}
	//err = storage.

	//构造存储配置
	configStr, err := util.StructToJsonStr(storage.GetConfig())
	if err != nil {
		return err
	}
	storageData := &model.Storage{
		Name:        param.Name,
		Category:    param.Category,
		Description: param.Description,
		Config:      configStr,
	}

	//查询名称是否存在
	get, err := storageData.Get(global.PanelDB)
	if err != nil {
		return err
	}
	if get.ID > 0 {
		return errors.New("storage name already exists")
	}

	//插入存储数据
	_, err = storageData.Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// StorageBucketList 获取存储桶列表
func (s *StorageService) StorageBucketList(param *request.StorageBucketListR) (bucketList []string, err error) {
	storage, err := s.NewStorageCore(param.Category, &param.StorageConfigR)
	if err != nil {
		return nil, err
	}
	list, err := storage.GetBucketList()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetLocalStoragePath 获取本地存储路径
func (s *StorageService) GetLocalStoragePath() string {
	return global.Config.System.DefaultBackupDirectory
}

// EditLocalStoragePath 编辑本地存储路径
func (s *StorageService) EditLocalStoragePath(path string) error {
	if util.IsDir(path) == false {
		return errors.New(fmt.Sprintf("path %v not a directory", path))
	}
	global.Config.System.DefaultBackupDirectory = path
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}

	return nil
}

// StorageList 获取存储列表
func (s *StorageService) StorageList(query string, category int, offset, limit int) (storageList []*model.Storage, total int64, err error) {
	whereT := model.ConditionsT{"ORDER": "create_time DESC"}
	whereOrT := model.ConditionsT{}
	if !util.StrIsEmpty(query) {
		query = "%" + query + "%"
		whereT["name LIKE ?"] = query
		whereOrT["description LIKE ?"] = query
	}
	if category > 0 {
		whereT["category"] = category
	}
	return (&model.Storage{}).List(global.PanelDB, &whereT, &whereOrT, offset, limit)

}

// EditStorage 编辑存储
func (s *StorageService) EditStorage(param *request.EditStorageR) error {
	storage, err := s.NewStorageCore(param.Category, &param.StorageConfigR)
	if err != nil {
		return err
	}
	//构造存储配置
	configStr, err := util.StructToJsonStr(storage.GetConfig())
	if err != nil {
		return err
	}
	storageData := &model.Storage{
		ID:          param.ID,
		Name:        param.Name,
		Category:    param.Category,
		Description: param.Description,
		Config:      configStr,
	}

	//查询名称是否存在
	get, err := storageData.Get(global.PanelDB)
	if err != nil {
		return err
	}
	if get.ID == 0 {
		return errors.New("storage does not exist")
	}
	if get.ID > 0 && get.ID != param.ID {
		return errors.New("storage name already exists")
	}

	//更新存储数据
	err = storageData.Update(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// NewStorageCore 创建存储核心
func (s *StorageService) NewStorageCore(category int, config *request.StorageConfigR) (storageCore.Storage, error) {
	var storageConfig interface{}
	switch category {
	case constant.StorageCategoryByTencentCOS:
		storageConfig = config.TencentCosConfig
	case constant.StorageCategoryByS3:
		storageConfig = config.S3Config
	case constant.StorageCategoryByQiniuKodo:
		storageConfig = config.QiniuKodoConfig
	case constant.StorageCategoryByAliOSS:
		storageConfig = config.AliOssConfig
	case constant.StorageCategoryByMinio:
		storageConfig = config.MinioConfig
	default:
		return nil, errors.New(fmt.Sprintf("unknown storage category: %d", category))
	}
	return storageCore.New(category, storageConfig)
}
