package model

//type TTWafOperationLog struct {
//	ID            int64  `gorm:"column:id" db:"id" json:"id" form:"id"`
//	Type          int    `gorm:"column:type" db:"type" json:"type" form:"type"`
//	Log           string `gorm:"column:log" db:"log" json:"log" form:"log"`
//	CreateTime    int64  `gorm:"column:create_time;autoCreateTime" db:"create_time" json:"create_time" form:"create_time"`
//	Uid           int64  `gorm:"column:uid" db:"uid" json:"uid" form:"uid"`
//	Username      string `gorm:"column:username" db:"username" json:"username" form:"username"`
//	IP            string `gorm:"column:ip" db:"ip" json:"ip" form:"ip"`
//	IPAttribution string `gorm:"column:ip_attribution" db:"ip_attribution" json:"ip_attribution" form:"ip_attribution"`
//}
//
//func (TTWafOperationLog) TableName() string {
//	return "operation_log"
//}
//
//func (s *TTWafOperationLog) Get(db *gorm.DB) (*TTWafOperationLog, error) {
//	var operationLog TTWafOperationLog
//	if s.ID > 0 {
//		db = db.Where("id", s.ID)
//	} else {
//		return nil, gorm.ErrRecordNotFound
//	}
//	err := db.Limit(1).Find(&operationLog).Error
//	if err != nil {
//		return &operationLog, err
//	}
//
//	return &operationLog, nil
//}
//func (s *TTWafOperationLog) List(db *gorm.DB, conditions *ConditionsT, offset, limit int) ([]*TTWafOperationLog, int64, error) {
//	var operationLog []*TTWafOperationLog
//	var err error
//	var count int64
//	if offset >= 0 && limit > 0 {
//		db = db.Offset(offset).Limit(limit)
//	}
//
//	for k, v := range *conditions {
//		if k == "ORDER" {
//			db = db.Order(v)
//		} else if k == "FIXED" {
//			db = db.Where(v)
//		} else {
//			db = db.Where(k, v)
//		}
//	}
//	if err = db.Find(&operationLog).Error; err != nil {
//		return nil, count, err
//	}
//
//	if err = db.Model(s).Count(&count).Error; err != nil {
//		return nil, count, err
//	}
//	return operationLog, count, nil
//}
//func (s *TTWafOperationLog) Create(db *gorm.DB) (*TTWafOperationLog, error) {
//	err := db.Create(&s).Error
//	return s, err
//}
//func (s *TTWafOperationLog) Update(db *gorm.DB) error {
//	return db.Model(&TTWafOperationLog{}).Where("id = ? ", s.ID).Save(s).Error
//}
//func (s *TTWafOperationLog) Delete(db *gorm.DB, conditions *ConditionsT) error {
//	for k, v := range *conditions {
//		db = db.Where(k, v)
//	}
//	return db.Delete(&TTWafOperationLog{}, s.ID).Error
//}
