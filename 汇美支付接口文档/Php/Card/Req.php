<?php
include_once("config.php");

$keyvalue = $key;//用户key
$parter = $id;//用户ID

$orderid =  date("Ymdhis");//用户订单号（必须唯一）
$price =  $_POST["Price2"];//订单金额
$cardtype =  $_POST["cardtype"];//卡类型ID（见文档）
$cardno =  $_POST["cardno"];//卡号
$cardpwd =  $_POST["cardpwd"];//卡密
$callbackurl =  $callback_url;//同步接收返回URL连接
$attach = "123";//备注 如有中文 urlencode编码

$sign = "parter=".$parter."&cardtype=".$cardtype."&cardno=".$cardno."&cardpwd=".$cardpwd."&orderid=".$orderid."&callbackurl=".$callbackurl."&restrict=0"."&price=".$price;
$signs = md5($sign.$key);//签名数据 32位小写的组合加密验证串
$url=$submiturl."?".$sign."&attach=".$attach."&sign=".$signs;
$result =file_get_contents($url); //提交

if($result=="0")
{
	echo "订单提交成功!";
}
else
{
	echo $result;
}
?>