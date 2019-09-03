package utils

import (
	"browser/model"
	"fmt"
	"testing"
)

func TestInsertBlock(t *testing.T) {
	client,err := InitSql()
	if err != nil {
		fmt.Println(err.Error())
	}
		block := model.BlockHeader{
			Number:uint64(1),
			PreviousHash:"11111",
			DataHash:"2222",
		}
		err = client.InsertBlock(block)
		if err != nil {
			fmt.Println(err.Error())
		}

	client.CloseSql()
}

func TestInsertToken(t *testing.T) {
	client,err := InitSql()
	if err != nil {
		fmt.Println(err.Error())
	}

		token := model.Token{
			Amount:0,
			Desc:"BABA.T.T Token burned amount :30000.00",
			Issuer:"mmadmin",
			Name:"BABA.T.T",
			Status:true,
			Type:"token",
		}
		err = client.InsertToken(token)
		if err != nil {
			fmt.Println(err.Error())
		}

	client.CloseSql()
}
