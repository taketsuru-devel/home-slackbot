package endpoint

import (
	"github.com/slack-go/slack"
)

func commandList() slack.MsgOption {

	headerText := slack.NewTextBlockObject("mrkdwn", "お呼びでございましょうか\n*ご用件は？*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	lightOnBtnTxt := slack.NewTextBlockObject("plain_text", "点灯", false, false)
	lightOnBtn := slack.NewButtonBlockElement("", "iot:light:on", lightOnBtnTxt)

	lightOffBtnTxt := slack.NewTextBlockObject("plain_text", "消灯", false, false)
	lightOffBtn := slack.NewButtonBlockElement("", "iot:light:off", lightOffBtnTxt)

	lockBtnTxt := slack.NewTextBlockObject("plain_text", "施錠", false, false)
	lockBtn := slack.NewButtonBlockElement("", "iot:key:lock", lockBtnTxt)

	unlockBtnTxt := slack.NewTextBlockObject("plain_text", "解錠", false, false)
	unlockBtn := slack.NewButtonBlockElement("", "iot:key:unlock", unlockBtnTxt)

	devWakeupBtnTxt := slack.NewTextBlockObject("plain_text", "開発マシン起動", false, false)
	devWakeupBtn := slack.NewButtonBlockElement("", "ec2:dev:start", devWakeupBtnTxt)

	devDownBtnTxt := slack.NewTextBlockObject("plain_text", "開発マシン終了", false, false)
	devDownBtn := slack.NewButtonBlockElement("", "ec2:dev:stop", devDownBtnTxt)

	actionBlock := slack.NewActionBlock("", lightOnBtn, lightOffBtn, lockBtn, unlockBtn, devWakeupBtn, devDownBtn)

	msg := slack.MsgOptionBlocks(
		headerSection,
		actionBlock,
	)
	return msg
}
