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
'商户密钥，用户自行替换
if orderid="" or restate="" or ovalue="" or sign="" then
	response.Write("参数错误")
	response.End()
end if
mysign=asp_md5("orderid="&orderid&"&restate="&restate&"&ovalue="&ovalue&u_userkey)
if mysign <> Lcase(sign) then
'验证签名是否正确
	response.Write "签名不正确"
	response.Write mysign &"<br />"&sign
	response.End()
end if
dim msg
msg = ""
select case trim(Restate)
	case "0"'表示支付成功
		'+++++++++++++++++++++++++++
		'商户在此可以处理自己的逻辑
		'+++++++++++++++++++++++++++
		msg = "支付成功!<br>成功金额："&ovalue&"<br>订单号:"&orderid
	case else
		msg = "支付成功!<br>订单号:"&orderid
end select
response.Write msg
%>