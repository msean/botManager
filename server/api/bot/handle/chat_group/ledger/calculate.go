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

func (c *CalcHandler) Match(botModel bot.Bot, update botapi.Update) (match bool) {

	msg := update.Message
	input := strings.TrimSpace(msg.Text)

	c.botModel = botModel
	c.chatGroupID = msg.Chat.ID
	c.msg = update

	input = normalizeExpr(input)

	if !strings.ContainsAny(input, "+-*/") {
		return false
	}

	reg := regexp.MustCompile(`^[0-9\.\+\-\*/\(\) ]+$`)
	if !reg.MatchString(input) {
		return false
	}

	numReg := regexp.MustCompile(`^\s*-?\d+(\.\d+)?\s*$`)
	if numReg.MatchString(input) {
		return false
	}

	return true
}
func (c *CalcHandler) Handle() (err error) {
	input := strings.TrimSpace(c.msg.Message.Text)
	input = normalizeExpr(input)
	result, err := calc(input)
	if err != nil {
		global.GVA_LOG.Error("CalcHandler calc error", zap.Error(err))
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

func calc(expr string) (float64, error) {

	if len(expr) > 80 {
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

func normalizeExpr(expr string) string {

	expr = strings.TrimSpace(expr)

	expr = strings.ReplaceAll(expr, "×", "*")
	expr = strings.ReplaceAll(expr, "✖", "*")
	expr = strings.ReplaceAll(expr, "÷", "/")
	expr = strings.ReplaceAll(expr, "（", "(")
	expr = strings.ReplaceAll(expr, "）", ")")

	return expr
}
