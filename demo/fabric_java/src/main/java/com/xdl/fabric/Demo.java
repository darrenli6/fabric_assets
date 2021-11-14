package com.xdl.fabric;

import java.util.Collection;
import java.util.Properties;

import org.hyperledger.fabric.sdk.ChaincodeID;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.HFClient;
import org.hyperledger.fabric.sdk.ProposalResponse;
import org.hyperledger.fabric.sdk.QueryByChaincodeRequest;
import org.hyperledger.fabric.sdk.TransactionProposalRequest;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric_ca.sdk.HFCAClient;

import com.xdl.utils.CertUtils;
import com.xdl.utils.PropertiesUtil;

public class Demo {
	public static void main(String[] args) throws Exception {
		// 用户注册，保存证书和私钥
		enroll("admin", "adminpw", "cert");
		// 账本更新
		// update("set", "b", "256");
		// 账本查询
		query("get", "b");
	}

	// 做用户注册，保存证书和私钥
	// certDir指定目录
	public static void enroll(String username, String password, String certDir)
			throws Exception {
		// HFClient对象
		// 创建客户端实例
		HFClient client = HFClient.createNewInstance();
		// 密码学相关套件
		CryptoSuite cs = CryptoSuite.Factory.getCryptoSuite();
		// 认证
		client.setCryptoSuite(cs);
		// 获取Properties对象
		Properties prop = new Properties();
		// 将verify（认证）写入prop
		prop.put("verify", false);
		// 得到认证后的客户端实例
		HFCAClient caClient = HFCAClient.createNewInstance(
				PropertiesUtil.read("demo.properties", "ca"), prop);
		// 认证
		caClient.setCryptoSuite(cs);
		// 得到Enrollment对象
		Enrollment enrollment = caClient.enroll(username, password);
		// 打印证书
		System.out.println(enrollment.getCert());
		// 将证书和私钥保存到本地硬盘
		CertUtils.saveEnrollment(enrollment, certDir, username);
	}

	// 初始化通道
	public static Channel initChannel(HFClient client) throws Exception {
		// 获取密码学相关套件
		CryptoSuite cs = CryptoSuite.Factory.getCryptoSuite();
		client.setCryptoSuite(cs);
		// 设置User上下文对象
		client.setUserContext(new MyUser("admin", CertUtils.loadEnrollment(
				"cert", "admin")));
		// 初始化channel
		Channel channel = client.newChannel("mychannel");
		// 加入peer
		channel.addPeer(client.newPeer("peer",
				PropertiesUtil.read("demo.properties", "peer")));
		// 指定排序节点的地址
		// 无论后面是否执行增删改查，这里必须指定排序节点
		channel.addOrderer(client.newOrderer("orderer",
				PropertiesUtil.read("demo.properties", "orderer")));
		channel.initialize();
		return channel;
	}

	// 更新账本
	// fun:调用的链码方法
	// 可变参：调用方法的参数是不确定的个数
	public static void update(String fun, String... args) throws Exception {
		// 拿到客户端对象
		HFClient client = HFClient.createNewInstance();
		// 进行初始化通道
		Channel channel = initChannel(client);
		// 构建一个提案
		TransactionProposalRequest req = client.newTransactionProposalRequest();
		// 指定要调用的链码，chaincode
		// 得到ChaincodeID对象
		ChaincodeID cid = ChaincodeID.newBuilder()
				.setName(PropertiesUtil.read("demo.properties", "chaincode"))
				.build();
		// 给提案携带参数
		req.setChaincodeID(cid);
		// 调用的方法名
		req.setFcn(fun);
		req.setArgs(args[0], args[1]);
		// 发送提案，从channel发送
		// 最终返回一个集合对象，泛型类型是提案响应对象
		Collection<ProposalResponse> resps = channel
				.sendTransactionProposal(req);
		// 将背书节点返回结果，提交到排序节点
		channel.sendTransaction(resps);
		System.out.println("更新完成");
	}

	// 查询账本
	public static void query(String fun, String... args) throws Exception {
		// 构建客户端对象
		HFClient client = HFClient.createNewInstance();
		// 初始化通道
		Channel channel = initChannel(client);
		// 构建提案，创建查询请求
		QueryByChaincodeRequest req = client.newQueryProposalRequest();
		// 指定调用的链码
		ChaincodeID cid = ChaincodeID.newBuilder()
				.setName(PropertiesUtil.read("demo.properties", "chaincode"))
				.build();
		// 给请求设置参数
		req.setChaincodeID(cid);
		req.setFcn(fun);
		req.setArgs(args[0]);
		// 账本查询
		Collection<ProposalResponse> resps = channel.queryByChaincode(req);
		// 遍历集合，取值
		for (ProposalResponse resp : resps) {
			// 最终返回一个Payload
			String payload = new String(
					resp.getChaincodeActionResponsePayload());
			System.out.println("response:" + payload);
		}
	}
}
