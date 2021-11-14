'use strict'

//查询汽车信息
//引入模块
var Fabric_Client = require('fabric-client');
var path = require('path');

//定义变量
var fabric_client = new Fabric_Client();

//连接fabric网络
//初始化通道
var channel = fabric_client.newChannel('mychannel');
//peer节点设置
var peer = fabric_client.newPeer('grpc://192.168.20.224:7051');
//加入通道
channel.addPeer(peer);
var member_user = null;
var store_path = path.join(__dirname, 'hfc-key-store');
console.log('store path:' + store_path);

//创建客户端进行操作
Fabric_Client.newDefaultKeyValueStore({
    path: store_path
}).then((state_store) => {
    //设置客户端存储
    fabric_client.setStateStore(state_store);
    //密码学相关套件获取
    var crypto_suite = Fabric_Client.newCryptoSuite();
    //存储路径
    var crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
    //认证相关
    crypto_suite.setCryptoKeyStore(crypto_store);
    fabric_client.setCryptoSuite(crypto_suite);
    //返回检查user1是否已经注册
    return fabric_client.getUserContext('user1', true);
}).then((user_from_store) => {
    if (user_from_store && user_from_store.isEnrolled()) {
        console.log('成功加载用户');
        member_user = user_from_store;
    } else {
        throw new Error('获取user1失败，请先运行registerUser.js');
    }
    //user1已经拿到了
    //调用链码查询
    //调用链码实际上本质还是一个request
    //封装请求参数
    const request = {
        //链码名字
        chaincodeId: 'node',
        //方法
        //fcn: 'queryAllCars',
        fcn:'queryCar',
        //请求参数
        //args: ['']
        args:['CAR1']
    };
    //发送提案
    return channel.queryByChaincode(request);
}).then((query_responses) => {
    console.log('查询完成');
    //检验查询结果
    if (query_responses && query_responses.length == 1) {
        if (query_responses[0] instanceof Error) {
            console.error("查询错误：", query_responses[0]);
        } else {
            console.log("返回结果是：", query_responses[0].toString());
        }
    } else {
        console.log('查询结果无效');
    }
}).catch(() => {
});

















