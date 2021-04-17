package endpoint

import (
	"github.com/slack-go/slack"
)

func CommandList() slack.MsgOption {

	headerText := slack.NewTextBlockObject("mrkdwn", "お呼びでございましょうか\n*ご用件は？*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	lightOnBtnTxt := slack.NewTextBlockObject("plain_text", "点灯", false, false)
	lightOnBtn := slack.NewButtonBlockElement("", "light:on", lightOnBtnTxt)

	lightOffBtnTxt := slack.NewTextBlockObject("plain_text", "消灯", false, false)
	lightOffBtn := slack.NewButtonBlockElement("", "light:off", lightOffBtnTxt)

	lockBtnTxt := slack.NewTextBlockObject("plain_text", "施錠", false, false)
	lockBtn := slack.NewButtonBlockElement("", "key:lock", lockBtnTxt)

	unlockBtnTxt := slack.NewTextBlockObject("plain_text", "解錠", false, false)
	unlockBtn := slack.NewButtonBlockElement("", "key:unlock", unlockBtnTxt)

	actionBlock := slack.NewActionBlock("", lightOnBtn, lightOffBtn, lockBtn, unlockBtn)

	msg := slack.MsgOptionBlocks(
		headerSection,
		actionBlock,
	)
	return msg
}
