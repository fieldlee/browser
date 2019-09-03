package utils

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedis()*RedisClient{
	redisClient := &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "p@ss1234",
			//OnConnect: func(conn *redis.Conn) error {
			//	intCmd := conn.DBSize()
			//	dbINt,err := intCmd.Result()
			//	if err != nil {
			//		return err
			//	}
			//	fmt.Printf("dbSize:%d\n",dbINt)
			//
			//	strcmd := conn.ClientGetName()
			//	nameStr,err := strcmd.Result()
			//	if err != nil {
			//		fmt.Errorf("ClientGetName Err :%s",err.Error())
			//		return err
			//	}
			//	fmt.Printf("client Name:%s\n",nameStr)
			//	return nil
			//},
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

func (r *RedisClient)GetKey(key string)(string,error){
	cmd:=r.Client.Get(key)
	if cmd.Err() != nil {
		return "",cmd.Err()
	}
	return cmd.Result()
}

func (r *RedisClient)PushList(key string,value string)(int64,error){
	cmd := r.Client.LPush(key,value)
	if cmd.Err() != nil {
		return int64(0),cmd.Err()
	}
	return cmd.Result()
}

func (r *RedisClient)PopList(key string)(string,error){
	cmd := r.Client.RPop(key)
	if cmd.Err() != nil {
		return "",cmd.Err()
	}
	return cmd.Result()
}

func (r *RedisClient)HSet(hkey string,mkey string,v interface{})(bool,error){
	resultcmd := r.Client.HSet(hkey,mkey,v)
	if resultcmd.Err() != nil {
		return false, resultcmd.Err()
	}
	return resultcmd.Result()
}

func (r *RedisClient)HGet(hkey string,mkey string)(string,error){

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

func (r *RedisClient)DelKey(keys ...string)error{
	cmd:=r.Client.Del(keys...)
	if cmd.Err() != nil{
		return cmd.Err()
	}
	return nil
}

func (r *RedisClient)Close(){
	err := r.Client.Close()
	if err != nil{
		fmt.Errorf(err.Error())
	}
}