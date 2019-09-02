package model

type BlockHeader struct {
	Number               uint64
	PreviousHash         string
	DataHash             string
}

type Block struct {
	Number               uint64
	PreviousHash         string
	DataHash             string
	TxList 	[]TransactionDetail
}

type TransactionDetail struct {
	TransactionId string
	CreateTime    string
	Args          []string
}

type Token struct {
	Action string `json:"action"`
	Amount float64 `json:"amount"`
	Desc   string `json:"desc"`
	Issuer string `json:"issuer"`
	Name   string `json:"name"`
	Status bool `json:"status"`
	Type   string `json:"type"`
}
type RecordToken struct {
	Key  string `json:"Key"`
	Record Token `json:"Record"`
}
type HistoryToken struct {
	TxId string `json:"TxId"`
	Value Token `json:"Value"`
}
type Payload struct {
	Status  int `json:"status"`
	Message string `json:"message"`
}

type Account struct {
	Address string `json:"address"`
	CN string `json:"cn"`
	Code string `json:"code"`
	MspId string `json:"mspid"`
	Name string `json:"name"`
	Status bool `json:"status"`
	Type string `json:"type"`
}

type CouchAccountList struct {
	Docs []Account `json:"docs"`
}