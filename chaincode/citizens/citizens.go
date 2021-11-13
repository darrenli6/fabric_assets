package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"log"
	"encoding/json"
)

//个人基本信息
type People struct {
	//区分数据类型
	DataType string `json:"dataType"`
	//身份证号码
	Id string `json:"id"`
	//性别
	Sex string `json:"sex"`
	//姓名
	Name string `json:"name"`
	//出生地
	BrithLocation Location `json:"birthLocation"`
	//现居住地
	LiveLocation Location `json:"liveLocation"`
	//母亲身份证号
	MotherId string `json:"montherID"`
	//父亲身份证号
	FatherId string `json:"fatherID"`
}

//位置
type Location struct {
	//国家
	Country string `json:"country"`
	//省
	Province string `json:"province"`
	//城市
	City string `json:"city"`
	//镇
	Town string `json:"town"`
	//详细地址
	Detail string `json:"detail"`
}

//公民链
type CitizensChain struct {
}

//初始化方法
func (c *CitizensChain) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, _ := stub.GetFunctionAndParameters()
	if function != "init" {
		return shim.Error("方法名错误")
	}
	log.Println("初始化成功")
	return shim.Success(nil)
}

//链码交互入口
func (c *CitizensChain) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//接收参数判断
	function, args := stub.GetFunctionAndParameters()
	//判断
	if function == "register" {
		//录入公民信息
		return c.register(stub, args)
	} else if function == "query" {
		//查询公民信息
		return c.query(stub, args)
	} else {
		return shim.Error("无效的方法名")
	}
	return shim.Success(nil)
}

//录入公民信息
//-c Args[register,身份证号,json]
//参数1：身份证号，是存储的key
//参数2：个人信息，当成value取存储
func (c *CitizensChain) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断参数个数
	if len(args) != 2 {
		return shim.Error("参数错误")
	}
	//接收身份证号
	key := args[0]
	//接收公民信息（json）
	value := args[1]
	perple := People{}
	//转换
	err := json.Unmarshal([]byte(value), &perple)
	if err != nil {
		return shim.Error("注册失败,参数无法解析")
	}
	//更新世界状态
	stub.PutState(key, []byte(value))
	return shim.Success(nil)
}

//查询公民信息
func (c *CitizensChain) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("参数错误")
	}
	//接收查询的key
	key := args[0]
	//去世界状态中查数据
	result, err := stub.GetState(key)
	if err != nil {
		return shim.Error("查询失败")
	}
	return shim.Success(result)
}

func main() {
	err := shim.Start(new(CitizensChain))
	if err != nil {
		log.Println(err)
	}
}
