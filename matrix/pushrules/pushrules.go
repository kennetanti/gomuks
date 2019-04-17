package pushrules

import (
	"encoding/json"
	"net/url"

	"github.com/tulir/mautrix-go"
)

// GetPushRules returns the push notification rules for the global scope.
func GetPushRules(client *mautrix.Client) (*PushRuleset, error) {
	return GetScopedPushRules(client, "global")
}

// GetScopedPushRules returns the push notification rules for the given scope.
func GetScopedPushRules(client *mautrix.Client, scope string) (resp *PushRuleset, err error) {
	u, _ := url.Parse(client.BuildURL("pushrules", scope))
	// client.BuildURL returns the URL without a trailing slash, but the pushrules endpoint requires the slash.
	u.Path += "/"
	_, err = client.MakeRequest("GET", u.String(), nil, &resp)
	return
}

type contentWithRuleset struct {
	Ruleset *PushRuleset `json:"global"`
}

// EventToPushRules converts a m.push_rules event to a PushRuleset by passing the data through JSON.
func EventToPushRules(event *mautrix.Event) (*PushRuleset, error) {
	content := &contentWithRuleset{}
	err := json.Unmarshal(event.Content.VeryRaw, content)
	if err != nil {
		return nil, err
	}

	return content.Ruleset, nil
}
