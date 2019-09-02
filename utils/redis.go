package utils

import "github.com/go-redis/redis"

type RedisClient struct {
	Client *redis.Client
}

func InitRedis()*RedisClient{
	redisClient := &RedisClient{
		Client:redis.NewClient(&redis.Options{
			Addr:"",
			DB:0,
			Password:"",
		}),
	}
	return redisClient
}

func (r *RedisClient)SetKey(key string,value string)error{
	cmd := r.Client.Set(key,value,0)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}