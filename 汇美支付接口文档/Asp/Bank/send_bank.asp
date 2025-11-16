<!--#include file="asp_md5.asp"-->
<!--#include file="merchant.asp"-->
<%
dim parter
'商户ID
dim banktype
'银行类型
dim price
'支付金额
dim orderid
'商户订单号
dim callbackurl
'支付回调页面
dim hrefbackurl
'支付完成后的跳转页面
dim sign
'签名
userkey=u_userkey
'商户密钥，用户自行替换
parter=u_parter
'测试商户ID，用户自行替换
callbackurl=u_callbackurl
'支付回调页面
hrefbackurl=u_hrefbackurl
'支付跳转回页面
price=Trim(request("Price2"))
'支付金额
bankid=Trim(request("bankid"))
'银行卡ID
orderid = Trim(request("orderid"))
'商户系统订单号
if orderid = "" then
	Randomize 
	rnds = Int((900 * Rnd) + 100)
	orderid=year(now())&month(now())&day(now())&hour(now())&minute(now())&second(now())&rnds'生成商户订单号,商户可自行定义
end if
if price = "" or not(isnumeric(price)) then
	response.Write "请选择充值的金额"
	response.End()
end if
sign=asp_md5("parter="&parter&"&type="&bankid&"&orderid="&orderid&"&callbackurl="&callbackurl&userkey)
url = u_sendurl
sendurl = url & "?type="&bankid&"&parter="&parter&"&value="&price&"&orderid="&orderid&"&callbackurl="&callbackurl&"&hrefbackurl="&hrefbackurl&"&sign="&sign
'response.Write(sendurl)
response.Redirect sendurl
response.End()
%>