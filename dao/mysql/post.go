package mysql

import (
	"project/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetPostListByIDs 根据给定的 ID 列表查询帖子的数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
		from post
		where post_id in(?)
		order by FIND_IN_SET(post_id,?)
	`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

// GetPostById 根据 id 查询单个帖子数据
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
		from post
		where post_id=?
	`
	err = db.Get(post, sqlStr, pid)
	return
}

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
                 post_id,title,content,author_id,community_id)
    values(?,?,?,?,?)
    `
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return

}
