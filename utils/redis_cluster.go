package utils

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

type RedisClusterClient struct {
	Client *redis.ClusterClient
}

func InitRedisCluster()RedisClusterClient{
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"127.0.0.1:7000","127.0.0.1:7001","127.0.0.1:7002","127.0.0.1:7003","127.0.0.1:7004","127.0.0.1:7005"},
	})

	return RedisClusterClient{
		Client:client,
	}
}


func (r *RedisClusterClient)SetKey(key string,value string)error{
	cmd := r.Client.Set(key,value,0)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r *RedisClusterClient)GetKey(key string)(string,error){
	cmd:=r.Client.Get(key)
	if cmd.Err() != nil {
		return "",cmd.Err()
	}
	return cmd.Result()
}

func (r *RedisClusterClient)PushList(key string,value string)(int64,error){
	cmd := r.Client.LPush(key,value)
	if cmd.Err() != nil {
		return int64(0),cmd.Err()
	}
	return cmd.Result()
}

func (r *RedisClusterClient)PopList(key string)(string,error){
	cmd := r.Client.RPop(key)
	if cmd.Err() != nil {
		return "",cmd.Err()
	}
	return cmd.Result()
}

func (r *RedisClusterClient)HSet(hkey string,mkey string,v interface{})(bool,error){
	resultcmd := r.Client.HSet(hkey,mkey,v)
	if resultcmd.Err() != nil {
		return false, resultcmd.Err()
	}
	return resultcmd.Result()
}

func (r *RedisClusterClient)HGet(hkey string,mkey string)(string,error){

	b := r.Client.HExists(hkey,mkey)
	if b.Err() != nil {
		return "",b.Err()
	}
	if !b.Val() {
		return "",errors.New(fmt.Sprintf("the %s key not exist",mkey))
	}

	vcmd:=r.Client.HGet(hkey,mkey)
	if vcmd.Err() != nil {
		return "",b.Err()
	}
	return vcmd.Val(),nil
}

func (r *RedisClusterClient)DelKey(keys ...string)error{
	cmd:=r.Client.Del(keys...)
	if cmd.Err() != nil{
		return cmd.Err()
	}
	return nil
}

func (r *RedisClusterClient)Close(){
	err := r.Client.Close()
	if err != nil{
		fmt.Errorf(err.Error())
	}
}
