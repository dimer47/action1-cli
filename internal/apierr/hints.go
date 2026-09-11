// Package apierr enrichit les erreurs renvoyées par l'API Action1 avec des
// indications pratiques.
//
// L'API rejette certaines valeurs ou champs sans expliquer ce qu'elle attend.
// Les indications ci-dessous viennent d'essais réels sur l'API de production ;
// elles ne remplacent pas sa validation, elles la commentent.
//
// Principe : ne jamais bloquer un appel en amont. L'API reste seule autorité —
// la liste des valeurs acceptées peut évoluer, et une validation locale trop
// stricte refuserait à tort des valeurs valides.
//
// Convention de fiabilité : une indication issue d'un constat empirique non
// confirmé par la documentation est préfixée de [observé]. Le lecteur sait
// ainsi ce qui est certain et ce qui reste une piste. Sans préfixe, le fait a
// été vérifié directement (message explicite de l'API, ou test reproductible).
package apierr

import (
	"fmt"
	"strings"
)

// hint associe un motif détecté dans le message d'erreur à une explication.
type hint struct {
	// match renvoie true si l'indication s'applique au message (déjà en minuscules).
	match func(lowerMsg string) bool
	// text est l'explication ajoutée sous l'erreur brute.
	text string
}

// containsAll renvoie un prédicat vrai quand tous les fragments sont présents.
func containsAll(fragments ...string) func(string) bool {
	return func(msg string) bool {
		for _, f := range fragments {
			if !strings.Contains(msg, f) {
				return false
			}
		}
		return true
	}
}

// hints est consulté dans l'ordre ; la première correspondance gagne.
var hints = []hint{
	{
		// Le message de l'API couvre le dépassement de maximum. On ajoute les
		// repères de conversion et l'autre cause de refus, moins évidente.
		match: containsAll("retry_minutes"),
		text: `Fenêtre pendant laquelle un poste hors tension exécutera la tâche à sa
reconnexion.
Repères : 1440 = 1 jour · 10080 = 7 j · 20160 = 14 j · 43200 = 30 j (maximum).

Si la valeur respecte le maximum, le refus vient probablement d'ailleurs :
[observé] une automatisation déjà exécutée (last_run renseigné) semble refuser
toute modification de retry_minutes — recréer l'automatisation dans ce cas.`,
	},
	{
		match: containsAll("randomize_start"),
		text: `randomize_start est en lecture seule : le retirer du payload de création.
Il apparaît dans les réponses de l'API mais ne peut pas être défini.`,
	},
	{
		match: containsAll("success_codes"),
		text:  `success_codes n'est pas accepté à la création d'un script : le retirer du payload.`,
	},
	{
		// L'API renvoie soit la regex brute (sans jamais nommer "settings"),
		// soit le libellé court "Invalid schedule settings".
		match: func(msg string) bool {
			return strings.Contains(msg, "invalid schedule") ||
				(strings.Contains(msg, "regex") && strings.Contains(msg, "enabled|disabled"))
		},
		text: `Format attendu pour settings :
    ENABLED|DISABLED  <fréquence>  AT:hh-mm-ss  [DATE:AAAA-MM-JJ]
Fréquences : ONCE[:AAAA-MM-JJ] · EVERY:<n> · WEEKLY:<jours> · MONTHLY:<n>
             MONTHLYWEEK:<1-5>:<jour>
Jours : Sun Mon Tue Wed Thu Fri Sat — séparés par des virgules, sans espace.
Exemples :
    ENABLED ONCE AT:22-13-38 DATE:2026-09-09
    ENABLED WEEKLY:Mon AT:11-00-00
    ENABLED WEEKLY:Mon,Thu AT:07-30-00
Attention : WEEKLY:Mon, et non DAYS:MON.
Pour un horaire stable été comme hiver, utiliser settings_timezone: "Europe/Paris"
plutôt que "UTC".`,
	},
	{
		match: containsAll("script_text"),
		text:  `Le corps d'un script se passe dans script_text (et non body).`,
	},
	{
		match: func(msg string) bool {
			return strings.Contains(msg, "platform") && strings.Contains(msg, "must be set")
		},
		text: `Champs obligatoires à la création d'un script :
    name · platform ("Windows") · language ("PowerShell") · status ("Published") · script_text
[observé] language semble devoir valoir "PowerShell" et non "PowerShell - Windows" :
cette dernière valeur a coïncidé avec une erreur 500 dont le message ne désigne
aucun champ.`,
	},
	{
		match: func(msg string) bool {
			return strings.Contains(msg, "status") && strings.Contains(msg, "must be set")
		},
		text: `status est obligatoire à la création d'un script : utiliser "Published".`,
	},
	{
		match: containsAll("language"),
		text: `[observé] language semble attendre "PowerShell" et non "PowerShell - Windows" :
cette dernière valeur a coïncidé avec une erreur 500 dont le message ne désigne
aucun champ.`,
	},
}

// Hint renvoie l'indication correspondant au message d'erreur, ou "" s'il n'y en a pas.
func Hint(message string) string {
	lower := strings.ToLower(message)
	for _, h := range hints {
		if h.match(lower) {
			return h.text
		}
	}
	return ""
}

// Wrap met en forme une erreur API en y ajoutant, le cas échéant, une indication.
func Wrap(statusCode int, message string) error {
	if h := Hint(message); h != "" {
		return fmt.Errorf("API error %d: %s\n\n%s", statusCode, message, h)
	}
	return fmt.Errorf("API error %d: %s", statusCode, message)
}
