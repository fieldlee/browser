package utils

import (
	"fmt"
	"testing"
)

var clusterKey = "clusterKey"
func BenchmarkRedisClusterClient_SetKey(b *testing.B) {
	clusterCli := InitRedisCluster()
	defer clusterCli.Close()
	for i:=0;i<100 ;i++  {
		value := fmt.Sprintf("clusterTestValue:%d",i)
		err := clusterCli.SetKey(clusterKey,value)
		if err != nil{
			fmt.Errorf("BenchmarkRedisClusterClient_SetKey ERR :%s",err.Error())
			continue
		}
	}
}

func BenchmarkRedisClusterClient_GetKey(b *testing.B) {
	clusterCli := InitRedisCluster()
	defer clusterCli.Close()
	for i:=0;i<100 ;i++  {
		v,err := clusterCli.GetKey(clusterKey)
		if err != nil{
			fmt.Errorf("BenchmarkRedisClusterClient_SetKey ERR :%s",err.Error())
			continue
		}
		fmt.Printf("BenchmarkRedisClusterClient_GetKey Value:%s \n",v)
	}
}

func BenchmarkRedisClusterClient_DelKey(b *testing.B) {
	clusterCli := InitRedisCluster()
	defer clusterCli.Close()
	err := clusterCli.DelKey(clusterKey)
	if err != nil{
		fmt.Errorf("BenchmarkRedisClusterClient_DelKey ERR :%s",err.Error())
	}
}

func BenchmarkRedisClusterClient_PushList(b *testing.B) {
	clusterCli := InitRedisCluster()
	defer clusterCli.Close()
	for i:=0;i<10000 ;i++  {
		value := fmt.Sprintf("clusterTestValue:%d",i)
		v,err := clusterCli.PushList(clusterKey,value)
		if err != nil{
			fmt.Errorf("BenchmarkRedisClusterClient_SetKey ERR :%s",err.Error())
			continue
		}
		fmt.Printf("PushList Value:%d \n",v)
	}
}

func BenchmarkRedisClusterClient_PopList(b *testing.B) {
	clusterCli := InitRedisCluster()
	defer clusterCli.Close()
	for i:=0;i<10000 ;i++  {
		v,err := clusterCli.PopList(clusterKey)
		if err != nil{
			fmt.Errorf("BenchmarkRedisClusterClient_SetKey ERR :%s",err.Error())
			continue
		}
		fmt.Printf("PopList Value:%s \n",v)
	}
}