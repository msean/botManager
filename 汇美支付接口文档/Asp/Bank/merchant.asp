<%
Dim u_parter
'商户ID
Dim u_userkey
'通信密钥
Dim u_callbackurl
'支付回调页面
Dim u_hrefbackurl
'支付完成后跳转的页面
Dim u_sendurl
'提交地址请替换
u_parter = "10031"
'商户ID，支付时请替换为正式的商户ID
u_userkey = "082F4C92E5F403EDC4F0D5F71ACAE0CA"
'密钥，支付时请替换正式的密钥
u_callbackurl = "http://"&request.ServerVariables("HTTP_HOST")&"/NetApi/Test/Bank/callBack.asp"
' 支付回调页面
u_hrefbackurl = "http://"&request.ServerVariables("HTTP_HOST")&"/NetApi/Test/Bank/show_order.asp"
'支付完成后跳转的页面
u_sendurl = "http://pay.an216.com/interface/chargebank.aspx"
'提交地址请替换
%>