# Netdisk

## 登录
http:url:port/login
### Showlogin
重写get方法，显示login.html页面
### HandleLogin
重写post方法，处理用户提交的表单

通过c.Data["string"]向前端传递字符串，前端页面通过{{.string}}接收

登录前检查该帐户是否已删除

登录成功后重定向到主页

### plan:
登录页面添加验证码

## 注册
http:url:port/register
### ShowRegister
重写get方法，显示register.html页面
### HandleRegister
重写post方法，处理用户提交表单

用户登录提供用户名、密码、重复密码、邮箱

注册成功后重定向到登录页面

### plan
注册页面添加验证码

## 注销功能
通过session实现，暂时应该无法区分多用户，后面解决

## 文件上传、搜索、删除功能
待添加
