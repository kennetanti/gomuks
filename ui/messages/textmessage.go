// gomuks - A terminal Matrix client written in Go.
// Copyright (C) 2019 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package messages

import (
	"fmt"
	"time"

	"github.com/tulir/mautrix-go"

	"github.com/kennetanti/gomuks/config"
	"github.com/kennetanti/gomuks/ui/messages/tstring"
)

type TextMessage struct {
	BaseMessage
	cache   tstring.TString
	MsgText string
}

// NewTextMessage creates a new UITextMessage object with the provided values and the default state.
func NewTextMessage(event *mautrix.Event, displayname string, text string) UIMessage {
	return &TextMessage{
		BaseMessage: newBaseMessage(event, displayname),
		MsgText:     text,
	}
}

func NewServiceMessage(text string) UIMessage {
	return &TextMessage{
		BaseMessage: BaseMessage{
			MsgSenderID:  "*",
			MsgSender:    "*",
			MsgTimestamp: time.Now(),
			MsgIsService: true,
		},
		MsgText: text,
	}
}

func (msg *TextMessage) Clone() UIMessage {
	return &TextMessage{
		BaseMessage: msg.BaseMessage.clone(),
		MsgText:     msg.MsgText,
	}
}

func (msg *TextMessage) getCache() tstring.TString {
	if msg.cache == nil {
		switch msg.MsgType {
		case "m.emote":
			msg.cache = tstring.NewColorTString(fmt.Sprintf("* %s %s", msg.MsgSender, msg.MsgText), msg.TextColor())
			msg.cache.Colorize(0, len(msg.MsgSender)+2, msg.SenderColor())
		default:
			msg.cache = tstring.NewColorTString(msg.MsgText, msg.TextColor())
		}
	}
	return msg.cache
}

func (msg *TextMessage) SetIsHighlight(isHighlight bool) {
	msg.BaseMessage.SetIsHighlight(isHighlight)
	msg.cache = nil
}

func (msg *TextMessage) NotificationContent() string {
	return msg.MsgText
}

func (msg *TextMessage) PlainText() string {
	return msg.MsgText
}

func (msg *TextMessage) CalculateBuffer(prefs config.UserPreferences, width int) {
	msg.CalculateReplyBuffer(prefs, width)
	msg.calculateBufferWithText(prefs, msg.getCache(), width)
}
