package goroku

import "testing"

func TestSecurityCheckDoesNotReloadRightsEveryCall(t *testing.T) {
	db := NewDatabase(42)
	db.data["goroku.security"] = map[string]interface{}{
		"owner":     []interface{}{int64(42)},
		"all_users": []interface{}{},
	}
	db.data["goroku.main"] = map[string]interface{}{
		"command_prefixes": map[string]interface{}{
			"999": []interface{}{"."},
		},
	}

	sm := NewSecurityManager(&CustomTelegramClient{TGID: 42}, db)

	// Simulate prefixes being written after startup. Check() is a hot path and
	// must not run reloadRights()/cleanup on every message.
	db.data["goroku.main"]["command_prefixes"] = map[string]interface{}{
		"999": []interface{}{"."},
	}

	if !sm.Check(&Message{SenderID: 42, ChatID: 1, Out: true}, "ping") {
		t.Fatal("owner/outgoing message should pass security check")
	}

	prefixes, ok := db.Get("goroku.main", "command_prefixes", nil).(map[string]interface{})
	if !ok {
		t.Fatalf("command_prefixes has unexpected type: %T", db.Get("goroku.main", "command_prefixes", nil))
	}
	if _, ok := prefixes["999"]; !ok {
		t.Fatal("Check() reloaded rights and cleaned command_prefixes on the hot path")
	}
}
