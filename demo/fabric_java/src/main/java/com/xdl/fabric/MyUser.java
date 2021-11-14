package com.xdl.fabric;

import java.util.Collections;
import java.util.Set;

import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.User;

import com.xdl.utils.PropertiesUtil;

public class MyUser implements User {
	private String name;
	// 权限验证等信息
	private Enrollment enrollment;

	public MyUser(String name, Enrollment enrollment) {
		this.name = name;
		this.enrollment = enrollment;
	}

	public String getName() {
		return this.name;
	}

	public Set<String> getRoles() {
		return Collections.emptySet();
	}

	public String getAccount() {
		return "";
	}

	public String getAffiliation() {
		return "";
	}

	public Enrollment getEnrollment() {
		return this.enrollment;
	}

	public String getMspId() {
		return PropertiesUtil.read("demo.properties", "mspId");
	}
}
