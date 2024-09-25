package models

import (
	"GinChat/utils"
	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint // 谁的关系
	TargetId uint // 对应的谁
	Type     int  // 对应的类型 1好友 2群
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriends(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id =? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

func AddFriend(userId uint, targetId uint) (int, string) {
	user := UserBasic{}
	if targetId != 0 {
		user = FindUserByID(targetId)
		if user.Password != "" {
			if userId == user.ID {
				return -5, "无法添加自己"
			}

			temp := Contact{}
			utils.DB.Where("owner_id =? and target_id = ? and type=1", userId, targetId).Find(&temp)
			if temp.ID != 0 {
				return -3, "对方已经是你的好友"
			}

			tx := utils.DB.Begin() // 打开事物
			defer func() {
				// 事物一旦开始，不论最后什么异常都会回滚
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact := Contact{
				OwnerId:  userId,
				TargetId: targetId,
				Type:     1,
				Desc:     "好友",
			}
			if err := utils.DB.Create(&contact); err != nil {
				tx.Rollback()
				return -4, "创建好友关系发生错误"
			}
			contact = Contact{
				OwnerId:  targetId,
				TargetId: userId,
				Type:     1,
				Desc:     "好友",
			}
			if err := utils.DB.Create(&contact); err != nil {
				tx.Rollback()
				return -4, "创建好友关系发生错误"
			}
			tx.Commit()
			return 0, "添加好友成功"
		}
		return -1, "未找到该用户"
	}
	return -2, "输入不许为空"
}

func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
