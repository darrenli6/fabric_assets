package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"fmt"
)

type PaymentChaincode struct {
}

//初始化方法
//-c '{"Args":["init","第一个账户名","第一个账户余额",XXXXXX]}'
func (p *PaymentChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 4 {
		return shim.Error("必须是4个参数")
	}
	var err error
	//拿到参数转换
	_, err = GetArgsState(args[1])
	if err != nil {
		return shim.Error("第一个账户的金额错误")
	}
	_, err = GetArgsState(args[3])
	if err != nil {
		return shim.Error("第二个账户的金额错误")
	}
	//将初始化数据存到账本中
	//持久化
	err = stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error("第一个账户保存失败")
	}
	err = stub.PutState(args[2], []byte(args[3]))
	if err != nil {
		return shim.Error("第二个账户保存失败")
	}
	fmt.Println("初始化成功")
	return shim.Success(nil)

}

//参数转换
//字符串转数字
func GetArgsState(value string) (int, error) {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return v, err
}

//链码交互的入口
func (p *PaymentChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//接收方法名
	fun, args := stub.GetFunctionAndParameters()
	//判断方法入口
	if fun == "query" {
		return query(stub, args)
	} else if fun == "invoke" {
		return invoke(stub, args)
	} else if fun == "set" {
		return set(stub, args)
	} else if fun == "get" {
		return get(stub, args)
	} else {
		return shim.Error("方法名错误")
	}
	return shim.Success(nil)
}

//根据指定账户查询
func query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("必须指定一个要查询的账户")
	}
	//查询操作
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据指定账户查询失败")
	}
	if result == nil {
		return shim.Error("没有查到数据")
	}
	return shim.Success(result)
}

//转账
//-c '{"Args":["invoke","原账户","目标账户","转账金额"]}'
func invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断参数
	if len(args) != 3 {
		return shim.Error("参数个数错误")
	}
	//判断金额
	v, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("转账金额错误，请重新设置")
	}
	//判断原账户有没有钱
	v1, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("查询失败")
	}
	if v1 == nil {
		return shim.Error("原账户未查询到数据")
	}
	//将查询的数据转换
	v2, err := strconv.Atoi(string((v1)))
	if err != nil {
		return shim.Error("金额类似错误")
	}
	if v2 < v {
		return shim.Error("原账户的金额不够")
	} else {
		v2 = v2 - v
	}

	var tarv int

	//查询目标账户余额
	bv, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error("目标账户查询错误")
	}
	if bv == nil {
		err = stub.PutState(args[1], []byte(args[2]))
		if err != nil {
			return shim.Error("目标账户更新错误")
		}
	} else {
		tarv, err = strconv.Atoi(string(bv))
		if err != nil {
			return shim.Error("转账金额类型错误")
		}
		//目标账户余额=目标账户余额+转账过来的钱
		tarv = tarv + v
	}
	//更新目标账户的金额
	err = stub.PutState(args[1], []byte(strconv.Itoa(tarv)))
	if err != nil {
		return shim.Error("更新目标账户失败")
	}
	//更新账户余额
	err = stub.PutState(args[0], []byte(strconv.Itoa(v2)))
	if err != nil {
		return shim.Error("更新原账户失败")
	}
	return shim.Success([]byte("转账成功"))
}

//向指定账户存钱
//-c Args [set,目标账户，金额]
func set(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("参数个数错误")
	}
	//判断金额是否类型正确
	v, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("金额类型出错")
	}
	//查询金额
	result, err := getValueWithKey(stub, args[0])
	if err != nil {
		return shim.Error("查询错误")
	}

	//转换
	v1, err := transByteToInt(result)
	if err != nil {
		return shim.Error("转换失败")
	}
	//转账后的金额
	v = v + v1
	//将新的金额存到账户中
	err = stub.PutState(args[0], []byte(strconv.Itoa(v)))
	if err != nil {
		return shim.Error("更新失败")
	}
	return shim.Success([]byte("保存成功"))
}

//账户查询
func getValueWithKey(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	result, err := stub.GetState(key)
	if err != nil {
		return nil, fmt.Errorf("查询账户出错")
	}
	if result == nil {
		return nil, fmt.Errorf("账户未查询到")
	}
	return result, nil
}

//将[]byte转int
func transByteToInt(value []byte) (int, error) {
	v, err := strconv.Atoi(string(value))
	if err != nil {
		return 0, fmt.Errorf("转换格式失败")
	}
	return v, nil
}

//取钱
// -c Args [get,目标账户，金额]
func get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("参数个数错误")
	}
	//取值
	result, err := getValueWithKey(stub, args[0])
	if err != nil {
		return shim.Error("查询失败")
	}
	//转换金额
	v1, err := transByteToInt([]byte(args[1]))
	if err != nil {
		return shim.Error("转换失败")
	}
	v, err := transByteToInt(result)
	if err != nil {
		return shim.Error("转换失败")
	}
	if v < v1 {
		//账户余额不足
		return shim.Error("余额不足")
	} else {
		//账户余额=账户余额-取的钱
		v = v - v1
	}

	//保存世界状态
	err = stub.PutState(args[0], []byte(strconv.Itoa(v)))
	if err != nil {
		return shim.Error("更新错误")
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(PaymentChaincode))
	if err != nil {
		fmt.Println("启动链码失败")
	}
}
