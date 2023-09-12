package mysql

import (
	"database/sql"
	"project/models"
)

// GetCommunityList 返回社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	err = db.Select(&communityList, sqlStr)
	return
}

// GetCommunityDetailByID 根据社区 id 返回社区详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select 
				community_id,community_name,introdection,create_time
				from community
				where community_id=?
				`
	if err := db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return community, err
}
