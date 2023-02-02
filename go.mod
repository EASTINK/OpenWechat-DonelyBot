module donelybot

go 1.19

require openwechat v0.0.0

require (
	github.com/go-co-op/gocron v1.18.0 // indirect
	gpt v0.0.0-00010101000000-000000000000 // indirect
)

require (
	github.com/robfig/cron/v3 v3.0.1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	plug v0.0.0-00010101000000-000000000000
)

replace openwechat => ../OpenWechat-DonelyBot/openwechat

replace plug => ../OpenWechat-DonelyBot/plug

replace gpt => ../OpenWechat-DonelyBot/gpt
