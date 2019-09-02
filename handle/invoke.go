package handle

import (
	"browser/model"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

func (f FabSdk)GetTokens()([]model.Token,error) {
	client, err := channel.New(f.channeProvider)
	if err != nil {
		fmt.Printf("Failed to create new channel client: %s", err.Error())
		return nil , err
	}
	response,err := client.Query(channel.Request{
		ChaincodeID:chaincodeid,
		Fcn:"token_list",
		Args:[][]byte{},
	},channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return nil , err
	}

	payload := model.Payload{}
	err = json.Unmarshal(response.Payload,&payload)
	if err != nil {
		return nil , err
	}
	records := make([]model.RecordToken,0)
	err = json.Unmarshal([]byte(payload.Message),&records)
	if err != nil {
		return nil , err
	}

	listTOken := make([]model.Token,0)
	for _,v := range records{
		listTOken = append(listTOken, v.Record)
	}
	return listTOken ,nil
}


func (f FabSdk)GetTokenHistory(tokenName string)([]model.HistoryToken,error) {
	client, err := channel.New(f.channeProvider)
	if err != nil {
		fmt.Printf("Failed to create new channel client: %s", err.Error())
		return nil,err
	}
	fmt.Println(tokenName)
	args := make([][]byte,0)
	args = append(args,[]byte(tokenName))
	response,err := client.Query(channel.Request{
		ChaincodeID:chaincodeid,
		Fcn:"token_history",
		Args:args,
	})
	if err != nil {
		return nil, err
	}

	payload := model.Payload{}
	err = json.Unmarshal(response.Payload,&payload)
	if err != nil {
		return nil , err
	}

	records := make([]model.HistoryToken,0)
	err = json.Unmarshal([]byte(payload.Message),&records)
	if err != nil {
		return nil , err
	}

	return  records,nil
}