package utils

import (
	"browser/model"
	"fmt"
	"testing"
)

func TestInsertBlock(t *testing.T) {
	init:=InitSql()
	if init {

		block := model.BlockHeader{
			Number:uint64(1),
			PreviousHash:"11111",
			DataHash:"2222",
		}
		err := InsertBlock(block)
		if err != nil {
			fmt.Println(err.Error())
		}
	}else{
		fmt.Println("not init")
	}
	CloseSql()
}

func TestInsertToken(t *testing.T) {
	init:=InitSql()
	if init {

		token := model.Token{
			Amount:0,
			Desc:"BABA.T.T Token burned amount :30000.00",
			Issuer:"mmadmin",
			Name:"BABA.T.T",
			Status:true,
			Type:"token",
		}
		err := InsertToken(token)
		if err != nil {
			fmt.Println(err.Error())
		}
	}else{
		fmt.Println("not init")
	}
	CloseSql()
}
