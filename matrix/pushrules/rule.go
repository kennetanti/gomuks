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

package pushrules

import (
	"encoding/gob"

	"github.com/tulir/mautrix-go"

	"github.com/kennetanti/gomuks/lib/glob"
)

func init() {
	gob.Register(PushRuleArray{})
	gob.Register(PushRuleMap{})
}

type PushRuleCollection interface {
	GetActions(room Room, event *mautrix.Event) PushActionArray
}

type PushRuleArray []*PushRule

func (rules PushRuleArray) SetType(typ PushRuleType) PushRuleArray {
	for _, rule := range rules {
		rule.Type = typ
	}
	return rules
}

func (rules PushRuleArray) GetActions(room Room, event *mautrix.Event) PushActionArray {
	for _, rule := range rules {
		if !rule.Match(room, event) {
			continue
		}
		return rule.Actions
	}
	return nil
}

type PushRuleMap struct {
	Map  map[string]*PushRule
	Type PushRuleType
}

func (rules PushRuleArray) SetTypeAndMap(typ PushRuleType) PushRuleMap {
	data := PushRuleMap{
		Map:  make(map[string]*PushRule),
		Type: typ,
	}
	for _, rule := range rules {
		rule.Type = typ
		data.Map[rule.RuleID] = rule
	}
	return data
}

func (ruleMap PushRuleMap) GetActions(room Room, event *mautrix.Event) PushActionArray {
	var rule *PushRule
	var found bool
	switch ruleMap.Type {
	case RoomRule:
		rule, found = ruleMap.Map[event.RoomID]
	case SenderRule:
		rule, found = ruleMap.Map[event.Sender]
	}
	if found && rule.Match(room, event) {
		return rule.Actions
	}
	return nil
}

func (ruleMap PushRuleMap) Unmap() PushRuleArray {
	array := make(PushRuleArray, len(ruleMap.Map))
	index := 0
	for _, rule := range ruleMap.Map {
		array[index] = rule
		index++
	}
	return array
}

type PushRuleType string

const (
	OverrideRule  PushRuleType = "override"
	ContentRule   PushRuleType = "content"
	RoomRule      PushRuleType = "room"
	SenderRule    PushRuleType = "sender"
	UnderrideRule PushRuleType = "underride"
)

type PushRule struct {
	// The type of this rule.
	Type PushRuleType `json:"-"`
	// The ID of this rule.
	// For room-specific rules and user-specific rules, this is the room or user ID (respectively)
	// For other types of rules, this doesn't affect anything.
	RuleID string `json:"rule_id"`
	// The actions this rule should trigger when matched.
	Actions PushActionArray `json:"actions"`
	// Whether this is a default rule, or has been set explicitly.
	Default bool `json:"default"`
	// Whether or not this push rule is enabled.
	Enabled bool `json:"enabled"`
	// The conditions to match in order to trigger this rule.
	// Only applicable to generic underride/override rules.
	Conditions []*PushCondition `json:"conditions,omitempty"`
	// Pattern for content-specific push rules
	Pattern string `json:"pattern,omitempty"`
}

func (rule *PushRule) Match(room Room, event *mautrix.Event) bool {
	if !rule.Enabled {
		return false
	}
	switch rule.Type {
	case OverrideRule, UnderrideRule:
		return rule.matchConditions(room, event)
	case ContentRule:
		return rule.matchPattern(room, event)
	case RoomRule:
		return rule.RuleID == event.RoomID
	case SenderRule:
		return rule.RuleID == event.Sender
	default:
		return false
	}
}

func (rule *PushRule) matchConditions(room Room, event *mautrix.Event) bool {
	for _, cond := range rule.Conditions {
		if !cond.Match(room, event) {
			return false
		}
	}
	return true
}

func (rule *PushRule) matchPattern(room Room, event *mautrix.Event) bool {
	pattern, err := glob.Compile(rule.Pattern)
	if err != nil {
		return false
	}
	return pattern.MatchString(event.Content.Body)
}
