package utils

import (
	"fmt"
	"testing"
)

var RKey = "TestListKey"


func ExampleRedisClient_PushList() {
	redisCLient := InitRedis()
	defer redisCLient.Close()
	for i:=0;i<100 ;i++  {
		listValue := fmt.Sprintf("listValue_%d",i)
		v ,err := redisCLient.PushList(RKey,listValue)
		if err != nil {
			fmt.Errorf("redis push list value : %s ,err : %s \n",listValue,err.Error())
			continue
		}
		fmt.Printf("redis push list return value : %d \n",v)
	}
}

func TestRedisClient_PushList(t *testing.T) {
	redisCLient := InitRedis()
	defer redisCLient.Close()
	for i:=0;i<100 ;i++  {
		listValue := fmt.Sprintf("listValue_%d",i)
		v ,err := redisCLient.PushList(RKey,listValue)
		if err != nil {
			fmt.Errorf("redis push list value : %s ,err : %s \n",listValue,err.Error())
			continue
		}
		fmt.Printf("redis push list return value : %d \n",v)
	}
}

func BenchmarkRedisClient_PushList(b *testing.B) {
	redisCLient := InitRedis()
	defer redisCLient.Close()
	for i:=0;i<100 ;i++  {
		listValue := fmt.Sprintf("listValue_%d",i)
		v ,err := redisCLient.PushList(RKey,listValue)
		if err != nil {
			fmt.Errorf("redis push list value : %s ,err : %s \n",listValue,err.Error())
			continue
		}
		fmt.Printf("redis push list return value : %d \n",v)
	}

}

func TestRedisClient_PopList(t *testing.T) {
	redisCLient := InitRedis()
	defer redisCLient.Close()
	for i:=0;i<100 ;i++  {
		v,err := redisCLient.PopList(RKey)
		if err != nil {
			fmt.Errorf("redis pop list err : %s \n",err.Error())

		}
		fmt.Printf("redis pop list return value : %s \n",v)
	}
}

func BenchmarkRedisClient_PopList(b *testing.B) {
	redisCLient := InitRedis()
	defer redisCLient.Close()
	for i:=0;i<100 ;i++  {
		v,err := redisCLient.PopList(RKey)
		if err != nil {
			fmt.Errorf("redis pop list err : %s \n",err.Error())

		}
		fmt.Printf("redis pop list return value : %s \n",v)
	}
}