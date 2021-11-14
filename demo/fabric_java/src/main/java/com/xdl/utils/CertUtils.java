package com.xdl.utils;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.spec.PKCS8EncodedKeySpec;

import javax.xml.bind.DatatypeConverter;

import org.hyperledger.fabric.sdk.Enrollment;

//专门用于证书和密钥相关的工具类
public class CertUtils {
	// 私有构造方法
	private CertUtils() {
	}

	// 将密钥和证书保存到cert目录下
	public static void saveEnrollment(Enrollment enrollment, String dir,
			String name) {
		if (null == enrollment) {
			try {
				throw new Exception("enrollment不能为空");
			} catch (Exception e) {
				e.printStackTrace();
			}
		}
		// 保存的cert证书的目录+文件名
		// File.separator：根据不同操作系统，去自动拼接 /或者\\
		String certFileName = String.join("", dir, File.separator, name,
				".cert");
		FileOutputStream certOut = null;
		try {
			certOut = new FileOutputStream(certFileName);
			certOut.write(enrollment.getCert().getBytes());
		} catch (Exception e) {
			e.printStackTrace();
		} finally {
			try {
				certOut.close();
			} catch (IOException e) {
				e.printStackTrace();
			}
		}

		// 保存私钥
		String keyFileName = String
				.join("", dir, File.separator, name, ".priv");
		// 文件输出流对象
		FileOutputStream keyOut = null;
		try {
			keyOut = new FileOutputStream(keyFileName);
			StringBuilder sb = new StringBuilder(300);
			sb.append("---BEGIN PRIVATE KEY---\n");
			String priKey = DatatypeConverter.printBase64Binary(enrollment
					.getKey().getEncoded());
			int len = priKey.length();
			for (int i = 0; i < len; ++i) {
				sb.append(priKey.charAt(i));
				// 每行64个字符输出
				if ((i + 1) % 64 == 0) {
					sb.append('\n');
				}
			}
			sb.append("\n---END PRIVATE KEY---\n");
			keyOut.write(sb.toString().getBytes());
		} catch (Exception e) {
			e.printStackTrace();
		} finally {
			try {
				keyOut.close();
			} catch (IOException e) {
				e.printStackTrace();
			}
		}
	}

	// 从文件中读取身份信息
	public static Enrollment loadEnrollment(String dir, String name)
			throws Exception {
		// 读取证书文件
		byte[] certBuf = Files.readAllBytes(Paths.get(dir, name + ".cert"));
		String cert = new String(certBuf);

		// 读取私钥对象，构造PrivateKey对象
		PrivateKey key = loadPrivateKey(Paths.get(dir, name + ".priv"));
		// 直接构造身份信息的封装
		return new MyEnrollment(key, cert);
	}

	private static PrivateKey loadPrivateKey(Path fileName) {
		PrivateKey key = null;
		FileInputStream is = null;
		BufferedReader br = null;
		try {
			is = new FileInputStream(fileName.toFile());
			// 缓冲流
			br = new BufferedReader(new InputStreamReader(is));
			StringBuilder builder = new StringBuilder();

			// 标记
			boolean inKey = false;
			for (String line = br.readLine(); line != null; line = br
					.readLine()) {
				if (!inKey) {
					if (line.startsWith("---BEGIN")
							&& line.endsWith("PRIVATE KEY---")) {
						inKey = true;
					}
					continue;
				} else {
					if (line.startsWith("---END")
							&& line.endsWith("PRIVATE KEY---")) {
						inKey = false;
						break;
					}
					builder.append(line);
				}
			}
			// 将StringBuilder转换为PrivateKey对象进行返回
			byte[] encode = DatatypeConverter.parseBase64Binary(builder
					.toString());
			PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(encode);
			KeyFactory kf = KeyFactory.getInstance("EC");
			key = kf.generatePrivate(keySpec);
		} catch (Exception e) {
			e.printStackTrace();
		} finally {
			try {
				br.close();
			} catch (IOException e) {
				e.printStackTrace();
			}
			try {
				is.close();
			} catch (IOException e) {
				e.printStackTrace();
			}
		}
		return key;
	}
}

// 需要一个实现了Enrollment接口的实现类
class MyEnrollment implements Enrollment {

	private PrivateKey privateKey;
	private String cert;

	public MyEnrollment(PrivateKey privateKey, String cert) {
		this.privateKey = privateKey;
		this.cert = cert;
	}

	@Override
	public PrivateKey getKey() {
		return this.privateKey;
	}

	@Override
	public String getCert() {
		return this.cert;
	}
}
