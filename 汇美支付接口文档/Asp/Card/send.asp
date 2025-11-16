<!--#include file="asp_md5.asp"-->
<!--#include file="merchant.asp"-->
<%
dim cardType
'卡类型
dim parter
'商户ID
dim cardNo
'卡号
dim cardPwd
'卡密
dim values
'提交金额
dim restrict
'使用范围
dim orderid
'商户定单号
dim callbackurl
'与满支付平台下行通知商户地址
dim userkey
'智多赢科技支付分配给商户密钥
parter = u_parter
'商户ID
userkey = u_userkey
'商户密钥
cardtype = Trim(request("cardtype"))
'点卡类型
cardno = Trim(request("cardno"))
'卡号
cardpwd = Trim(request("cardpwd"))
'卡密
price = Trim(request("price2"))
'充值金额
restrict = 0
'使用范围，仅限制于神州行充值卡
orderid = Trim(request("num"))
if orderid = "" then
	Randomize 
	rnds = Int((900 * Rnd) + 100)
	orderid=year(now())&month(now())&day(now())&hour(now())&minute(now())&second(now())&rnds''''生成商户订单号,商户可自行定义
end if
'商户系统订单号
userkey = u_userkey
'商户密钥，用户自行替换
parter = u_parter
'测试商户ID，用户自行替换
callbackurl = u_callbackurl
'支付结果回调地址
if restrict = "" or not(isnumeric(restrict)) then restrict = 0
'判断参数是否有效
if  parter="" or cardno="" or cardpwd="" or price="" or  restrict="" or  callbackurl="" or userkey="" or orderid = "" then
	response.Write "参数不正确"
	response.End()
end if
sendurl = u_sendurl
'提交到智多赢科技支付点卡接口的地址
dim sign'参数进行加密签名
signtxt = "parter="&parter&"&cardtype="&cardtype&"&cardno="&cardno&"&cardpwd="&cardpwd&"&orderid="&orderid&"&callbackurl="&callbackurl&"&restrict="&restrict&"&price="&price&userkey
sign=asp_md5(signtxt)
send_xyurl = sendurl&"?parter="&parter&"&cardtype="&cardtype&"&cardno="&cardNo&"&cardpwd="&cardPwd&"&orderid="&orderid&"&callbackurl="&callbackurl&"&restrict="&restrict&"&price="&price&"&sign="&sign
response.write(send_xyurl)
response.Redirect(send_xyurl)
response.End()
'++++++++++++++++++++++++++++
'商户写入自己平台的逻辑
'++++++++++++++++++++++++++++
resultStr=xmlGet(send_xyurl)'异步提交，并返回消息
'resultStr=cstr(right(resultStr,1))
select case trim(resultStr)
	case "0" 
		'提交成功以后。跳转到显示支付状态的页面
		response.Write "ok"
	case "1"
		'参数无效直接通知商户，不在进入下行逻辑
		response.Write "参数无效,请检查之后重新提交"
	case "2"
		'签名验证无效，不在进入下行逻辑
		response.Write "签名验证无效"
	case "3"
		'卡密为重复提交
		response.Write "卡密为重复提交"
	case "4"
		'卡面值不符合规则
		response.Write "卡面值不符合规则"
	case else
		'未知的网络原因
		response.Write resultStr
end select
response.end()
Function xmlGet(url)
	set objXML = server.CreateObject("MSXML2.ServerXMLHTTP.3.0")
	objXML.open "GET", url, false
	objXML.setRequestHeader "Accept","image/gif, image/x-xbitmap, image/jpeg, image/pjpeg, application/vnd.ms-excel, application/vnd.ms-powerpoint, application/msword, application/x-shockwave-flash, */*"
	objXML.setRequestHeader "Accept-Language","zh-cn"
	objXML.setRequestHeader "Content-Type", "application/x-www-form-urlencoded"
	objXML.setRequestHeader "Accept-Encoding", "gzip, deflate"
	objXML.setRequestHeader "User-Agent","Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)"
	objXML.setRequestHeader "Connection","Keep-Alive" 
	objXML.setRequestHeader "Cache-Control","no-cache"
	objXML.send()
	if objXML.readystate<>4 then
		exit function
	end if
	set oStream = server.CreateObject("ADODB.Stream")
	oStream.Type=1
	oStream.Mode=3
	oStream.Open()
	oStream.Write objXML.responseBody
	oStream.Position= 0
	oStream.Type= 2
	oStream.Charset="gb2312"
	ReturnText = oStream.readtext()
	oStream.Close()
	set oStream = nothing
	set objXML = nothing
	xmlGet = ReturnText	
End function
%>