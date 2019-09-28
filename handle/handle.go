package handle

import (
	"browser/model"
	"browser/utils"
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/viper"
	"os"
)

var channelid, user, org,chaincodeid = "","","",""

type FabSdk struct {
	sdk *fabsdk.FabricSDK
	channeProvider context.ChannelProvider
}

func InitSdk()FabSdk{
	configFile := "/var/yaml/config_event.yaml"
	viper.AddConfigPath(configFile)
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		fmt.Printf("实例化Fabric SDK失败: %v\n", err)
		os.Exit(-1)
	}

	ymlCon:= utils.InitConfig()
	channelid = ymlCon.V.GetString("channel")
	user = ymlCon.V.GetString("name")
	org = ymlCon.V.GetString("orgname")
	chaincodeid = ymlCon.V.GetString("chaincodeid")

	chanProvider := sdk.ChannelContext(channelid,fabsdk.WithUser(user),fabsdk.WithOrg(org))

	if chanProvider == nil {
		fmt.Println("根据指定的组织名称与管理员创建资源管理客户端Context失败")
		os.Exit(-1)
	}
	return FabSdk{
		sdk:sdk,
		channeProvider:chanProvider,
	}
}

func (f FabSdk)Close(){
	f.sdk.Close()
	f.channeProvider = nil
}

func (f FabSdk)GetInfo() (*fab.BlockchainInfoResponse,error) {
	client , err := ledger.New(f.channeProvider)
	if err != nil {
		fmt.Errorf(err.Error())
		return &fab.BlockchainInfoResponse{},err
	}
	blcokinfo , err := client.QueryInfo()
	if err != nil {
		fmt.Errorf(err.Error())
		return &fab.BlockchainInfoResponse{},err
	}
	return blcokinfo,nil
}

func (f FabSdk)GetBlocks(height uint64)(model.Block,error){
	client , err := ledger.New(f.channeProvider)
	if err != nil {
		fmt.Errorf(err.Error())
		return model.Block{},err
	}
	blockinfo,err := client.QueryBlock(height)
	if err != nil {
		fmt.Errorf(err.Error())
		return model.Block{},err
	}

	listTx := make([]model.TransactionDetail,0)
	for _,data := range blockinfo.Data.Data{
		txDetail,err := utils.GetTransactionInfoFromData(data,true)
		if err != nil {
			fmt.Println(err.Error())
		}
		listTx = append(listTx,*txDetail)
	}

	block := model.Block{}
	block.DataHash = hex.EncodeToString(blockinfo.Header.GetDataHash())
	block.PreviousHash = hex.EncodeToString(blockinfo.Header.GetPreviousHash())
	block.Number = blockinfo.Header.Number
	// add create time
	if len(listTx)>0{
		block.CreateTime = listTx[0].CreateTime
	}

	block.TxList = listTx
	return block,nil
}

func (f FabSdk)GetBlocksByHash(hash string)(model.Block,error){
	client , err := ledger.New(f.channeProvider)
	if err != nil {
		fmt.Errorf(err.Error())
		return model.Block{},err
	}
	hashByte,err := hex.DecodeString(hash)
	if err != nil {
		fmt.Errorf(err.Error())
		return model.Block{},err
	}
	blockinfo,err := client.QueryBlockByHash(hashByte)
	if err != nil {
		fmt.Errorf(err.Error())
		return model.Block{},err
	}

	listTx := make([]model.TransactionDetail,0)
	for _,data := range blockinfo.Data.Data{
		txDetail,err := utils.GetTransactionInfoFromData(data,true)
		if err != nil {
			fmt.Println(err.Error())
		}
		listTx = append(listTx,*txDetail)
	}
	block := model.Block{}
	block.DataHash = hex.EncodeToString(blockinfo.Header.GetDataHash())
	block.PreviousHash = hex.EncodeToString(blockinfo.Header.PreviousHash)
	block.Number = blockinfo.Header.Number
	// add create time
	if len(listTx)>0{
		block.CreateTime = listTx[0].CreateTime
	}
	block.TxList = listTx
	return block,nil
}
func (f FabSdk)GetBlocksByTxId(hash string)(model.Block,error){
	client , err := ledger.New(f.channeProvider)
	if err != nil {
		fmt.Errorf(err.Error())
		return model.Block{},err
	}
	blockinfo,err := client.QueryBlockByTxID(fab.TransactionID(hash))
	if err != nil {
		fmt.Errorf("QueryBlockByTxID err :%s",err.Error())
		return model.Block{},err
	}
	listTx := make([]model.TransactionDetail,0)
	for _,data := range blockinfo.Data.Data{
		txDetail,err := utils.GetTransactionInfoFromData(data,true)
		if err != nil {
			fmt.Errorf("QueryBlockByTxID err :%s",err.Error())
			fmt.Println(err.Error())
		}
		listTx = append(listTx,*txDetail)
	}
	block := model.Block{}
	block.DataHash = hex.EncodeToString(blockinfo.Header.GetDataHash())
	block.PreviousHash = hex.EncodeToString(blockinfo.Header.PreviousHash)
	block.Number = blockinfo.Header.Number
	// add create time
	if len(listTx)>0{
		block.CreateTime = listTx[0].CreateTime
	}
	block.TxList = listTx
	return block,nil
}
func (f FabSdk)GetTransactionByTxId(txid string)(model.TransactionDetail,error){
	result := model.TransactionDetail{}
	client , err := ledger.New(f.channeProvider)
	if err != nil {
		fmt.Errorf(err.Error())
		return result,err
	}
	transaction,err := client.QueryTransaction(fab.TransactionID(txid))
	if err != nil {
		fmt.Errorf(err.Error())
		return result,err
	}
	args,err := utils.GetTransaction(transaction.TransactionEnvelope)
	if err != nil {
		fmt.Errorf(err.Error())
		return result,err
	}

	result.Args = args
	result.TransactionId = txid
	return result,nil
}

