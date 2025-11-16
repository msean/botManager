<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=gb2312" />
<title>- 在线银行付款</title>
</head>
<style>
*{ font-family:"微软雅黑";font-size:12px}
.STYLE2{font-weight: bold}
</style>
<body>
<html>
<div style="text-align:center; line-height:22px; margin-top:50px"><strong>请选择网上银行</strong></div>
<form name="form1" action="send_bank.asp" method="post">
  <table width="396" border="0" align="center" cellpadding="0" cellspacing="0" bgcolor="#FFFFFF" style="border:#0099FF solid 5px">
    <tr>
      <td height="431" colspan="2" align="center" bordercolor="#00CCFF"><table width="90%" border="0" cellpadding="0" cellspacing="1" bgcolor="#FFFFFF" style="border:#666666 dashed 1px; padding:5px">
        <tr>
          <td width="32%" height="25" align="left" bgcolor="#FFFFFF"><input name="bankid" type="radio" value="967" checked="checked" />
            <span class="STYLE2">中国工商银行</span></td>
          <td width="32%" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="964"  />
            中国农业银行</td>
        </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="965" />
            中国建设银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="971" />
邮政储蓄</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="969" />
            浙江稠州商业银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="970" />
            招商银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="972" />
            兴业银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="973"  />
            顺德农村信用合作社</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="975" />
            上海银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="976" />
            上海农村商业银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="978" />
            平安银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="980" />
            民生银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="982" />
            华夏银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="983" />
            杭州银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="985" />
            广东发展银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="986" />
            光大银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="988" />
            渤海银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="989" />
            北京银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="979" />
            南京银行</td>
          <td align="left" bgcolor="#FFFFFF"><input name="bankid" type="radio"  value="962" />
中信银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="963"  />
中国银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="974" />
深圳发展银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="968" />
浙商银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="977"  />
浦东发展银行</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="981"  />
交通银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="984"  />
广州市农村信用社</td>
          </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="987"  />
            东亚银行</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="991"  />
            北京农村商业银行</td>
        </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="TENPAY"  />
          财付通</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="KUAIJIE"  />
            快捷支付</td>
        </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="WEIXIN"  />
          微信支付</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="ALIPAY"  />
          支付宝余额支付</td>
        </tr>
        <tr>
          <td height="25" align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="WXWAP"  />
            手机微信</td>
          <td align="left" bgcolor="#FFFFFF"><input type="radio" name="bankid" value="ALIWAP"  />
            手机支付宝</td>
        </tr>
      </table></td>
    </tr>
    <tr>
      <td width="109" height="34" align="center"><strong>填写金额</strong></td>
      <td width="277">
      <input name="Price2" type="text" size="10" maxlength="10" />
        元</td>
    </tr>
    <tr>
      <td height="57" colspan="2" align="center"><input type="submit" name="submit2" value="确认付款" onClick="return checkMoney()" /> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<input type="button" name="submit0" value="返回上一步" onClick="history.go(-1)" /> </td>
    </tr>
  </table>
</form>

</html>
</body>
<script language="javascript">
function checkMoney(){
	if(document.form1.Price2.value ==""){
		alert('请输入充值的金额');
		return false;
	}
	document.form1.submit2.disabled;
}
</script>