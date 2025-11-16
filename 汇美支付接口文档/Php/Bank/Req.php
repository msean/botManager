<?php
include_once("config.php");

$keyvalue = $key;//用户key
$parter = $id;//用户ID

$orderid =  date("Ymdhis");//用户订单号（必须唯一）
$value =  $_POST["Price2"];//订单金额
$type =  $_POST["bankid"];//银行ID（见文档）
$callbackurl =  $callback_url;//异步接收返回URL连接
$hrefbackurl =  $$hrefback_url;//同步接收返回URL连接
$attach = "123";//备注 如有中文 urlencode编码
$sign = "parter=".$parter."&type=".$type."&orderid=".$orderid."&callbackurl=".$callbackurl;
$sign = md5($sign.$key);//签名数据 32位小写的组合加密验证串
$url=$submiturl."?parter=".$parter."&type=".$type."&value=".$value."&orderid=".$orderid."&callbackurl=".$callbackurl."&attach=".$attach."&hrefbackurl=".$hrefbackurl."&sign=".$sign;
header('Location:'.$url);

?>