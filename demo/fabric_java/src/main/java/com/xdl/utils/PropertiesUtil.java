package com.xdl.utils;

import java.io.InputStream;
import java.util.Properties;

//专门用于读取配置文件
public class PropertiesUtil {
	// 静态方法
	public static String read(String path, String key) {
		Properties pro = new Properties();
		String value = null;
		try {
			// 获取一个流对象
			InputStream in = PropertiesUtil.class.getClassLoader()
					.getResourceAsStream(path);
			// 将流对象加载
			pro.load(in);
			value = pro.getProperty(key);
		} catch (Exception e) {
			e.printStackTrace();
		}
		return value;
	}
}
