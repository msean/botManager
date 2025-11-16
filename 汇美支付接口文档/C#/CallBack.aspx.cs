using System;
using System.Configuration;
using System.Web.Security;

public partial class Receive : System.Web.UI.Page
{
    protected void Page_Load(object sender, EventArgs e)
    {
        string key = ConfigurationManager.AppSettings["userkey"];//配置文件密钥
        //返回参数
        string orderid = Request["orderid"];//返回订单号
        string restate = Request["restate"];//返回处理结果
        string ovalue = Request["ovalue"];//返回实际充值金额
        string sign = Request["sign"];//返回签名
        string attach = Request["attach"];//上行附加信息

        string param = string.Format("orderid={0}&restate={1}&ovalue={2}{3}", orderid, restate, ovalue, key);//组织参数
        //比对签名是否有效
        if (sign.Equals(FormsAuthentication.HashPasswordForStoringInConfigFile(param, "MD5").ToLower()))
        {
            Response.Write("ok");
            //执行操作方法
            if (restate.Equals("0"))
            {
                //支付成功 更新数据库操作
            }
            else 
            {
                //支付失败 更新数据库操作
            }

        }
        else
        { 
            //签名无效
        }
    }
}
