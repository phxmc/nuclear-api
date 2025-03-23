package domain

import "testing"

func TestTelegram(t *testing.T) {
	t.Run("ENTER_EMAIL", func(t *testing.T) {
		state := StateEnterEmail
		got := state.Valid()
		if !got {
			t.Fatalf("want: true, got: %v\n", got)
		}
	})

	t.Run("missing", func(t *testing.T) {
		state := ChatState("MISSING_CHAT_STATE")
		got := state.Valid()
		if got {
			t.Fatalf("want: false, got: %v\n", got)
		}
	})
}
