//定义严格模式
'use strict';

//注册管理员用户
//引入模块
var Fabric_Client = require('fabric-client');
var Fabric_CA_Client = require('fabric-ca-client');
var path = require('path');

//定义变量
var fabric_client = new Fabric_Client();
var fabric_ca_client = null;
var admin_user = null;

//定义存放证书和私钥的位置
var store_path = path.join(__dirname, 'hfc-key-store');
console.log('Store path:' + store_path);

//创建一个Client客户端
Fabric_Client.newDefaultKeyValueStore({
    path: store_path
}).then((state_store) => {
    //设置客户端存储
    fabric_client.setStateStore(state_store);
    //密码学相关套件的获取
    var crypto_suite = Fabric_Client.newCryptoSuite();
    //存到相同的位置
    var crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
    //认证
    crypto_suite.setCryptoKeyStore(crypto_store);
    fabric_client.setCryptoSuite(crypto_suite);

    //verify：验证
    //跳过tls验证
    var tlsOptions = {
        trustedRoots: [],
        verify: false
    };

    // cs：端口7054
    fabric_ca_client = new Fabric_CA_Client('http://192.168.20.224:7054', tlsOptions, 'ca.example.com', crypto_suite);

    //返回检查admin是否已经注册
    return fabric_client.getUserContext('admin', true);
}).then((user_from_store) => {
    //判断是否已经注册
    if (user_from_store && user_from_store.isEnrolled()) {
        console.log('成功加载admin');
        admin_user = user_from_store;
        return null;
    } else {
        //注册admin
        //注册ca服务器
        //拼装请求参数
        return fabric_ca_client.enroll({
            enrollmentID: 'admin',
            enrollmentSecret: 'adminpw'
        }).then((enrollment) => {
            console.log('成功注册admin')
            return fabric_client.createUser({
                username: 'admin',
                mspid: 'Org1MSP',
                cryptoContent: {
                    privateKeyPEM: enrollment.key.toBytes(),
                    signedCertPEM: enrollment.certificate
                }
            });
        }).then((user) => {
            admin_user = user;
            //设置user的上下文对象
            return fabric_client.setUserContext(admin_user);
        }).catch(() => {
            console.error('注册或保存admin失败');
            throw new Error('注册admin失败');
        });
    }
}).then(() => {
    console.log('注册admin用户到fabric client:' + admin_user.toString());
}).catch((err) => {
    console.log('注册admin失败:' + err)
});



















