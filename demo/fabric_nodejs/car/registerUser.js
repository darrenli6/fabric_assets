'use strict';

//引入模块
var Fabric_Client = require('fabric-client');
var Fabric_CA_Client = require('fabric-ca-client');
var path = require('path');

var fabric_client = new Fabric_Client();
var fabric_ca_client = null;
var admin_user = null;
var member_user = null;
var store_path = path.join(__dirname, 'hfc-key-store');
console.log('Store path:' + store_path)

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
    //设置ca访问
    fabric_ca_client = new Fabric_CA_Client('http://192.168.20.224:7054', null, '', crypto_suite);
    //返回检查admin是否已经注册
    return fabric_client.getUserContext('admin', true);
}).then((user_from_store) => {
    //判断是否已经注册
    if (user_from_store && user_from_store.isEnrolled()) {
        console.log('成功加载admin');
        admin_user = user_from_store;
    } else {
        throw new Error('获取admin失败，请先运行enrollAdmin.js');
    }
    //向ca注册一个用户user1
    //enrollmentID：就是要注册的用户名
    //第二个参数传的是管理员账户
    return fabric_ca_client.register({
        enrollmentID: 'user1',
        affiliation: 'org1.department1'
    }, admin_user);
}).then((secret) => {
    console.log('成功注册user1：' + secret);
    return fabric_ca_client.enroll({
        enrollmentID: 'user1',
        enrollmentSecret: secret
    });
}).then((enrollment) => {
    console.log('成功持久化用户user1');
    return fabric_client.createUser({
        username: 'user1',
        mspid: 'Org1MSP',
        cryptoContent: {
            privateKeyPEM: enrollment.key.toBytes(),
            signedCertPEM: enrollment.certificate
        }
    });
}).then((user) => {
    member_user = user;
    return fabric_client.setUserContext(member_user);
}).then(() => {
    console.log('成功注册用户user1，并添加到fabric上下文对象')
}).catch(() => {
});


















