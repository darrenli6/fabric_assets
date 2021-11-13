package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

type HelloChaincode struct {
}

//链码初始化
func (h *HelloChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//初始化链码时调用，校验初始化参数
	args := stub.GetStringArgs()
	//校验参数
	if len(args) != 2 {
		return shim.Error("初始化链码时发生错误，参数必须为2个")
	}
	//更新世界状态
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error("更新世界状态失败")
	}
	fmt.Println("链码初始化成功")
	return shim.Success(nil)
}

//链码的执行
func (h *HelloChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetStringArgs()
	if len(args) != 1 {
		return shim.Error("传递参数必须为1个")
	}
	//查询
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("查询失败")
	}
	//返回查询结果
	return shim.Success(result)
}

func main() {
	err := shim.Start(new(HelloChaincode))
	if err != nil {
		fmt.Printf("start chaincode faild %s", err)
	}
}
