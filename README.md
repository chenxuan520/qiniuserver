# 七牛云文件快速上传

- **chenxuan**

## 作用

- 方便的七牛云客户端实现上传文件,管理文件方便

- 支持粘贴,拖入,选择文件上传

## 效果

![demo](http://cdn.androidftp.top/test/202404429163058pasteboard.paste)

## Quick Start

### 直接通过release版本下载

1. release 下载对应[操作系统版本](https://github.com/chenxuan520/qiniuserver/releases/)

2. 下载对应版本,解压,修改config文件夹内的demo.json(主要是添加七牛云参数)并重新命名为config.json

3. 运行qiniuserver

4. 打开浏览器

### 下载源码编译

1. 运行 `git clone https://github.com/chenxuan520/qiniuserver`

2. 运行 `cd src;go build .;mv ./qiniuserver ..;cd ..`

3. 修改config文件夹内的demo.json(主要是添加七牛云参数)并重新命名为config.json

4. 运行 qiniuserver

5. 打开浏览器

## 参数获取

- [access_key 和 secret_key 获取](https://portal.qiniu.com/developer/user/key)

- [bucket 的获取](https://portal.qiniu.com/kodo/bucket)

- [zone 的获取,取值为Huanan,Huabei,Huadong,Xingjiapo](https://portal.qiniu.com/kodo/bucket)

- upload_path 和 file_name 是自定义的

- [Api 文档](https://developer.qiniu.com/kodo/sdk/go-v6)

## 参数说明
```
{
	"host":{
		"ip":"127.0.0.1", // 绑定ip
		"port":"5200", // 绑定端口
		"password": "123" // 设置登录密码,如果为空意味着无需密码
	},
	"qiniu":{
		"domain":"demo", // cdn的域名
		"access_key":"demo", // qiniu云的accesskey
		"secret_key":"demo", // qiniu云的secretkey
		"bucket":"demo", // qiniu云的bucket
		"upload_path":"demo", // qiniu云的上传路径
		"file_name":"", // 自定义上传文件名,为空就是保持原文件名,%d表示日期,%f表示原文件名,%r表示6位随机字符串,目前只支持这几个
		"zone":"Huadong" // qiniu云的地区
	}
}
```
