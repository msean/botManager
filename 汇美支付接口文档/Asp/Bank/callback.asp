<!--#include file="asp_md5.asp"-->
<!--#include file="merchant.asp"-->
<%
dim orderid
'支付订单号
dim Restate
'返回值
dim ovalue
'提交实际金额
dim sign
'签名
orderid=Trim(request("orderid"))
restate=Trim(request("restate"))
ovalue=Trim(request("ovalue"))
sign=Trim(request("sign"))
userkey = u_userkey
'商户密钥，用户自行替换
if orderid="" or Restate="" or ovalue="" or sign="" then
	response.Write("参数错误")
	response.End()
end if
mysign=asp_md5("orderid="&orderid&"&restate="&restate&"&ovalue="&ovalue&userkey)

if mysign <> sign then
'验证签名是否正确
	response.Write "签名不正确"
	response.End()
end if
select case trim(Restate)
	case "0"'表示支付成功
		'+++++++++++++++++++++++++++
		'商户在此可以处理自己的逻辑
		'+++++++++++++++++++++++++++
		esultstr="ok"
		resultstat=0
	case "-1"
		esultstr="ok"
		resultstat=-1
	case "-2"
		esultstr="ok"
		resultstat=-2
	case "-3"
		esultstr="ok"
		resultstat=-3
	case "-4"
		esultstr="ok"
		resultstat=-4
	case "-5"
		esultstr="ok"
		resultstat=-5
	case else
		esultstr="err"
		resultstat=-6
end select
response.Write esultstr
'商户系统收到通知以后请输出 小写"ok"表示接受通知成功
response.End()

%>