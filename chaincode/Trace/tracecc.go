package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"encoding/json"
)

//定义贷款和还款
const (
	Bank_Flag_Loan      = 1
	Bank_Flag_Repayment = 2
)

//贷款
//-c Args:[loan,账户身份证号,银行名字,金额]
func loan(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//判断参数
	if len(args) != 3 {
		return shim.Error("参数个数错误")
	}
	//判断类型
	v, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("类型错误")
	}

	//组装数据
	bank := Bank{
		BankName:  args[1],
		Amount:    v,
		Flag:      Bank_Flag_Loan,
		StartTime: "2011-01-10",
		EndTime:   "2021-01-09",
	}
	account := Account{
		CardNo: args[0],
		Aname:  "jack",
		Gender: "男",
		Mobile: "159999999",
		Bank:   bank,
	}

	//保存状态
	b := putAccount(stub, account)
	//判断
	if ! b {
		return shim.Error("保存贷款数据失败")
	}
	return shim.Success([]byte("保存贷款数据成功"))
}

//序列化并保存
//参数将要保存的账户传过来，返回布尔
func putAccount(stub shim.ChaincodeStubInterface, account Account) bool {
	//序列化
	accBytes, err := json.Marshal(account)
	if err != nil {
		return false
	}
	//保存数据
	err = stub.PutState(account.CardNo, accBytes)
	if err != nil {
		return false
	}
	return true
}

//还款
//-c Args:[repayment,账户身份证号,银行名字,金额]
func repayment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("参数个数错误")
	}
	v, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("指定的还款金额错误")
	}
	//组装数据
	bank := Bank{
		BankName:  args[1],
		Amount:    v,
		Flag:      Bank_Flag_Repayment,
		StartTime: "2011-01-10",
		EndTime:   "2021-01-09",
	}
	account := Account{
		CardNo: args[0],
		Aname:  "jack",
		Gender: "男",
		Mobile: "159999999",
		Bank:   bank,
	}
	b := putAccount(stub, account)
	if !b {
		return shim.Error("存款失败")
	}
	return shim.Success([]byte("存款成功"))
}
