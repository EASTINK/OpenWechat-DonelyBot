package plug

import (
	"encoding/json"
	"fmt"
	"gpt/gtp"
	"log"
	"net/http"
	"openwechat"
	"strings"
	"time"

	//定时
	"github.com/go-co-op/gocron"
)

/*
早安语Api by hin
*/
func linestr() string {
	reg, err := http.Get(`https://apis.tianapi.com/zaoan/index?key=f32a181fb14f87489adc39fd951d4843`)
	if err != nil {
		return ""
	}
	defer reg.Body.Close()

	formData := make(map[string]interface{})
	json.NewDecoder(reg.Body).Decode(&formData)
	//for key, value := range formData {
	return formData["result"].(map[string]interface{})["content"].(string)
}

/*
计算倒计时 by hin
*/
func Studysecond(current, rdate time.Time) int {
	//Current := time.Now()
	//2月
	if current.Month() == 2 {
		if current.Day() <= rdate.Day() {

			if current.Day() == rdate.Day() {
				hour := rdate.Hour() - current.Hour()
				miniute := rdate.Minute() - current.Minute()
				seconds := rdate.Second() - current.Second()
				num := (hour*60+miniute)*60 + seconds
				return num
			}

			if current.Day() < rdate.Day() {
				day := rdate.Day() - current.Day()           //1day
				hour := rdate.Hour() - current.Hour()        //-10h
				miniute := rdate.Minute() - current.Minute() //-22m
				seconds := rdate.Second() - current.Second() //-32s
				num := ((day*24+hour)*60+miniute)*60 + seconds
				return num
			}

		}
	}
	return 0
}

func Cron(groups openwechat.Groups, self *openwechat.Self) {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	s := gocron.NewScheduler(timezone)
	s.Every(1).Day().At("08:30").Do(func() {
		linestr := linestr()
		str := fmt.Sprintf(
			"同学们早安，离寒假结束，新学期伊始还剩下：%d 秒.\n\n%s\n\n这是一条实验性的定时早报,手动查询可发送群指令 #Doomsday\n >退订回#T 祝您依然学习愉快.",
			Studysecond(time.Now(), time.Date(2023, 2, 18, 8, 00, 00, 00, timezone)),
			linestr,
		)
		// 获取所有的群组
		for x := range groups {
			self.SendTextToGroup(groups[x], str)
		}
	})
	s.StartBlocking()
}

func IschatGPT(msg string) string {
	chat := strings.Split(msg, "\n")
	res := ""
	if chat[0] == "#gpt" {
		for x := range chat {
			if x != 0 {
				res += chat[x]
			}
		}
	}
	return res
}

func ChatGPT(question string, msg *openwechat.Message) {
	reply, err := gtp.Completions(question)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("我的妈诶，服务收到不可抗力干扰，请稍后再来尝试吧。")
		return
	}
	if reply == "" {
		log.Printf("没有回应: %v \n", err)
		return
	}
	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		log.Printf("get sender in group error :%v \n", err)
		return
	}
	//格式化回复内容
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	atText := "@" + groupSender.NickName
	replyText := atText + ": \n" + reply
	_, err = msg.ReplyText(replyText)
	if err != nil {
		log.Printf("response group error: %v \n", err)
	}
}

func Doomsday(msg *openwechat.Message) {
	//判断是否是过期消息
	if time.Now().Unix()-msg.CreateTime < 5 {
		if msg.IsText() && msg.IsSendByGroup() {
			switch msg.Content {
			case "#Doomsday":
				timezone, _ := time.LoadLocation("Asia/Shanghai")
				str := fmt.Sprintf("距离工贸开学还有：%d 秒", Studysecond(time.Now(), time.Date(2023, 2, 18, 8, 00, 00, 00, timezone)))
				msg.ReplyText(str)
				break
			case "#T":
				msg.ReplyText("抱歉，该业务还未上线")
				break
			default:
				if IschatGPT(msg.Content) != "" {
					//异步
					msg.ReplyText("正在请求...")
					go ChatGPT(IschatGPT(msg.Content), msg)
				}
			}
		}
	}
}
