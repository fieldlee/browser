package utils

import (
	"browser/model"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type SqlCliet struct {
	DB *sql.DB
}

func InitSql()(*SqlCliet,error) {
	ymlCon := InitConfig()
	dbhost := ymlCon.V.GetString("dbhost")
	dbuser := ymlCon.V.GetString("dbuser")
	dbpwd  := ymlCon.V.GetString("dbpwd")
	dbname := ymlCon.V.GetString("dbname")
	var err error
	var db *sql.DB
	db,err = sql.Open("mysql",dbuser+":"+dbpwd+"@tcp("+dbhost+")/"+dbname+"?charset=utf8&parseTime=true")
	if err != nil {
		return nil,err
	}
	err = db.Ping()
	if err != nil {
		return nil,err
	}
	return &SqlCliet{
		DB:db,
	},nil
}

func (s *SqlCliet)InsertBlock(block model.BlockHeader)error{
	stmt, err := s.DB.Prepare("INSERT INTO blocks(height,createtime,prehash,hash) VALUES (?,?,?,?) ")
	defer stmt.Close()
	if err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02 15:04:05",block.CreateTime)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(block.Number,t,block.PreviousHash,block.DataHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)UpdateBlockHash(height int,hash string)error{
	stmt, err := s.DB.Prepare("update blocks set hash = ? where height = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(hash,height)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)UpdateTxHash(curhash string,hash string)error{
	stmt, err := s.DB.Prepare("update transactions set blockhash = ? where blockhash = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(curhash,hash)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)InsertTx(blockhash string,tx model.TransactionDetail)error{
	stmt, err := s.DB.Prepare("INSERT INTO transactions(txhash,blockhash,createtime,method,args,signed) VALUES(?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02 15:04:05",tx.CreateTime )
	if err != nil {
		return err
	}

	method , args , signed  := "","",""

	if tx.Args != nil &&  tx.Args[0] != "" {
		method = tx.Args[0]
	}
	if len(tx.Args) > 1 {
		if tx.Args != nil &&  tx.Args[1] != "" {
			args = tx.Args[1]
		}
	}
	if len(tx.Args) > 2 {
		if tx.Args != nil &&  tx.Args[2] != "" {
			signed = tx.Args[2]
		}
	}

	_, err = stmt.Exec(tx.TransactionId,blockhash,t,method,args,signed)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqlCliet)QueryBlockHeight()(int,error){
	stmt,err := s.DB.Prepare("select height from blocks order by height desc limit 1")
	defer stmt.Close()
	if err != nil {
		return 0,err
	}
	row := stmt.QueryRow()
	if row != nil {
		height := 0
		err = row.Scan(&height)
		if err != nil {
			fmt.Printf(err.Error())
			return 0,err
		}
		return height,nil
	}else{
		return 0 , errors.New("no row in the table")
	}
}

func (s *SqlCliet)QueryBlockByHeight(height int)(model.BlockHeader,error){
	stmt,err := s.DB.Prepare("select prehash,hash,height,createtime from blocks where height = ?")
	defer stmt.Close()
	if err != nil {
		return model.BlockHeader{} , err
	}
	row := stmt.QueryRow(height)
	if err != nil {
		return model.BlockHeader{} , err
	}
	if row != nil {
		blockheader := model.BlockHeader{}
		var time = time.Now()
		err = row.Scan(&blockheader.PreviousHash, &blockheader.DataHash, &blockheader.Number,&time)
		if err != nil {
			fmt.Printf(err.Error())
			return model.BlockHeader{} , err
		}
		blockheader.CreateTime = time.Format("2006-01-02 15:04:05")
		return blockheader,nil
	}else{
		return model.BlockHeader{} , errors.New("the transaction not exist")
	}
}

func (s *SqlCliet)QueryBlockByHash(hash string)(model.BlockHeader,error){
	stmt,err := s.DB.Prepare("select prehash,hash,height,createtime from blocks where hash = ?")
	defer stmt.Close()
	if err != nil {
		return model.BlockHeader{} , err
	}
	row := stmt.QueryRow(hash)
	if err != nil {
		return model.BlockHeader{} , err
	}
	if row != nil {
		blockheader := model.BlockHeader{}
		var time = time.Now()
		err = row.Scan(&blockheader.PreviousHash, &blockheader.DataHash, &blockheader.Number,&time)
		if err != nil {
			fmt.Printf(err.Error())
			return model.BlockHeader{} , err
		}
		blockheader.CreateTime = time.Format("2006-01-02 15:04:05")
		return blockheader,nil
	}else{
		return model.BlockHeader{} , errors.New("the transaction not exist")
	}
}

func (s *SqlCliet)QueryBlocksByRange(curHeight int,limit int)([]model.BlockHeader,error){
	start := curHeight - limit
	if start <= 0{
		start = 0
	}
	stmt,err := s.DB.Prepare("select prehash,hash,height,createtime from blocks where height >= ? and height <= ? order by height desc ")
	defer stmt.Close()
	if err != nil {
		return nil , err
	}
	rows,err := stmt.Query(start,curHeight)
	if err != nil {
		return nil , err
	}
	listBLOCK := make([]model.BlockHeader,0)
	for rows.Next(){
		blockheader := model.BlockHeader{}
		var time = time.Now()
		err = rows.Scan(&blockheader.PreviousHash, &blockheader.DataHash, &blockheader.Number,&time)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		blockheader.CreateTime = time.Format("2006-01-02 15:04:05")

		listBLOCK = append(listBLOCK,blockheader)
	}
	return listBLOCK,nil
}

func (s *SqlCliet)QueryTxsByBlockHash(hash string)([]model.TransactionDetail,error){
	stmt,err := s.DB.Prepare("select txhash,method,args,signed,createtime  from transactions where blockhash = ?")
	defer stmt.Close()
	if err != nil {
		return nil , err
	}
	rows,err := stmt.Query(hash)
	if err != nil {
		return nil , err
	}
	listTX := make([]model.TransactionDetail,0)
	for rows.Next(){
		tx := model.TransactionDetail{}
		var time = time.Now()
		var method = ""
		var args = ""
		var signed = ""

		err = rows.Scan(&tx.TransactionId, &method, &args,&signed,&time)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		tx.CreateTime = time.Format("2006-01-02 15:04:05")
		var argslist = make([]string,0)
		argslist = append(argslist,method)
		argslist = append(argslist,args)
		argslist = append(argslist,signed)
		tx.Args = argslist
		listTX = append(listTX,tx)
	}
	return listTX,nil
}

func (s *SqlCliet)QueryBlockHashByTxId(hash string)(string,error){
	stmt,err := s.DB.Prepare("select blockhash  from transactions where txhash = ?")
	defer stmt.Close()
	if err != nil {
		return "" , err
	}
	row := stmt.QueryRow(hash)
	if row != nil {
		var blockhash = new(string)
		err = row.Scan(&blockhash)
		if err != nil {
			return "" , err
		}
		return *blockhash ,nil
	}else{
		return "" , errors.New("the hash transaction not exist")
	}
}

func (s *SqlCliet)QueryTxs(hash string)(model.TransactionDetail,error){
	stmt,err := s.DB.Prepare("select txhash,method,args,signed,createtime  from transactions where txhash = ?")
	defer stmt.Close()
	if err != nil {
		return model.TransactionDetail{} , err
	}
	row := stmt.QueryRow(hash)

	if row != nil {
		tx := model.TransactionDetail{}
		var time = time.Now()
		var method = ""
		var args = ""
		var signed = ""

		err = row.Scan(&tx.TransactionId, &method, &args,&signed,&time)
		if err != nil {
			return model.TransactionDetail{} , err
		}
		tx.CreateTime = time.Format("2006-01-02 15:04:05")
		var argslist = make([]string,0)
		argslist = append(argslist,method)
		argslist = append(argslist,args)
		argslist = append(argslist,signed)
		tx.Args = argslist
		return tx ,nil
	}else{
		return model.TransactionDetail{} , errors.New("the hash transaction not exist")
	}
}

func (s *SqlCliet)QueryTxsByAccount(account string)([]model.TransactionDetail,error){
	stmt,err := s.DB.Prepare("select txhash,method,args,signed,createtime  from transactions where method = 'transfer' and args like ?")
	defer stmt.Close()
	if err != nil {
		return nil , err
	}
	laccount := "%"+account+"%"
	rows,err := stmt.Query(laccount)
	if err != nil {
		return nil , err
	}
	listTX := make([]model.TransactionDetail,0)
	for rows.Next(){
		tx := model.TransactionDetail{}
		var time = time.Now()
		var method = ""
		var args = ""
		var signed = ""

		err = rows.Scan(&tx.TransactionId, &method, &args,&signed,&time)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		tx.CreateTime = time.Format("2006-01-02 15:04:05")
		var argslist = make([]string,0)
		argslist = append(argslist,method)
		argslist = append(argslist,args)
		argslist = append(argslist,signed)
		tx.Args = argslist
		listTX = append(listTX,tx)
	}
	return listTX,nil
}


func (s *SqlCliet)QueryTxsByToken(token string)([]model.TransactionDetail,error){
	stmt,err := s.DB.Prepare("select txhash,method,args,signed,createtime  from transactions where method = 'transfer' and args like ?")
	defer stmt.Close()
	if err != nil {
		return nil , err
	}
	laccount := "%"+token+"%"
	rows,err := stmt.Query(laccount)
	if err != nil {
		return nil , err
	}
	listTX := make([]model.TransactionDetail,0)
	for rows.Next(){
		tx := model.TransactionDetail{}
		var time = time.Now()
		var method = ""
		var args = ""
		var signed = ""

		err = rows.Scan(&tx.TransactionId, &method, &args,&signed,&time)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		tx.CreateTime = time.Format("2006-01-02 15:04:05")
		var argslist = make([]string,0)
		argslist = append(argslist,method)
		argslist = append(argslist,args)
		argslist = append(argslist,signed)
		tx.Args = argslist
		listTX = append(listTX,tx)
	}
	return listTX,nil
}

func (s *SqlCliet)QueryTxsNum()(int,error){
	stmt,err := s.DB.Prepare(" select count(*) as txcount  from transactions")
	defer stmt.Close()
	if err != nil {
		return 0 , err
	}
	row := stmt.QueryRow()
	var count = 0
	err = row.Scan(&count)
	if err != nil {
		return 0 , err
	}
	return count,nil
}


func (s *SqlCliet)InsertToken(token model.Token)error{
	stmt, err := s.DB.Prepare("INSERT INTO tokens (name_,amount,issuer,status,type_,action_,desc_) VALUES (?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(token.Name,float64(token.Amount),token.Issuer,token.Status,token.Type,token.Action,token.Desc)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)RemoveToken()error{
	stmt, err := s.DB.Prepare("delete from tokens")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)QueryTokensById(token string)([]model.Token,error){
	stmt,err := s.DB.Prepare("select name_,amount,issuer,status,type_,action_,desc_  from tokens where name_ like ?")
	defer stmt.Close()
	if err != nil {
		return nil , err
	}
	tokenlike := "%"+token+"%"
	rows,err := stmt.Query(tokenlike)
	if err != nil {
		return nil , err
	}
	tokenList := make([]model.Token,0)
	for rows.Next(){
		name := ""
		amount := float64(0)
		issuer := ""
		status := false
		type_ := ""
		action := ""
		desc := ""
		err = rows.Scan(&name,&amount,&issuer,&status,&type_,&action,&desc)
		if err != nil {
			continue
		}
		tmpToken := model.Token{
			Amount:amount,
			Issuer:issuer,
			Name:name,
			Type:type_,
			Status:status,
			Action:action,
			Desc:desc,
		}
		tokenList = append(tokenList,tmpToken)
	}
	return tokenList,nil
}

func (s *SqlCliet)QueryTokens()([]model.Token,error){
	stmt,err := s.DB.Prepare("select name_,amount,issuer,status,type_,action_,desc_  from tokens ")
	defer stmt.Close()
	if err != nil {
		return nil , err
	}
	rows,err := stmt.Query()
	if err != nil {
		return nil , err
	}
	tokenList := make([]model.Token,0)
	for rows.Next(){
		name := ""
		amount := float64(0)
		issuer := ""
		status := false
		type_ := ""
		action := ""
		desc := ""
		err = rows.Scan(&name,&amount,&issuer,&status,&type_,&action,&desc)
		if err != nil {
			continue
		}
		tmpToken := model.Token{
			Amount:amount,
			Issuer:issuer,
			Name:name,
			Type:type_,
			Status:status,
			Action:action,
			Desc:desc,
		}
		tokenList = append(tokenList,tmpToken)
	}

	return tokenList,nil
}

func (s *SqlCliet)CloseSql(){
	s.DB.Close()
}
///CREATE DATABASE IF NOT EXISTS mmchannel DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
///CREATE TABLE blocks( height INT NOT NULL ,createtime DATETIME,  prehash VARCHAR(100) NOT NULL, hash VARCHAR(100) NOT NULL, PRIMARY KEY ( height ) )ENGINE=InnoDB DEFAULT CHARSET=utf8;
///CREATE TABLE transactions( txhash VARCHAR(100) NOT NULL ,  blockhash VARCHAR(100) NOT NULL, method VARCHAR(50) DEFAULT NULL, args VARCHAR(150) DEFAULT NULL, signed VARCHAR(200) DEFAULT NULL, createtime DATETIME, PRIMARY KEY ( txhash ) )ENGINE=InnoDB DEFAULT CHARSET=utf8;
///CREATE TABLE tokens( name_ VARCHAR(30) NOT NULL,amount float, issuer VARCHAR(30) DEFAULT NULL,action_ VARCHAR(30) DEFAULT NULL, desc_ VARCHAR(50) DEFAULT NULL,status bool, type_ VARCHAR(20) DEFAULT NULL, PRIMARY KEY ( name_) )ENGINE=InnoDB DEFAULT CHARSET=utf8;
