package utils

import (
	"browser/model"
	"encoding/json"
	"github.com/rhinoman/couchdb-go"
	"time"
)

type CouchClient struct {
	DB *couchdb.Database
}

func InitCouchClient()(*CouchClient,error){
	timeout := time.Duration(500 * time.Millisecond)
	conn, err := couchdb.NewConnection("192.168.1.31", 6984, timeout)
	if err != nil {
		return nil, err
	}
	auth := couchdb.BasicAuth{Username: "couchadmin", Password: "adminpwd"}
	db := conn.SelectDB("mmchannel_ledger", &auth)
	return &CouchClient{
		DB:db,
	},nil
}

func(c *CouchClient)GetAccounts()([]model.Account,error){
	selector := `{"type":"account"}`

	var selectorObj interface{}
	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return nil, err
	}
	params := couchdb.FindQueryParams{Selector: &selectorObj}
	accoutsResp := model.CouchAccountList{}
	err = c.DB.Find(&accoutsResp,&params)
	if err != nil {
		return nil, err
	}
	return accoutsResp.Docs,nil
}

func(c *CouchClient)GetAccount(name string)(model.Account,error){
	selector := `{"type":"account","cn":"`+name+`"}`

	var selectorObj interface{}
	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return model.Account{}, err
	}
	params := couchdb.FindQueryParams{Selector: &selectorObj}
	accoutsResp := model.CouchAccountList{}

	err = c.DB.Find(&accoutsResp,&params)
	if err != nil {
		return model.Account{}, err
	}
	if accoutsResp.Docs != nil && len(accoutsResp.Docs)>0 {
		return accoutsResp.Docs[0],nil
	}
	return model.Account{}, nil
}