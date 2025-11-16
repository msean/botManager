<?php
//=======================网银支付公用配置==================
//商户ID
$id		= '10000';
//通信密钥
$key		= '33B4DB3244BB5C2875D98CB1B79B68DA';	//hc6NOTDETVQe9Lgr
$submiturl = 'http://pay.an216.com/interface/CardReceive.aspx';  //点卡网关地址
//接收下行数据的地址, 该地址必须是可以再互联网上访问的网址
$callback_url  = "http://127.0.0.1/Pay/callback.php";   
?>