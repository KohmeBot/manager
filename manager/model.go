package manager

import (
	"gorm.io/gorm"
	"time"
)

type BanRecord struct {
	GroupId int64 `gorm:"primaryKey"`
	UserId  int64 `gorm:"primaryKey"`
	// 被ban的次数
	Count int
	// 上次被ban的时间
	LastBan time.Time
}

func (r *BanRecord) NeedBan(banCd time.Duration, groupId, userId int64, db *gorm.DB) (needBan bool, err error) {
	r.GroupId = groupId
	r.UserId = userId
	err = db.Find(&r).Error
	if err != nil {
		return false, err
	}
	nowTime := time.Now()

	// 24小时内触发
	if nowTime.Sub(r.LastBan) <= banCd {
		needBan = true
	} else {
		r.Count = 0
	}

	r.Count++
	r.LastBan = nowTime
	err = db.Save(&r).Error
	return
}
