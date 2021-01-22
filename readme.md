# Project Development

## Task Assignment
* clf：后端开发，golang＋beego
* sy、mzw：sgx开发（app+Enclave）
* lzh：后端+sgx+衔接部分

## 1.15-1.22
* clf、lzh：golang与beego框架的学习，22号即可开始开发
* sy、mzw：熟悉sgx开发的基本流程，包括tantivy的使用（索引构建、查询等等）
* sy、mzw、lzh：rust的学习
* beanchmark代码：端口7585，其它默认

## 1.22-1.28
* clf、lzh：
  * 完成用户注册、登录、文件上传、文件删除，可参照https://tank.eyeblue.cn/matter/list
  * 前端可暂缓
  * 了解json
  * mysql：root，sgx12345

| Models  |                                                              |
| ------- | ------------------------------------------------------------ |
| User    | 用户名、密码、头像、邮箱（后期可邮箱验证码注册、邮箱登录）、待扩展... |
| Article | User（一对多）、标题、待扩展...                              |

* sy、mzw、lzh：
  * 熟悉sgx开发的基本流程，读懂sample代码
  * 学长给的库，sgx部分用rust实现，包括tantivy的使用，https://gitee.com/ggdG/datacenter-apiserver/tree/master/tantivy-sgx-part/sgx
  * rust的学习，尝试写出rust的helloworld

|      | 服务器                             |
| ---- | ---------------------------------- |
| lzh  | ssh -p 7585  emison@... |
| sy   | ssh -p 7585  take@...   |
| clf  | ssh -p 7585  clf@...    |
| mzw  | ssh -p 7585  hamer@...  |

* 主要内容：添加、查询（按标题搜索、按关键字搜索）、删除
