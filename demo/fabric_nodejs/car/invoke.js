'use strict';

//更新汽车信息
//引入模块
var Fabric_Client = require('fabric-client');
var path = require('path');

//定义变量
var fabric_client = new Fabric_Client();
var member_user = null;
var store_path = path.join(__dirname, 'hfc-key-store');
console.log('store path:' + store_path);
//声明一个事务id
var tx_id = null;

//连接网络
var channel = fabric_client.newChannel('mychannel');
//peer节点
var peer = fabric_client.newPeer('grpc://192.168.20.224:7051');
//加入通道，通道连接peer节点
channel.addPeer(peer);
//orderer节点的设置
var order = fabric_client.newOrderer('grpc://192.168.20.224:7050');
//通道连接orderer节点
channel.addOrderer(order);

//创建一个Client
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
        console.log('成功加载到user1');
        member_user = user_from_store;
    } else {
        throw new Error('获取user1失败，请先运行registerUser.js');
    }
    //已经拿到了user1对象
    tx_id = fabric_client.newTransactionID();
    console.log('分配的tx_id:', tx_id._transaction_id);
    //封装请求参数
    var request = {
        chaincodeId: 'node',
        fcn: 'changeCarOwner',
        args: ['CAR1', 'sjc'],
        chainId: 'mychannel',
        txId: tx_id
    }

    //提交事务
    //本质就是一个request
    return channel.sendTransactionProposal(request);
}).then((results) => {
    // results[0]就是提案结果
    var proposalResponse = results[0];
    // results[1]就是提案内容
    var proposal = results[1];
    let isProposalGood = false;
    //检验返回结果
    if (proposalResponse && proposalResponse[0].response && proposalResponse[0].response.status === 200) {
        //提案成功
        isProposalGood = true;
        console.log('提案事务执行成功');
    } else {
        console.error('提案事务执行失败');
    }
    //当提案成功时
    if (isProposalGood) {
        console.log('提案成功');

        //封装请求参数
        //背书节点执行的结果和内容的封装
        var request = {
            proposalResponses: proposalResponse,
            proposal: proposal
        };

        //创建事务提案并发送给orderer
        //获取事务id
        var transaction_id_string = tx_id.getTransactionID();
        //通道发送事务，本质也是一个请求
        var sendPromise = channel.sendTransaction(request);
        //定义一个数组，用于发送事务数据
        var promises = [];
        //加入数组
        promises.push(sendPromise);

        //创建EventHub对象
        let event_hub = channel.newChannelEventHub('192.168.20.224:7051');
        //resolve：成功后执行的
        //reject：失败后执行的
        let txPromise = new Promise((resolve, reject) => {
            let handle = setTimeout(() => {
                event_hub.disconnect();
                resolve({event_status: 'TIMEOUT'});
            }, 3000);
            event_hub.connect();
            event_hub.registerTxEvent(
                transaction_id_string, (tx, code) => {
                    clearTimeout(handle);
                    event_hub.unregisterTxEvent(transaction_id_string);
                    event_hub.disconnect();
                    var return_status = {
                        event_status: code, tx_id: transaction_id_string
                    };
                    if (code !== 'VALID') {
                        console.error('交易无效,code=' + code);
                        //返回执行状态
                        resolve(return_status);
                    } else {
                        console.log('事务已经提交');
                        resolve(return_status);
                    }
                }, (err) => {
                    //事务回调失败
                    reject(new Error('执行失败,' + err));
                }
            );
        });
        promises.push(txPromise);
        //all()：确保promises被查出来，起到同步的作用
        return Promise.all(promises);
    } else {
        console.error('发送或接收失败');
        throw new Error('发送或接收失败');
    }
}).then((results) => {
    console.log('发送交易和事件监听完成');
    //数据校验
    if (results && results[0] && results[0].status == 'SUCCESS') {
        console.log('成功发送事务到orderer');
    } else {
        console.log('发送事务到orderer失败');
    }
}).catch(() => {
});










