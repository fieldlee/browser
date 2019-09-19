package handle

import (
	"browser/model"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

func (f FabSdk)Invoke(funcName string,args []string)(model.Payload,string,error){
	client, err := channel.New(f.channeProvider)
	if err != nil {
		fmt.Printf("Failed to create new channel client: %s", err.Error())
		return  model.Payload{},"",err
	}
	argList := make([][]byte,0)

	for _,v := range args{
		argList = append(argList,[]byte(v))
	}
	response,err := client.Execute(channel.Request{
		ChaincodeID:chaincodeid,
		Fcn:funcName,
		Args:argList,
	},channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return model.Payload{},"",err
	}
	payload := model.Payload{}
	err = json.Unmarshal(response.Payload,&payload)
	if err != nil {
		return  model.Payload{},"",err
	}
	return payload,string(response.TransactionID),nil
}
