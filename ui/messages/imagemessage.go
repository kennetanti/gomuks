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
	"bytes"
	"fmt"
	"image/color"

	"github.com/tulir/mautrix-go"
	"github.com/tulir/tcell"

	"github.com/kennetanti/gomuks/config"
	"github.com/kennetanti/gomuks/debug"
	"github.com/kennetanti/gomuks/interface"
	"github.com/kennetanti/gomuks/lib/ansimage"
	"github.com/kennetanti/gomuks/ui/messages/tstring"
)

type ImageMessage struct {
	BaseMessage
	Body       string
	Homeserver string
	FileID     string
	data       []byte

	matrix ifc.MatrixContainer
}

// NewImageMessage creates a new ImageMessage object with the provided values and the default state.
func NewImageMessage(matrix ifc.MatrixContainer, event *mautrix.Event, displayname string, body, homeserver, fileID string, data []byte) UIMessage {
	return &ImageMessage{
		newBaseMessage(event, displayname),
		body,
		homeserver,
		fileID,
		data,
		matrix,
	}
}

func (msg *ImageMessage) Clone() UIMessage {
	data := make([]byte, len(msg.data))
	copy(data, msg.data)
	return &ImageMessage{
		BaseMessage: msg.BaseMessage.clone(),
		Body:        msg.Body,
		Homeserver:  msg.Homeserver,
		FileID:      msg.FileID,
		data:        data,
		matrix:      msg.matrix,
	}
}

func (msg *ImageMessage) RegisterMatrix(matrix ifc.MatrixContainer) {
	msg.matrix = matrix

	if len(msg.data) == 0 {
		go msg.updateData()
	}
}

func (msg *ImageMessage) NotificationContent() string {
	return "Sent an image"
}

func (msg *ImageMessage) PlainText() string {
	return fmt.Sprintf("%s: %s", msg.Body, msg.matrix.GetDownloadURL(msg.Homeserver, msg.FileID))
}

func (msg *ImageMessage) updateData() {
	defer debug.Recover()
	debug.Print("Loading image:", msg.Homeserver, msg.FileID)
	data, _, _, err := msg.matrix.Download(fmt.Sprintf("mxc://%s/%s", msg.Homeserver, msg.FileID))
	if err != nil {
		debug.Printf("Failed to download image %s/%s: %v", msg.Homeserver, msg.FileID, err)
		return
	}
	debug.Print("Image", msg.Homeserver, msg.FileID, "loaded.")
	msg.data = data
}

func (msg *ImageMessage) Path() string {
	return msg.matrix.GetCachePath(msg.Homeserver, msg.FileID)
}

// CalculateBuffer generates the internal buffer for this message that consists
// of the text of this message split into lines at most as wide as the width
// parameter.
func (msg *ImageMessage) CalculateBuffer(prefs config.UserPreferences, width int) {
	if width < 2 {
		return
	}
	msg.CalculateReplyBuffer(prefs, width)

	if prefs.BareMessageView || prefs.DisableImages {
		msg.calculateBufferWithText(prefs, tstring.NewTString(msg.PlainText()), width)
		return
	}

	image, err := ansimage.NewScaledFromReader(bytes.NewReader(msg.data), 0, width, color.Black)
	if err != nil {
		msg.buffer = []tstring.TString{tstring.NewColorTString("Failed to display image", tcell.ColorRed)}
		debug.Print("Failed to display image:", err)
		return
	}

	msg.buffer = image.Render()
}
