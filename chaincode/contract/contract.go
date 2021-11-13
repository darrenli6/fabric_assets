package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
)

//Invoke
//发布合约
//响应合约
//合约成交
//合约关闭
//合约交易查询
//合约交易历史数据查询

//定义合约结构体
//字段有的可能用不到
//用于发布合约
type Bill struct {
	//合约id
	TaskId string `json:"task_id"`
	//用户代码
	UserCode string `json:"user_code"`
	//合同代码
	ContractCode string `json:"contract_code"`
	//采购方用户代码
	PurchaserUserCode string `json:"purchaser_user_code"`
	//投标开始时间
	BidingStartTime string `json:"biding_start_time"`
	//投标结束时间
	BidingEndTime string `json:"biding_end_time"`
	//关闭时间
	CloseTime string `json:"close_time"`
	//合同开始时间
	ContractStartTime string `json:"contract_start_time"`
	//合同结束时间
	ContractEndTime string `json:"contract_end_time"`
	//合同状态
	//发布合同：yes  关闭合同: close  已成交: deal
	ContractStatus string `json:"contract_status"`
}

//定义用于响应合约的结构体
type BidingBill struct {
	//账单ID
	TaskId string `json:"task_id"`
	//用户代码
	UserCode string `json:"user_code"`
	//合同代码
	ContractCode string `json:"contract_code"`
	//合同状态
	//发布合同：yes  关闭合同: close  已成交: deal
	ContractStatus string `json:"contract_status"`
}

//合约成交
type BillDeal struct {
	//id
	TaskId string `json:"task_id"`
	//用户代码
	UserCode string `json:"user_code"`
	//合同代码
	ContractCode string `json:"contract_code"`
	//成交时间
	DealTime string `json:"deal_time"`
	//合同状态
	//发布合同：yes  关闭合同: close  已成交: deal
	ContractStatus string `json:"contract_status"`
}

//用于合约关闭
type BillClose struct {
	//id
	TaskId string `json:"task_id"`
	//用户代码
	UserCode string `json:"user_code"`
	//合同代码
	ContractCode string `json:"contract_code"`
	//结束时间
	CloseTime string `json:"close_time"`
	//合同状态
	//发布合同：yes  关闭合同: close  已成交: deal
	ContractStatus string `json:"contract_status"`
}

//链码的返回结构
type chaincodeRet struct {
	//1代表成功，0代表失败
	Result int `json:"result"`
	//0代表没有错误
	//1000代表参数错误
	//2000代表内容格式错误
	ErrorCode int `json:"error_code"`
	//错误信息
	ErrorMsg string `json:"error_msg"`
}

//合约查询
type QueryBill struct {
	TaskId       string `json:"task_id"`
	UserCode     string `json:"user_code"`
	ContractCode string `json:"contract_code"`
	//last：查询是的最新的合约交易，whole：查询的是全部的合约的交易信息
	VersionType string `json:"version_type"`
}

//定义用于返回查询结果的结构体
type queryRet struct {
	//1代表成功，0代表失败
	Result int `json:"result"`
	//0代表没有错误
	//1000代表参数错误
	//2000代表内容格式错误
	ErrorCode int `json:"error_code"`
	//错误信息
	ErrorMsg string `json:"error_msg"`
	//返回合同列表
	DataList []Bill `json"data_list"`
}

//结构体
type BillChaincode struct {
}

//初始化
func (a *BillChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

//链码入口
func (a *BillChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//接收方法名
	function, args := stub.GetFunctionAndParameters()
	//判断
	if function == "link_contract_create" {
		//发布合约
		return a.LinkContractCreate(stub, args)
	} else if function == "link_contract_biding" {
		//响应合约
		return a.LinkContractBiding(stub, args)
	} else if function == "link_contract_deal" {
		//合约成交
		return a.LinkContractDeal(stub, args)
	} else if function == "link_contract_close" {
		//合约关闭
		return a.LinkContractClose(stub, args)
	} else if function == "query" {
		//合约查询
		return a.query(stub, args)
	} else {
		//处理错误
		res := getRetString(0, 1000, "无效的方法名")
		return shim.Error(res)
	}
	return shim.Success(nil)
}

//根据传入的参数，处理异常
func getRetString(result int, code int, msg string) string {
	var r chaincodeRet
	//1代表成功，0代表失败
	r.Result = result
	//1000代表参数错误
	r.ErrorCode = code
	r.ErrorMsg = msg
	//序列化
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println("序列化失败")
		return ""
	}
	//返回字符串
	return string(b[:])
}

//根据传入的参数，处理异常
func getRetByte(result int, code int, msg string) []byte {
	var r chaincodeRet
	//1代表成功，0代表失败
	r.Result = result
	//1000代表参数错误
	r.ErrorCode = code
	r.ErrorMsg = msg
	//序列化
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println("序列化失败")
		return nil
	}
	return b
}

//发布合约
func (a *BillChaincode) LinkContractCreate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断参数个数
	if len(args) != 1 {
		//处理错误
		res := getRetString(0, 1000, "参数必须是一个")
		return shim.Error(res)
	}

	//将输入参数解析到结构体
	arg := []byte(args[0])
	bill := Bill{}
	err := json.Unmarshal(arg, &bill)
	if err != nil {
		//处理错误
		res := getRetString(0, 2000, "合约解析失败")
		return shim.Error(res)
	}
	//如果合约代码为空则返回错误
	if bill.ContractCode == "" {
		//处理错误
		res := getRetString(0, 2000, "合约代码不能为空")
		return shim.Error(res)
	}
	//进行校验，判断世界状态中合约是否已经存在
	_, existbl := a.getBill(stub, bill.ContractCode)
	//合约已经存在
	if existbl {
		//处理错误
		res := getRetString(0, 2000, "合约已经存在")
		return shim.Error(res)
	}

	//保存合约
	_, bl := a.putBill(stub, bill)
	if !bl {
		//保存失败
		//处理错误
		res := getRetString(0, 2000, "合约保存失败")
		return shim.Error(res)
	}
	//打印合约信息
	fmt.Println(bill)
	res := getRetByte(1, 0, "发布合约成功")
	return shim.Success(res)
}

//根据合约号取出合约
func (a *BillChaincode) getBill(stub shim.ChaincodeStubInterface, bill_No string) (Bill, bool) {
	var bill Bill
	key := bill_No
	//获取合约
	b, err := stub.GetState(key)
	if b == nil {
		return bill, false
	}
	err = json.Unmarshal(b, &bill)
	if err != nil {
		return bill, false
	}
	//返回
	return bill, true
}

//保存合约
func (a *BillChaincode) putBill(stub shim.ChaincodeStubInterface, bill Bill) ([]byte, bool) {
	//处理json
	byte, err := json.Marshal(bill)
	if err != nil {
		return nil, false
	}
	//保存
	err = stub.PutState(bill.ContractCode, byte)
	if err != nil {
		return nil, false
	}
	return byte, true
}

//响应合约
func (a *BillChaincode) LinkContractBiding(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断参数个数
	if len(args) != 1 {
		res := getRetString(0, 1000, "参数必须是一个")
		return shim.Error(res)
	}

	//解析
	arg := []byte(args[0])
	biding_bill := &BidingBill{}
	err := json.Unmarshal(arg, biding_bill)
	if err != nil {
		res := getRetString(0, 2000, "json解析失败")
		return shim.Error(res)
	}

	//判断合约是否存在
	key_id := biding_bill.ContractCode
	bill, bl := a.getBill(stub, key_id)
	if !bl {
		//没查到
		res := getRetString(0, 1000, "合同代码不存在")
		return shim.Error(res)
	}

	//判断合约是否已经成交，若成交，不能再响应
	if bill.ContractStatus == "deal" {
		res := getRetString(0, 2000, "合约已经成交")
		return shim.Error(res)
	} else if bill.ContractStatus == "close" {
		res := getRetString(0, 2000, "合约已经关闭")
		return shim.Error(res)
	}
	fmt.Println(bill)

	//保存
	_, bl = a.putBill(stub, bill)
	//
	bill, bl = a.getBill(stub, key_id)
	fmt.Println(bill)
	if !bl {
		res := getRetString(0, 2000, "合约保存失败")
		return shim.Error(res)
	}

	res := getRetByte(1, 0, "响应成功")
	return shim.Success(res)
}

//合约成交
func (a *BillChaincode) LinkContractDeal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断输入参数个数
	if len(args) != 1 {
		res := getRetString(0, 1000, "参数必须是一个")
		return shim.Error(res)
	}

	//解析
	arg := []byte(args[0])
	//
	billdeal := &BillDeal{}
	err := json.Unmarshal(arg, billdeal)
	if err != nil {
		res := getRetString(0, 1000, "解析失败")
		return shim.Error(res)
	}

	//判断合约是否存在
	bill, existbl := a.getBill(stub, billdeal.ContractCode)
	if !existbl {
		res := getRetString(0, 1000, "合约不存在")
		return shim.Error(res)
	}

	//判断合约是否已经关闭或成交
	if bill.ContractStatus == "deal" {
		res := getRetString(0, 2000, "合约已经成交")
		return shim.Error(res)
	} else if bill.ContractStatus == "close" {
		res := getRetString(0, 2000, "合约已经关闭")
		return shim.Error(res)
	}

	//修改合约状态
	bill.ContractStatus = "deal"
	//保存合约
	_, bl := a.putBill(stub, bill)
	if !bl {
		res := getRetString(0, 2000, "合约保存失败")
		return shim.Error(res)
	}
	fmt.Println(bill)

	res := getRetByte(1, 0, "合约成交成功")
	return shim.Success(res)
}

//合约关闭
func (a *BillChaincode) LinkContractClose(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断输入参数个数
	if len(args) != 1 {
		res := getRetString(0, 1000, "参数必须是一个")
		return shim.Error(res)
	}

	//解析
	arg := []byte(args[0])
	//
	billclose := &BillClose{}
	err := json.Unmarshal(arg, billclose)
	if err != nil {
		res := getRetString(0, 1000, "解析失败")
		return shim.Error(res)
	}

	//合约关闭
	err = stub.PutState(billclose.ContractCode, arg)
	if err != nil {
		return shim.Error("合约关闭失败")
	}

	//查找合约是否存在
	bill, exitbl := a.getBill(stub, billclose.ContractCode)
	if !exitbl {
		res := getRetString(0, 1000, "合约关闭失败，合约不存在")
		return shim.Error(res)
	}

	//更改合约状态和时间
	bill.ContractStatus = billclose.ContractStatus
	bill.CloseTime = billclose.CloseTime
	_, bl := a.putBill(stub, bill)
	if !bl {
		res := getRetString(0, 1000, "合约关闭失败")
		return shim.Error(res)
	}

	res := getRetByte(1, 0, "合约关闭成功")
	return shim.Success(res)
}

//查询合约
//支持查询最新的，和查询所有的
func (a *BillChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断输入参数个数
	if len(args) != 1 {
		res := getRetString(0, 1000, "参数必须是一个")
		return shim.Error(res)
	}

	//解析
	query_bill := &QueryBill{}
	err := json.Unmarshal([]byte(args[0]), query_bill)
	if err != nil {
		res := getRetString(0, 1000, "解析失败")
		return shim.Error(res)
	}

	fmt.Println(query_bill)

	if query_bill.VersionType == "last" {
		//查询最新合约信息
		return a.queryLastBill(stub, query_bill.UserCode, query_bill.ContractCode)
	} else {
		//查询所有合约信息
		return a.queryWholeBill(stub, query_bill.UserCode, query_bill.ContractCode)
	}
}

//查询最新合约信息
func (a *BillChaincode) queryLastBill(stub shim.ChaincodeStubInterface, userCode string, contractCode string) pb.Response {
	key_id := contractCode
	//查询合约
	bill, bl := a.getBill(stub, key_id)
	if !bl {
		res := getRetString(0, 1000, "查询失败")
		return shim.Error(res)
	}

	var history []Bill
	var hist = bill
	history = append(history, hist)
	res := getQueryByte(1, 0, "", history)
	return shim.Success(res)
}

//查询所有合约信息
func (a *BillChaincode) queryWholeBill(stub shim.ChaincodeStubInterface, userCode string, contractCode string) pb.Response {
	key_id := contractCode
	//fabric的API查询历史
	resultsIterator, err := stub.GetHistoryForKey(key_id)
	if err != nil {
		res := getRetString(0, 1000, "查询合约历史失败")
		return shim.Error(res)
	}
	defer resultsIterator.Close()

	var history []Bill
	//循环遍历
	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			res := getRetString(0, 1000, "遍历合约历史失败")
			return shim.Error(res)
		}
		//定义合约结构
		var hist Bill
		if historyData.Value == nil {
			var emptyBill Bill
			hist = emptyBill
		} else {
			json.Unmarshal(historyData.Value, &hist)
			fmt.Println(hist)
		}
		history = append(history, hist)
		fmt.Println(history)
	}

	res := getQueryByte(1, 0, "", history)
	return shim.Success(res)
}

//根据传的状态码，返回查询的字节数组
func getQueryByte(result int, code int, msg string, hist []Bill) []byte {
	var r queryRet
	//组装数据
	r.Result = result
	r.ErrorCode = code
	r.ErrorMsg = msg
	r.DataList = hist
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println("序列化失败")
		return nil
	}
	return b
}

func main() {
	if err := shim.Start(new(BillChaincode)); err != nil {
		fmt.Printf("启动链码失败:%s", err)
	}
}
