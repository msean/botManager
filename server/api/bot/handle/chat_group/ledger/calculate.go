package ledger

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

	// 匹配表达式：数字 + 运算符 + 数字
	reg := regexp.MustCompile(`^\s*(-?\d+(\.\d+)?)\s*([\+\-\*/])\s*(-?\d+(\.\d+)?)\s*$`)

	global.GVA_LOG.Debug("CalcHandler", zap.Any(">>>>>>>", reg.MatchString(input)))
	if reg.MatchString(input) {
		return true
	}

	return false
}

func (c *CalcHandler) Handle() (err error) {

	input := strings.TrimSpace(c.msg.Message.Text)

	result, err := calc(input)
	if err != nil {
		return
	}

	botSender, err := botapi.NewBotAPI(c.botModel.Token)
	if err != nil {
		return
	}

	reply := fmt.Sprintf("%v", utils.FloatReserve(result, 2))

	msg := botapi.NewMessage(c.chatGroupID, reply)
	msg.ReplyToMessageID = c.msg.Message.MessageID

	_, err = botSender.Send(msg)
	return
}

func calc(expr string) (float64, error) {

	reg := regexp.MustCompile(`^\s*(-?\d+(\.\d+)?)\s*([\+\-\*/])\s*(-?\d+(\.\d+)?)\s*$`)
	match := reg.FindStringSubmatch(expr)

	if len(match) < 5 {
		return 0, fmt.Errorf("invalid expression")
	}

	a, _ := strconv.ParseFloat(match[1], 64)
	op := match[3]
	b, _ := strconv.ParseFloat(match[4], 64)

	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}

	return 0, fmt.Errorf("unsupported operator")
}
