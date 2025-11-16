using System;
using System.Configuration;
using System.Web.Security;
using System.Net;
using System.IO;
using System.Text;

public partial class Send : System.Web.UI.Page
{
    string orderid =null;
    string shop_id =null;
    string key = null;
    protected void Page_Load(object sender, EventArgs e)
    {

        orderid = Guid.NewGuid().ToString().Substring(0, 20).Replace("-", "");
        shop_id = ConfigurationManager.AppSettings["myid"]; //商户ID
        key = ConfigurationManager.AppSettings["mykey"]; //商户密钥
        #region 原代码
        if (Request.Form["bankCardType"] == "00")
        {
            ChargeBank(orderid);
        }
        else if (Request.Form["bankCardType"] == "01")
        {
            CardReceive(orderid);
        }
        #endregion

    }

    private void ChargeBank(string orderid)
    {
        //商户信息
        string bankgate = ConfigurationManager.AppSettings["bankgate"]; //
        string callbackurl = "http://" + Request.Url.Host + ":" + Request.Url.Port + "/pay/CallBack.aspx";
        string hrefbackurl = "http://" + Request.Url.Host + ":" + Request.Url.Port + "/pay/HrefBack.aspx";
        string attach = "";
        //银行提交获取信息
        string bank_Type = Request.Form["bankCode"];//银行类型
        string bank_payMoney = Request.Form["totalAmount"];//充值金额
        string param = string.Format("parter={0}&type={1}&orderid={2}&callbackurl={3}", shop_id, bank_Type, orderid, callbackurl);
        string PostUrl = string.Format(bankgate + "?{0}&value={1}&sign={2}&hrefbackurl={3}&attach={4}", param, bank_payMoney, FormsAuthentication.HashPasswordForStoringInConfigFile(param + key, "MD5").ToLower(), hrefbackurl, attach);
        Response.Redirect(PostUrl,true);//转发URL地址
        //Response.Write(PostUrl);
    }

    private void CardReceive(string orderid)
    {
        //商户信息
        string cardgate = ConfigurationManager.AppSettings["cardgate"]; //
        string callbackurl = "http://" + Request.Url.Host + ":" + Request.Url.Port + "/pay/CallBack.aspx";
        string hrefbackurl = "http://" + Request.Url.Host + ":" + Request.Url.Port + "/pay/HrefBack.aspx";
        //获取卡类提交信息
        string card_No = Request.Form["cardNo"];//卡号
        string card_pwd = Request.Form["cardPwd"];//卡密
        //string card_account = Request.Form["txtUserNameCard"];//充值账号
        string card_type = Request.Form["bankCode"].Split('，')[0];//卡类型
        string card_payMoney = Request.Form["totalAmount"];//充值金额
        string restrict = "0";//使用范围
        string attach = "test";//附加内容，下行原样返回

        //卡类支付
        string param = "parter=" + shop_id + "&cardtype=" + card_type + "&cardno=" + card_No + "&cardpwd=" + card_pwd + "&orderid=" + orderid + "&callbackurl=" + callbackurl + "&restrict=" + restrict + "&price=" + card_payMoney;
        string PostUrl = string.Format(cardgate+"?{0}&attach={1}&sign={2}", param, attach, FormsAuthentication.HashPasswordForStoringInConfigFile(param + key, "MD5").ToLower());

        HttpWebRequest httpWebRequest = (HttpWebRequest)WebRequest.Create(PostUrl);
        //获取响应结果 此过程大概需要5秒
        HttpWebResponse httpWebResponse = (HttpWebResponse)httpWebRequest.GetResponse();
        //获取响应流
        Stream stream = httpWebResponse.GetResponseStream();
        //用指定的字符编码为指定的流初始化 StreamReader 类的一个新实例。
        using (StreamReader streamReader = new StreamReader(stream, Encoding.UTF8))
        {
            string useResult = streamReader.ReadToEnd();
            streamReader.Dispose();
            streamReader.Close();
            httpWebResponse.Close();

            if (useResult.Trim() == "0")
            {
                Response.Write("支付成功.");
            }
            else
            {
                Response.Write(useResult.Trim());
            }
        }
    }
}
