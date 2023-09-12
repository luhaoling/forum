package redis

const (
	Prefix             = "bluebell:"   // 项目前缀
	KeyPostTimeZSet    = "post:time"   // zset;帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  // zset;帖子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" // zset;记录用户及投票类型;参数时post id
	KeyCommunitySetPF  = "community:"  // set;保存每个分区下的帖子的 id
)

// 给 redis key 加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
