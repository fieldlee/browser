package utils

import (
	"browser/handle"
	"browser/model"
	"encoding/hex"

	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	putils "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/utils"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"time"
)

type Config struct {
	V *viper.Viper
}

var Con *Config

func InitConfig () *Config {
	Con := &Config{
		V:viper.New(),
	}
	//设置配置文件的名字
	Con.V.SetConfigName("config")
	//添加配置文件所在的路径,注意在Linux环境下%GOPATH要替换为$GOPATH

	Con.V.AddConfigPath("./")
	//设置配置文件类型
	Con.V.SetConfigType("yaml")
	if err := Con.V.ReadInConfig(); err != nil{
		log.Fatal(err.Error())
	}
	return Con
}

func GetWhiteIPs()[]string{
	ymlCon := InitConfig()
	whiteiplist := ymlCon.V.Get("whiteips")
	iplist := make([]string,0)
	v := reflect.ValueOf(whiteiplist)
	if v.Kind() == reflect.Slice {
		l := v.Len()
		for i:=0;i<l;i++  {
			iplist = append(iplist,fmt.Sprintf("%v",v.Index(i)))
		}
	}
	return iplist
}

func Contain(list []string,obj string)bool{
	for _,tmp := range list{
		if tmp == obj{
			return  true
		}
	}
	return false
}

// 从SDK中Block.BlockDara.Data中提取交易具体信息
func GetTransactionInfoFromData(data []byte, needArgs bool) (*model.TransactionDetail, error) {
	env, err := putils.GetEnvelopeFromBlock(data)
	if err != nil {
		return nil, err
	}
	if env == nil {
		return nil, errors.New("nil envelope")
	}
	payload, err := putils.GetPayload(env)
	if err != nil {
		return nil, err
	}
	channelHeaderBytes := payload.Header.ChannelHeader
	channelHeader := &cb.ChannelHeader{}

	if err := proto.Unmarshal(channelHeaderBytes, channelHeader); err != nil {
		return nil, err
	}
	var (
		args []string
	)
	if needArgs {
		tx, err := putils.GetTransaction(payload.Data)
		if err != nil {
			return nil, err
		}
		chaincodeActionPayload, err := putils.GetChaincodeActionPayload(tx.Actions[0].Payload)
		if err != nil {
			return nil,err
		}
		propPayload := &pb.ChaincodeProposalPayload{}
		if err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, propPayload); err != nil {
			return nil, err
		}
		invokeSpec := &pb.ChaincodeInvocationSpec{}
		err = proto.Unmarshal(propPayload.Input, invokeSpec)
		if err != nil {
			return nil, err
		}
		if invokeSpec.ChaincodeSpec != nil && invokeSpec.ChaincodeSpec.Input != nil  &&  invokeSpec.ChaincodeSpec.Input.Args != nil {
			for _, v := range invokeSpec.ChaincodeSpec.Input.Args {
				args = append(args, string(v))
			}
		}

	}
	result := &model.TransactionDetail{
		TransactionId: channelHeader.TxId,
		Args:          args,
		CreateTime:    time.Unix(channelHeader.Timestamp.Seconds, 0).Format("2006-01-02 15:04:05"),
	}
	return result, nil
}

func GetTransaction(e *cb.Envelope)([]string,error){
	args := make([]string,0)

	payload, err := putils.GetPayload(e)
	if err != nil {
		return nil, err
	}

	tx, err := putils.GetTransaction(payload.Data)
	if err != nil {
		return nil, err
	}
	chaincodeActionPayload, err := putils.GetChaincodeActionPayload(tx.Actions[0].Payload)
	if err != nil {
		return nil,err
	}

	propPayload := &pb.ChaincodeProposalPayload{}

	if err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, propPayload); err != nil {
		return nil, errors.New(fmt.Sprintf("Unmarshal ChaincodeProposalPayload Error:%s",err.Error()))
	}

	invokeSpec := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(propPayload.Input, invokeSpec)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unmarshal ChaincodeInvocationSpec Error:%s",err.Error()))
	}

	for _, v := range invokeSpec.ChaincodeSpec.Input.Args {
		args = append(args, string(v))
	}
	return args,nil
}

func UpdateBlockAndTx(block cb.Block)error{
	sqlClient,err := InitSql()
	if err != nil {
		return err
	}
	defer sqlClient.CloseSql()
	listTx := make([]model.TransactionDetail,0)
	for _,data := range block.Data.Data{
		txDetail,err := GetTransactionInfoFromData(data,true)
		if err != nil {
			fmt.Println(err.Error())
		}
		listTx = append(listTx,*txDetail)
	}

	bck := model.Block{}
	bck.DataHash = hex.EncodeToString(block.Header.GetDataHash())
	bck.PreviousHash = hex.EncodeToString(block.Header.PreviousHash)
	bck.Number = block.Header.Number
	bck.TxList = listTx

	bckHeader := model.BlockHeader{}
	bckHeader.DataHash = bck.DataHash
	bckHeader.Number = bck.Number
	bckHeader.PreviousHash = bck.PreviousHash

	//// update tx block
	//fmt.Println(fmt.Sprintf("update previous %d",int(bck.Number-1)))
	//prebck,err := sqlClient.QueryBlockByHeight(int(bck.Number-1))
	//if err != nil {
	//	return err
	//}
	//fmt.Println(fmt.Sprintf("previous block hash %s previous %s",prebck.DataHash,bck.PreviousHash))
	//err = sqlClient.UpdateTxHash(prebck.DataHash,bck.PreviousHash)
	//if err != nil {
	//	return err
	//}
	//// update block hash
	//fmt.Println(fmt.Sprintf("update previous block hash number %d",int(bck.Number-1)))
	//err = sqlClient.UpdateBlockHash(int(bck.Number-1),bck.PreviousHash)
	//if err != nil {
	//	return err
	//}

	err = sqlClient.InsertBlock(bckHeader)
	if err != nil {
		return err
	}

	for i := 0;i<len(bck.TxList);i++{
		err = sqlClient.InsertTx(bck.DataHash,bck.TxList[i])
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}

func GetTxDetail(block cb.Block)*model.TransactionDetail{
	listTx := make([]model.TransactionDetail,0)
	for _,data := range block.Data.Data{
		txDetail,err := GetTransactionInfoFromData(data,true)
		if err != nil {
			fmt.Println(err.Error())
		}
		listTx = append(listTx,*txDetail)
	}
	if len(listTx) == 0 {
		return nil
	}
	return &listTx[0]
}

func TypeSwitch(arg interface{}){
	vType := reflect.TypeOf(arg)
	switch vType.Name() {
	case "string":
		fmt.Printf("String:%s\n",vType.String())
		fmt.Printf("name:%s\n",vType.Name())

	}
}