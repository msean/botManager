<?php
include_once("config.php");


        $orderid  = trim($_GET['orderid']);
		$restate  = trim($_GET['restate']);
		$ovalue         = trim($_GET['ovalue']);
		$attach         = trim($_GET['attach']);
		$sign           = trim($_GET['sign']);
		
		$sign_text	= "orderid=".$orderid."&restate=".$restate."&ovalue=".$ovalue;
		$sign_md5 = md5($sign_text.$key);

        if (sign == localsign)
        { 
           if (opstate == "0")
           {
			    //成功逻辑处理，现阶段只发送成功的单据
            } 
        }
      else
	  {
		   //签名错误处理逻辑
	  }
       
?>