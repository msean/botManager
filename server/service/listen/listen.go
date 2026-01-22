package listen

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/listen"
)

type ListenSvc struct{}

func (svc *ListenSvc) Choice(ctx context.Context) (list []listen.BotChatMap, err error) {
	err = global.GVA_MYSQL.
		Model(&listen.BotChatMap{}).
		Order("updated_at desc").
		Find(&list).Error

	return
}

func (svc *ListenSvc) Query(
	ctx context.Context,
	req listen.ListenQueryReq,
) (list []listen.BotMessageVO, total int64, err error) {

	table := fmt.Sprintf(`bot_messages_%d`, abs(req.GroupID))
	offset := (req.Page - 1) * req.PageSize

	where := "1=1"
	args := []interface{}{}

	if req.Keyword != "" {
		where += " AND (text ILIKE ? OR caption ILIKE ?)"
		args = append(args, "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.StartTime != nil {
		where += " AND timestamp >= ?"
		args = append(args, *req.StartTime)
	}

	if req.EndTime != nil {
		where += " AND timestamp <= ?"
		args = append(args, *req.EndTime)
	}

	// ===== count =====
	countSQL := fmt.Sprintf(`SELECT count(*) FROM "%s" WHERE %s`, table, where)
	if err = global.GVA_MYSQL.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return
	}

	// ===== list =====
	sql := fmt.Sprintf(`
		SELECT
			m.id,
			? AS group_id,
			c.group_name,
			m.user_id,
			m.username,
			m.nick_name,
			m.message_type,
			m.text,
			m.caption,
			m.timestamp
		FROM "%s" m
		LEFT JOIN bot_chat_map c ON c.group_id = ?
		WHERE %s
		ORDER BY m.timestamp DESC
		LIMIT ? OFFSET ?
	`, table, where)

	queryArgs := append(
		[]interface{}{req.GroupID, req.GroupID},
		append(args, req.PageSize, offset)...,
	)

	err = global.GVA_MYSQL.Raw(sql, queryArgs...).Scan(&list).Error
	return
}

func (svc *ListenSvc) Export(
	ctx context.Context,
	req listen.ListenQueryReq,
) (string, error) {

	table := fmt.Sprintf(`bot_messages_%d`, abs(req.GroupID))

	where := "1=1"
	args := []interface{}{}

	if req.Keyword != "" {
		where += " AND (text ILIKE ? OR caption ILIKE ?)"
		args = append(args, "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	sql := fmt.Sprintf(`
		SELECT
			id,
			username,
			nick_name,
			text,
			caption,
			timestamp
		FROM "%s"
		WHERE %s
		ORDER BY timestamp DESC
		LIMIT 50000
	`, table, where)

	rows, err := global.GVA_MYSQL.Raw(sql, args...).Rows()
	if err != nil {
		return "", err
	}
	defer rows.Close()

	file := fmt.Sprintf("export_%d_%d.csv", req.GroupID, time.Now().Unix())
	f, _ := os.Create(file)
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Write([]string{"ID", "Username", "NickName", "Text", "Caption", "Time"})

	for rows.Next() {
		var r listen.BotMessageVO
		global.GVA_MYSQL.ScanRows(rows, &r)

		writer.Write([]string{
			strconv.FormatInt(r.ID, 10),
			r.Username,
			r.NickName,
			r.Text,
			r.Caption,
			r.Timestamp.Format("2006-01-02 15:04:05"),
		})
	}

	writer.Flush()
	return file, nil
}

func abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}
