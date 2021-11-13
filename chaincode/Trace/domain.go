package main

//银行、账户、定义交易历史
//定义银行
type Bank struct {
	//名字
	BankName string `json:"BankName"`
	//金额
	Amount int `json:"Amount"`
	//1：贷款   2：还款
	Flag int `json:"Flag"`
	//起始时间
	StartTime string `json:"StartTime"`
	//结束时间
	EndTime string `json:"EndTime"`
}

//定义账户
type Account struct {
	//身份证号
	CardNo string `json:"CardNo"`
	//用户名
	Aname string `json:"Aname"`
	//性别
	Gender string `json:"Gender"`
	//电话
	Mobile string `json:"Mobile"`
	//银行
	Bank Bank `json:"Bank"`
	//交易历史
	Histroys []HistoryItem
}

//定义交易历史
type HistoryItem struct {
	//交易id
	TxID string
	//账户
	Account Account
}
