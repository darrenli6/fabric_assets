package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"golang.org/x/net/context"
	"time"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"net/http"
	"bytes"
)

//定义SDK、channel等全局变量
var (
	sdk *fabsdk.FabricSDK


	//通道名称
	channelName = "assetschannel"
	//链码名称
	chaincodeName = "assets"
	//组织
	org = "org1"
	//用户
	user = "Admin"
	//配置文件位置，是当前位置下的配置
	configPath = "./config.yaml"
)

//初始化sdk

func init(){
	var err error
	fabsdk.New(config.FromFile(configPath))

	if err!=nil{
		panic(err)
	}


}


func main(){

	//使用http服务
	//使用了gin框架，是go-web框架
	router := gin.Default()
	//定义http服务的路由
	{
		//注册用户
		router.POST("/users", userRegister)
		//查询用户，根据id
		router.GET("/users/:id", queryUser)
		//删除用户,根据id
		router.DELETE("/users/:id", deleteUser)
		//查询资产
		router.GET("/assets/get/:id", queryAsset)
		//资产登记
		router.POST("/assets/enroll", assetsEnroll)
		//资产变更
		router.POST("/assets/exchange", assetsExchange)
		//资产变更历史查询
		router.GET("/assets/exchange/history", assetsExchangeHistory)
	}
	router.Run()

}


//SDK有4个模块，这里只用到了区块链交互
//区块链管理
//区块链数据查询
//区块链交互
//事件监听


//区块链管理
func manageBolockchain() {
	//表明身份
	//向SDK说明身份
	//获取上下文对象
	ctx := sdk.Context(fabsdk.WithOrg(org), fabsdk.WithUser(user))

	//区块链管理的方法在fabric-sdk-go\pkg\client\resmgmt下
	//创建Client对象
	cli, err := resmgmt.New(ctx)
	if err != nil {
		panic(err)
	}


	//具体操作
	cli.SaveChannel(resmgmt.SaveChannelRequest{},
		resmgmt.WithOrdererEndpoint("orderer.example.com"),
		resmgmt.WithTargetEndpoints())

}


//区块链查询
//在fabric-sdk-go\pkg\client\ledger包下
func queryBlockchain() {
	//获取上下文对象
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))
	//创建Client对象
	cli, err := ledger.New(ctx)
	if err != nil {
		panic(err)
	}

	//查询区块链当前状态
	resp, err := cli.QueryInfo(ledger.WithTargetEndpoints("peer0,org1.example.com"))
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

	//查询区块
	//方法1：
	cli.QueryBlockByHash(resp.BCI.CurrentBlockHash)
	//方法2：
	for i := uint64(0); i <= resp.BCI.Height; i++ {
		cli.QueryBlock(i)
	}
}

//事件监听
func eventHandle() {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))
	//创建Client对象
	//源码在fabric-sdk-go\pkg\client\channel下
	cli, err := event.New(ctx)
	if err != nil {
		panic(err)
	}

	//注册区块事件
	reg, blkevent, err := cli.RegisterBlockEvent()
	if err != nil {
		panic(err)
	}
	defer cli.Unregister(reg)
	timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for {
		select {
		case evt := <-blkevent:
			fmt.Printf("接收了一个块", evt)
		case <-timeoutctx.Done():
			fmt.Println("事件超时")
			return
		}
	}

}



//定义Execute，处理后面的注册、注销、资产相关
func channelExecute(fcn string, args [][]byte) (channel.Response, error) {
	//上下文对象
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))
	//创建Client对象
	//源码在fabric-sdk-go\pkg\client\channel下
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	//Execute()：更新状态，增删改
	//参数指定了请求内容和交易发送到的节点
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.example.com"))
	if err != nil {
		return channel.Response{}, err
	}

	//监听
	go func() {
		//使用channel模块的事件监听
		reg, ccevt, err := cli.RegisterChaincodeEvent(chaincodeName, "eventname")
		if err != nil {
			return
		}
		defer cli.UnregisterChaincodeEvent(reg)
		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		for {
			select {
			case evt := <-ccevt:
				fmt.Printf("接收到事件 %s: %+v", resp.TransactionID, evt)
			case <-timeoutctx.Done():
				fmt.Println("事件超时")
				return
			}
		}
	}()

	//交易状态事件监听
	go func() {
		eventcli, err := event.New(ctx)
		if err != nil {
			return
		}
		//注册交易事件
		reg, status, err := eventcli.RegisterTxStatusEvent(string(resp.TransactionID))
		defer eventcli.Unregister(reg)

		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		for {
			select {
			case evt := <-status:
				fmt.Printf("接收到事件 %s: %+v", resp.TransactionID, evt)
			case <-timeoutctx.Done():
				fmt.Println("事件超时")
				return
			}
		}
	}()
	return resp, nil
}


//一会儿用于调用的数据查询
func channelQuery(fcn string, args [][]byte) (channel.Response, error) {
	//上下文对象
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))
	//创建Client对象
	//源码在fabric-sdk-go\pkg\client\channel下
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	//状态查询
	return cli.Query(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.example.com"))
}

//定义接收的用户相关参数
//在客户端发送时必须存在，否则报错，gin框架
type UserRegisterRequest struct {
	Id   string `form:"id" binding:"required"`
	Name string `form:"name" binding:"required"`
}

//用户注册
func userRegister(ctx *gin.Context) {
	req := new(UserRegisterRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	resp, err := channelExecute("userRegister", [][]byte{
		[]byte(req.Name),
		[]byte(req.Id),
	})
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

//查询用户
func queryUser(ctx *gin.Context) {
	//取到id，根据id去查询
	userId := ctx.Param("id")
	//调用上面定义的查询的方法
	resp, err := channelQuery("queryUser", [][]byte{
		[]byte(userId),
	})
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

//用户注销
func deleteUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	resp, err := channelExecute("userDestory", [][]byte{
		[]byte(userId),
	})
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

//资产查询
func queryAsset(ctx *gin.Context) {
	assetId := ctx.Param("id")
	resp, err := channelQuery("queryAsset", [][]byte{
		[]byte(assetId),
	})
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

//定义接收的资产相关的参数
type AssetsEnrollRequest struct {
	AssetId   string `form:"assetsid" binding:"required"`
	AssetName string `form:"assetname" binding:"required"`
	Metadata  string `form:"metadata" binding:"required"`
	OwnerId   string `form:"ownerid" binding:"required"`
}

//资产登记
func assetsEnroll(ctx *gin.Context) {
	req := new(AssetsEnrollRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	//调用添加执行的方法
	resp, err := channelExecute("assetEnroll", [][]byte{
		[]byte(req.AssetName),
		[]byte(req.AssetId),
		[]byte(req.Metadata),
		[]byte(req.OwnerId),
	})

	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

//定义接收的资产变更的相关的参数
type AssetsExchangeRequest struct {
	AssetId        string `form:"assetsid" binding:"required"`
	OriginOwnerId  string `form:"originownerid" binding:"required"`
	CurrentOwnerId string `form:"currentownerid" binding:"required"`
}

//资产转让
func assetsExchange(ctx *gin.Context) {
	req := new(AssetsExchangeRequest)
	//参数绑定
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	resp, err := channelExecute("assetExchange", [][]byte{
		[]byte(req.OriginOwnerId),
		[]byte(req.AssetId),
		[]byte(req.CurrentOwnerId),
	})
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

//资产变更历史查询
func assetsExchangeHistory(ctx *gin.Context) {
	assetId := ctx.Query("assetid")
	queryType := ctx.Query("querytype")

	resp, err := channelQuery("queryAssetHistory", [][]byte{
		[]byte(assetId),
		[]byte(queryType),
	})
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
