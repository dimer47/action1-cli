package api

import "testing"

func TestEscapeID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "plain id is unchanged",
			id:   "Correctif_SSH___PC_CEDRINE_1789119931048_2026-09-11_09-47-30",
			want: "Correctif_SSH___PC_CEDRINE_1789119931048_2026-09-11_09-47-30",
		},
		{
			name: "timezone slash is double-encoded",
			id:   "Maintien_SSH_1788991709904_2026-09-14_11-00-00_Europe/Paris",
			want: "Maintien_SSH_1788991709904_2026-09-14_11-00-00_Europe%252FParis",
		},
		{
			name: "empty id stays empty",
			id:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeID(tt.id); got != tt.want {
				t.Errorf("EscapeID(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}
