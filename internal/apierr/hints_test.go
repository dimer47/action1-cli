package apierr

import "strings"

import "testing"

func TestHintRetryMinutes(t *testing.T) {
	msg := "Property name 'retry_minutes' is invalid or cannot be set."
	h := Hint(msg)
	if h == "" {
		t.Fatal("aucune indication pour retry_minutes")
	}
	// La valeur maximale doit apparaître : c'est l'information la plus utile.
	if !strings.Contains(h, "43200") {
		t.Errorf("le maximum 43200 manque dans l'indication:\n%s", h)
	}
}

func TestHintSettingsRegex(t *testing.T) {
	msg := `Regex format: ^(ENABLED|DISABLED)\s?((EVERY:(\d+)|ONCE...`
	h := Hint(msg)
	if h == "" {
		t.Fatal("aucune indication pour le format settings")
	}
	if !strings.Contains(h, "WEEKLY:Mon") {
		t.Errorf("la syntaxe WEEKLY:Mon manque dans l'indication:\n%s", h)
	}
}

func TestHintSettingsUserMessage(t *testing.T) {
	// L'API renvoie parfois le libellé court plutôt que la regex.
	if Hint("Invalid schedule settings") == "" {
		t.Error("aucune indication pour \"Invalid schedule settings\"")
	}
}

func TestHintScriptCreationFields(t *testing.T) {
	cases := map[string]string{
		"Invalid request: Property name 'platform' must be set when creating an object":    "platform",
		"Invalid request: Property name 'script_text' must be set when creating an object": "script_text",
		"Property name 'success_codes' does not supported":                                 "success_codes",
		"Property name 'randomize_start' is invalid or cannot be set.":                     "randomize_start",
	}
	for msg, label := range cases {
		if Hint(msg) == "" {
			t.Errorf("aucune indication pour %s (message: %q)", label, msg)
		}
	}
}

func TestHintUnknownMessage(t *testing.T) {
	// Un message sans rapport ne doit produire aucune indication : mieux vaut
	// pas d'aide qu'une aide hors sujet.
	if h := Hint("Endpoint not found"); h != "" {
		t.Errorf("indication inattendue:\n%s", h)
	}
}

func TestWrapWithoutHint(t *testing.T) {
	err := Wrap(404, "Endpoint not found")
	want := "API error 404: Endpoint not found"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestWrapWithHint(t *testing.T) {
	err := Wrap(400, "Property name 'retry_minutes' is invalid or cannot be set.")
	s := err.Error()
	// Le message brut de l'API doit rester intact en tête.
	if !strings.HasPrefix(s, "API error 400: Property name 'retry_minutes'") {
		t.Errorf("le message d'origine a été altéré:\n%s", s)
	}
	if !strings.Contains(s, "43200") {
		t.Errorf("l'indication n'a pas été ajoutée:\n%s", s)
	}
}

func TestHintIsCaseInsensitive(t *testing.T) {
	if Hint("PROPERTY NAME 'RETRY_MINUTES' IS INVALID") == "" {
		t.Error("la détection devrait ignorer la casse")
	}
}
