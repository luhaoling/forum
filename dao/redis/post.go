package redis

import (
	"project/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从 redis 获取 id
	// 1. 根据用户请求中携带的 order 参数确定要查询的 redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2.确定查询的索引起始点
	return getIDsFormKey(key, p.Page, p.Size)
}

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZRERANGE 按分数从大到小的顺序查询指定数量的元素(从 redis 中获取指定数量的 id 列表，降序排序返回)
	return client.ZRevRange(key, start, end).Result()
}

// GetPostVoteData 根据 ids 查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// 使用 Pipeline 机制来进行批量操作
	pipeline := client.Pipeline()
	// 循环遍历每个帖子 ID ，构建响应的 redis 键，计算集合中分数为 1 的元素数量
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区查询 ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostTimeZSet)
	}

	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存 key 减少 zinterstore 执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	// 如果缓存键不存在,则需要计算帖子 ID 列表
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算、
		pipeline := client.Pipeline()
		// 将社区 ID 对应的帖子 ID 集合与帖子时间的有序集合求交集
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		// 将结果存储在 key 中
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)
}
