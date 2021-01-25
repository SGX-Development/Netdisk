# Netdisk

## 登陆
http:url:port/login
### Showlogin
重写get方法，显示login.html页面
### HandleLogin
重写post方法，处理用户提交的表单

通过c.Data["string"]向前端传递字符串，前端页面通过{{.string}}接收

登陆成功后重定向到主页

### plan:
登陆页面添加验证码

通过session记录登陆状态

前端没写

## 注册
http:url:port/register
### ShowRegister
重写get方法，显示register.html页面
### HandleRegister
重写post方法，处理用户提交表单

用户登陆提供用户名、密码、重复密码、邮箱

注册成功后重定向到登陆页面

### plan
注册页面添加验证码

前端没写

## 注销功能
待添加

## 文件上传、搜索、删除功能
待添加
