package ledger

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/utils"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

type CalcHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	msg         botapi.Update
}

// ================= 匹配 =================

func (c *CalcHandler) Match(botModel bot.Bot, update botapi.Update) (match bool) {

	msg := update.Message
	input := strings.TrimSpace(msg.Text)

	c.botModel = botModel
	c.chatGroupID = msg.Chat.ID
	c.msg = update

	// 👉 预处理（兼容 TG 输入）
	input = normalizeExpr(input)

	// 👉 判断是不是表达式（核心）
	reg := regexp.MustCompile(`^[0-9\.\+\-\*/\(\) ]+$`)

	if reg.MatchString(input) {
		global.GVA_LOG.Debug("Calc match", zap.String("expr", input))
		return true
	}

	return false
}

// ================= 处理 =================

func (c *CalcHandler) Handle() (err error) {

	input := strings.TrimSpace(c.msg.Message.Text)

	// 👉 预处理
	input = normalizeExpr(input)

	result, err := calc(input)
	if err != nil {
		global.GVA_LOG.Error("calc error", zap.Error(err))
		return
	}

	botSender, err := botapi.NewBotAPI(c.botModel.Token)
	if err != nil {
		return
	}

	reply := fmt.Sprintf("%v", utils.FloatReserve(result, 6))

	msg := botapi.NewMessage(c.chatGroupID, reply)
	msg.ReplyToMessageID = c.msg.Message.MessageID

	_, err = botSender.Send(msg)
	return
}

// ================= 表达式计算 =================

func calc(expr string) (float64, error) {

	// 安全限制（防止乱输入）
	if len(expr) > 50 {
		return 0, fmt.Errorf("expression too long")
	}

	e, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, err
	}

	result, err := e.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	return result.(float64), nil
}

// ================= 预处理 =================

func normalizeExpr(expr string) string {

	expr = strings.TrimSpace(expr)

	// 中文符号替换
	expr = strings.ReplaceAll(expr, "×", "*")
	expr = strings.ReplaceAll(expr, "✖", "*")
	expr = strings.ReplaceAll(expr, "÷", "/")
	expr = strings.ReplaceAll(expr, "（", "(")
	expr = strings.ReplaceAll(expr, "）", ")")

	return expr
}
