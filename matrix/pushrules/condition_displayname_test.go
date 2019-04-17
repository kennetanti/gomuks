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

package pushrules_test

import (
	"github.com/tulir/mautrix-go"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushCondition_Match_DisplayName(t *testing.T) {
	event := newFakeEvent(mautrix.EventMessage, mautrix.Content{
		MsgType: mautrix.MsgText,
		Body:    "tulir: test mention",
	})
	event.Sender = "@someone_else:matrix.org"
	assert.True(t, displaynamePushCondition.Match(displaynameTestRoom, event))
}

func TestPushCondition_Match_DisplayName_Fail(t *testing.T) {
	event := newFakeEvent(mautrix.EventMessage, mautrix.Content{
		MsgType: mautrix.MsgText,
		Body:    "not a mention",
	})
	event.Sender = "@someone_else:matrix.org"
	assert.False(t, displaynamePushCondition.Match(displaynameTestRoom, event))
}

func TestPushCondition_Match_DisplayName_CantHighlightSelf(t *testing.T) {
	event := newFakeEvent(mautrix.EventMessage, mautrix.Content{
		MsgType: mautrix.MsgText,
		Body:    "tulir: I can't highlight myself",
	})
	assert.False(t, displaynamePushCondition.Match(displaynameTestRoom, event))
}

func TestPushCondition_Match_DisplayName_FailsOnEmptyRoom(t *testing.T) {
	emptyRoom := newFakeRoom(0)
	event := newFakeEvent(mautrix.EventMessage, mautrix.Content{
		MsgType: mautrix.MsgText,
		Body:    "tulir: this room doesn't have the owner Member available, so it fails.",
	})
	event.Sender = "@someone_else:matrix.org"
	assert.False(t, displaynamePushCondition.Match(emptyRoom, event))
}
