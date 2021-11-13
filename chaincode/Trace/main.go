package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
)

type TraceChaincode struct {
}

//初始化方法
func (t *TraceChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	//初始化测试数据
	initTest(stub)
	return shim.Success(nil)
}

//链码入口invoke
//loan：贷款
//repayment：还款
//initTest：测试初始化
func (t *TraceChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//得到方法名和参数
	fun, args := stub.GetFunctionAndParameters()
	//进行判断
	if fun == "loan" {
		//贷款
		return loan(stub, args)
	} else if fun == "repayment" {
		//还款
		return repayment(stub, args)
	} else if fun == "initTest" {
		return initTest(stub)
	} else {
		return shim.Error("方法名错误")
	}
}

//测试方法
func initTest(stub shim.ChaincodeStubInterface) peer.Response {
	bank := Bank{
		BankName: "icbc",
		Amount:   6000,
		//贷款
		Flag:      1,
		StartTime: "2010-01-01",
		EndTime:   "2020-01-01",
	}
	bank1 := Bank{
		BankName: "abc",
		Amount:   1000,
		//还款
		Flag:      2,
		StartTime: "2010-02-01",
		EndTime:   "2020-02-01",
	}
	account := Account{
		CardNo:   "1234",
		Aname:    "jack",
		Gender:   "男",
		Mobile:   "15900000",
		Bank:     bank1,
		Histroys: nil,
	}
	account1 := Account{
		CardNo:   "123444",
		Aname:    "jack2",
		Gender:   "男",
		Mobile:   "1590000000",
		Bank:     bank,
		Histroys: nil,
	}

	//对象序列化，存储
	accBytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error("序列化账户失败")
	}
	accBytes1, err := json.Marshal(account1)
	if err != nil {
		return shim.Error("序列化账户1失败")
	}

	//保存
	err = stub.PutState(account.CardNo, accBytes)
	if err != nil {
		return shim.Error("保存账户失败")
	}
	err = stub.PutState(account.CardNo, accBytes1)
	if err != nil {
		return shim.Error("保存账户1失败")
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(TraceChaincode))
	if err != nil {
		fmt.Println("启动链码时发生错误")
	}
}
