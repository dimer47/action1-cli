# Action1 CLI - Architecture des commandes

> CLI en Go pour Action1 — Design des commandes, sous-commandes et flags

---

## 1. Informations générales

| Élément | Valeur |
|---------|--------|
| **Binaire** | `action1` |
| **Langage** | Go |
| **Framework CLI** | `cobra` (github.com/spf13/cobra) |
| **Config** | `viper` (github.com/spf13/viper) |
| **Keyring** | `go-keyring` (github.com/zalando/go-keyring) |
| **HTTP** | `net/http` stdlib |
| **Sortie** | JSON, table, CSV, YAML |

---

## 2. Flags globaux

Ces flags sont disponibles sur **toutes** les commandes :

```
--org, -o           ID de l'organisation (override la config)
--region, -r        Région du serveur : na | eu | au (défaut: na)
--output, -O        Format de sortie : table | json | csv | yaml (défaut: table)
--quiet, -q         Mode silencieux (pas de headers, juste les données)
--no-color          Désactiver les couleurs
--verbose, -v       Mode verbeux (affiche les requêtes HTTP)
--config             Chemin du fichier de config (défaut: ~/.config/action1/config.yaml)
--profile, -p       Profil de configuration à utiliser (défaut: "default")
```

---

## 3. Commandes de configuration et authentification

### `action1 auth`

Gestion de l'authentification OAuth 2.0.

```
action1 auth login          Connexion interactive (demande client_id + client_secret)
action1 auth login --client-id <id> --client-secret <secret>
                            Connexion non-interactive
action1 auth logout         Supprime le token du keychain/config
action1 auth status         Affiche le statut d'authentification (token valide, expiration, user)
action1 auth token          Affiche le token d'accès courant (pour piping)
action1 auth refresh        Force le rafraîchissement du token
```

**Flags spécifiques :**
```
--client-id         Client ID de l'API
--client-secret     Client Secret de l'API
--no-keychain       Stocker dans le fichier config au lieu du keychain OS
```

**Stockage du token :**
- macOS : Keychain via `go-keyring`
- Windows : Windows Credential Manager via `go-keyring`
- Linux : Secret Service (GNOME Keyring / KWallet) via `go-keyring`
- Fallback : fichier config avec permissions 0600

---

### `action1 config`

Gestion de la configuration persistante.

```
action1 config init         Initialise un fichier de config interactif
action1 config show         Affiche la configuration courante
action1 config set <key> <value>
                            Définir une valeur (ex: action1 config set org abc123)
action1 config get <key>    Lire une valeur
action1 config unset <key>  Supprimer une valeur
action1 config list-profiles
                            Lister les profils disponibles
action1 config use-profile <name>
                            Changer le profil actif
```

**Clés de configuration :**
```yaml
# ~/.config/action1/config.yaml
current_profile: default
profiles:
  default:
    region: na          # na | eu | au
    org: "abc-123"      # Organisation par défaut
    output: table       # Format de sortie par défaut
    no_color: false
  production:
    region: eu
    org: "prod-456"
    output: json
```

---

## 4. Commandes principales

### 4.1 `action1 endpoint` — Gestion des endpoints

```
action1 endpoint list                   Lister tous les endpoints
action1 endpoint get <id>               Détails d'un endpoint
action1 endpoint status                 Statut des endpoints (online/offline/counts)
action1 endpoint update <id>            Modifier le nom ou commentaire
action1 endpoint delete <id>            Supprimer un endpoint
action1 endpoint move <id> --to-org <orgId>
                                        Déplacer vers une autre org
action1 endpoint missing-updates <id>   Mises à jour manquantes
action1 endpoint install-url <type>     URL d'installation de l'agent
```

**Flags `list` :**
```
--limit, -l     Nombre max de résultats (défaut: 50)
--filter, -f    Filtre OData (ex: --filter "name eq 'SRV01'")
--all           Récupérer tous les résultats (pagination auto)
```

**Flags `update` :**
```
--name          Nouveau nom
--comment       Nouveau commentaire
```

**Flags `install-url` :**
```
<type>          Type d'installation : windowsEXE | windowsMSI | macOS | linux
```

---

### 4.2 `action1 endpoint-group` — Groupes d'endpoints

```
action1 endpoint-group list                     Lister les groupes
action1 endpoint-group create <name>            Créer un groupe
action1 endpoint-group get <id>                 Détails d'un groupe
action1 endpoint-group update <id>              Modifier un groupe
action1 endpoint-group delete <id>              Supprimer un groupe
action1 endpoint-group members <id>             Lister les endpoints du groupe
action1 endpoint-group add <id> --endpoints <id1,id2,...>
                                                Ajouter des endpoints
action1 endpoint-group remove <id> --endpoints <id1,id2,...>
                                                Retirer des endpoints
```

---

### 4.3 `action1 remote-session` — Sessions de bureau à distance

```
action1 remote-session start <endpointId>       Démarrer une session
action1 remote-session get <endpointId> <sessionId>
                                                Détails d'une session
action1 remote-session switch-monitor <endpointId> <sessionId>
                                                Changer de moniteur
```

---

### 4.4 `action1 deployer` — Deployers

```
action1 deployer list                   Lister les deployers
action1 deployer get <id>               Détails d'un deployer
action1 deployer delete <id>            Supprimer un deployer
action1 deployer install-url            URL d'installation du deployer (Windows EXE)
```

---

### 4.5 `action1 agent-deployment` — Configuration du déploiement d'agent

```
action1 agent-deployment get            Paramètres de déploiement
action1 agent-deployment update         Modifier les paramètres
```

---

### 4.6 `action1 automation` — Automations (Schedules + Instances)

#### Sous-commandes schedules :

```
action1 automation schedule list                Lister les schedules
action1 automation schedule create              Créer un schedule
action1 automation schedule get <id>            Détails d'un schedule
action1 automation schedule update <id>         Modifier un schedule
action1 automation schedule delete <id>         Supprimer un schedule
action1 automation schedule deployment-status <id>
                                                Statuts de déploiement
action1 automation schedule remove-action <automationId> <actionId>
                                                Supprimer une action
```

#### Sous-commandes instances :

```
action1 automation instance list                Lister les instances
action1 automation instance run                 Appliquer une automation (run immédiat)
action1 automation instance get <id>            Détails d'une instance
action1 automation instance results <id>        Résultats par endpoint
action1 automation instance result-details <instanceId> <endpointId>
                                                Détails d'un endpoint
action1 automation instance stop <id>           Arrêter une automation
```

#### Sous-commandes action templates :

```
action1 automation template list                Lister les templates d'actions
action1 automation template get <id>            Détails d'un template
```

**Flags `schedule create` / `schedule update` :**
```
--name              Nom du schedule
--actions           Actions (JSON ou fichier @actions.json)
--scope             Scope/cible (JSON ou fichier @scope.json)
--schedule          Planification (JSON ou fichier @schedule.json)
--enabled           Activer/désactiver (true|false)
```

**Flags `instance run` :**
```
--actions           Actions à exécuter (JSON ou @fichier)
--scope             Scope cible (JSON ou @fichier)
--wait, -w          Attendre la fin de l'exécution
--timeout           Timeout en secondes (avec --wait)
```

---

### 4.7 `action1 report` — Rapports

#### Définitions de rapports :

```
action1 report list                             Lister les rapports
action1 report list <categoryId>                Lister par catégorie
action1 report create                           Créer un rapport custom
action1 report update <id>                      Modifier un rapport custom
action1 report delete <id>                      Supprimer un rapport custom
```

#### Données de rapports :

```
action1 report data <reportId>                  Obtenir les lignes du rapport
action1 report errors <reportId>                Erreurs du rapport
action1 report export <reportId>                Exporter un rapport
action1 report requery <reportId>               Re-exécuter un rapport
action1 report drilldown <reportId> <rowId>     Détails d'une ligne
action1 report drilldown-export <reportId> <rowId>
                                                Exporter les détails
```

**Flags `data` :**
```
--limit, -l     Nombre de lignes
--filter, -f    Filtre OData
--all           Toutes les lignes (pagination auto)
```

**Flags `export` :**
```
--format        Format d'export (csv, xlsx)
--output-file   Chemin du fichier de sortie
```

---

### 4.8 `action1 report-subscription` — Abonnements rapports par email

```
action1 report-subscription list                Lister les abonnements
action1 report-subscription create              Créer un abonnement
action1 report-subscription update <id>         Modifier
action1 report-subscription delete <id>         Supprimer
```

**Flags `create` / `update` :**
```
--type          Type de rapport : weekly_patch_statistics | daily_patch_statistics
--enabled       Activer/désactiver (true|false)
```

---

### 4.9 `action1 software` — Dépôt logiciel

```
action1 software list                           Lister les packages
action1 software create                         Créer un package
action1 software get <id>                       Détails d'un package
action1 software update <id>                    Modifier un package
action1 software delete <id>                    Supprimer un package custom
action1 software clone <id>                     Cloner un package
action1 software match-conflicts                Conflits de matching (nouveau)
action1 software match-conflicts <id>           Conflits d'un package existant
```

#### Versions :

```
action1 software version create <packageId>     Créer une version
action1 software version get <packageId> <versionId>
                                                Détails d'une version
action1 software version update <packageId> <versionId>
                                                Modifier une version
action1 software version delete <packageId> <versionId>
                                                Supprimer une version
action1 software version remove-action <packageId> <versionId> <actionId>
                                                Supprimer une action additionnelle
```

#### Upload :

```
action1 software upload <packageId> <versionId> <filePath>
                                                Upload d'un fichier d'installation
```

**Flags `upload` :**
```
--chunk-size    Taille des chunks en Mo (défaut: 10)
--resume        Reprendre un upload interrompu
```

---

### 4.10 `action1 update` — Mises à jour / Patches

```
action1 update list                             Mises à jour manquantes (toutes)
action1 update get <packageId>                  MAJ d'un package spécifique
action1 update endpoints <packageId> <versionId>
                                                Endpoints manquant une MAJ spécifique
```

---

### 4.11 `action1 installed-software` — Inventaire logiciel

```
action1 installed-software list                 Apps installées
action1 installed-software get <endpointId>     Apps d'un endpoint
action1 installed-software errors               Erreurs de collecte
action1 installed-software requery              Re-interroger tous les endpoints
action1 installed-software requery <endpointId> Re-interroger un endpoint
```

---

### 4.12 `action1 vulnerability` — Gestion des vulnérabilités

```
action1 vulnerability list                      Logiciels vulnérables
action1 vulnerability get <cveId>               Détails d'une CVE (org)
action1 vulnerability cve <cveId>               Détails d'une CVE (global, hors org)
action1 vulnerability endpoints <cveId>         Endpoints affectés
action1 vulnerability remediation list <cveId>  Remédiations passées
action1 vulnerability remediation create <cveId>
                                                Documenter un contrôle compensatoire
action1 vulnerability remediation update <cveId> <remediationId>
                                                Modifier une remédiation
action1 vulnerability remediation delete <cveId> <remediationId>
                                                Supprimer une remédiation
```

---

### 4.13 `action1 data-source` — Sources de données

```
action1 data-source list                Lister les sources
action1 data-source create              Créer une source
action1 data-source get <id>            Détails d'une source
action1 data-source update <id>         Modifier une source custom
action1 data-source delete <id>         Supprimer une source custom
```

---

### 4.14 `action1 script` — Bibliothèque de scripts

```
action1 script list                     Lister les scripts
action1 script create                   Créer un script custom
action1 script get <id>                 Détails d'un script
action1 script update <id>              Modifier un script custom
action1 script delete <id>              Supprimer un script custom
```

**Flags `create` / `update` :**
```
--name          Nom du script
--description   Description
--type          Type : powershell | cmd
--content       Contenu du script (inline)
--file, -f      Lire le contenu depuis un fichier
```

---

### 4.15 `action1 setting` — Paramètres avancés

```
action1 setting template list           Lister les templates de paramètres
action1 setting template get <id>       Détails d'un template
action1 setting list                    Lister les paramètres
action1 setting create                  Créer un paramètre
action1 setting get <id>               Détails d'un paramètre
action1 setting update <id>            Modifier un paramètre
action1 setting delete <id>            Supprimer un paramètre
```

---

## 5. Commandes d'administration

### 5.1 `action1 org` — Organisations

```
action1 org list                        Lister les organisations
action1 org create                      Créer une organisation
action1 org update <id>                 Modifier une organisation
action1 org delete <id>                 Supprimer une organisation
```

**Flags `create` / `update` :**
```
--name          Nom de l'organisation
--description   Description
```

---

### 5.2 `action1 user` — Gestion des utilisateurs

```
action1 user me                         Utilisateur courant
action1 user me update                  Modifier l'utilisateur courant
action1 user list                       Lister les utilisateurs
action1 user create                     Créer un utilisateur
action1 user get <id>                   Détails d'un utilisateur
action1 user update <id>                Modifier un utilisateur
action1 user delete <id>                Supprimer un utilisateur
action1 user roles <id>                 Rôles assignés à un utilisateur
```

**Flags `create` / `update` :**
```
--email         Adresse email
--name          Nom complet
--role          ID du rôle à assigner
```

---

### 5.3 `action1 role` — Contrôle d'accès (RBAC)

```
action1 role list                       Lister les rôles
action1 role create                     Créer un rôle
action1 role get <id>                   Détails d'un rôle
action1 role update <id>                Modifier un rôle
action1 role delete <id>                Supprimer un rôle
action1 role clone <id>                 Cloner un rôle
action1 role users <id>                 Utilisateurs d'un rôle
action1 role assign <roleId> <userId>   Assigner un utilisateur
action1 role unassign <roleId> <userId> Retirer un utilisateur
action1 role permissions                Lister les templates de permissions
```

---

### 5.4 `action1 enterprise` — Configuration enterprise

```
action1 enterprise get                  Paramètres enterprise
action1 enterprise update               Modifier les paramètres
action1 enterprise close                Fermer le compte (demande confirmation)
action1 enterprise revoke-closure       Révoquer la fermeture
```

---

### 5.5 `action1 subscription` — Licences et abonnements

```
action1 subscription info               Info licence enterprise
action1 subscription trial              Demander un essai gratuit
action1 subscription quote              Demander un devis
action1 subscription usage              Stats d'utilisation enterprise
action1 subscription usage-orgs         Stats par organisation
action1 subscription usage-org <orgId>  Stats d'une organisation
```

---

## 6. Commandes utilitaires

### 6.1 `action1 search` — Recherche rapide

```
action1 search <query>                  Recherche rapide (rapports, endpoints, apps)
```

**Flags :**
```
--type, -t      Filtrer par type de résultat
```

---

### 6.2 `action1 log` — Logs diagnostiques

```
action1 log get                         Obtenir les logs diagnostiques
```

---

### 6.3 `action1 audit` — Journal d'audit

```
action1 audit list                      Événements d'audit
action1 audit get <id>                  Détails d'un événement
action1 audit export                    Exporter le journal d'audit
```

**Flags `list` :**
```
--limit, -l     Nombre d'événements
--from          Date de début (RFC3339)
--to            Date de fin (RFC3339)
--filter, -f    Filtre
--all           Tous les événements
```

**Flags `export` :**
```
--output-file   Fichier de sortie
--format        Format d'export
```

---

## 7. Commandes spéciales

### `action1 completion`

```
action1 completion bash         Générer le script d'autocomplétion bash
action1 completion zsh          Générer le script d'autocomplétion zsh
action1 completion fish         Générer le script d'autocomplétion fish
action1 completion powershell   Générer le script d'autocomplétion PowerShell
```

### `action1 version`

```
action1 version                 Afficher la version de la CLI
```

---

## 8. Structure du projet Go

```
action1-cli/
├── cmd/
│   └── action1/
│       └── main.go                 Point d'entrée
├── internal/
│   ├── cli/
│   │   ├── root.go                 Commande racine + flags globaux
│   │   ├── auth.go                 action1 auth *
│   │   ├── config.go               action1 config *
│   │   ├── endpoint.go             action1 endpoint *
│   │   ├── endpoint_group.go       action1 endpoint-group *
│   │   ├── remote_session.go       action1 remote-session *
│   │   ├── deployer.go             action1 deployer *
│   │   ├── agent_deployment.go     action1 agent-deployment *
│   │   ├── automation.go           action1 automation *
│   │   ├── report.go               action1 report *
│   │   ├── report_subscription.go  action1 report-subscription *
│   │   ├── software.go             action1 software *
│   │   ├── update.go               action1 update *
│   │   ├── installed_software.go   action1 installed-software *
│   │   ├── vulnerability.go        action1 vulnerability *
│   │   ├── data_source.go          action1 data-source *
│   │   ├── script.go               action1 script *
│   │   ├── setting.go              action1 setting *
│   │   ├── org.go                  action1 org *
│   │   ├── user.go                 action1 user *
│   │   ├── role.go                 action1 role *
│   │   ├── enterprise.go           action1 enterprise *
│   │   ├── subscription.go         action1 subscription *
│   │   ├── search.go               action1 search
│   │   ├── log.go                  action1 log *
│   │   └── audit.go                action1 audit *
│   ├── api/
│   │   ├── client.go               Client HTTP (base URL, auth, pagination)
│   │   ├── oauth.go                Gestion OAuth 2.0 (token, refresh)
│   │   ├── endpoints.go            Appels API endpoints
│   │   ├── automations.go          Appels API automations
│   │   ├── reports.go              Appels API reports
│   │   ├── software.go             Appels API software repository
│   │   ├── vulnerabilities.go      Appels API vulnerabilities
│   │   ├── users.go                Appels API users/roles
│   │   ├── organizations.go        Appels API organizations
│   │   └── ...                     (un fichier par domaine API)
│   ├── auth/
│   │   ├── store.go                Interface de stockage de credentials
│   │   ├── keyring.go              Implémentation go-keyring (macOS/Win/Linux)
│   │   └── file.go                 Implémentation fichier (fallback)
│   ├── config/
│   │   ├── config.go               Gestion de la config YAML + profiles
│   │   └── paths.go                Chemins par défaut selon l'OS
│   ├── output/
│   │   ├── formatter.go            Interface de formatage
│   │   ├── table.go                Sortie en table
│   │   ├── json.go                 Sortie JSON
│   │   ├── csv.go                  Sortie CSV
│   │   └── yaml.go                 Sortie YAML
│   └── upload/
│       └── chunked.go              Upload multi-chunk resumable
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── docs/
    ├── api-analysis.md
    └── cli-design.md
```

---

## 9. Récapitulatif des commandes de premier niveau

| Commande | Nb sous-commandes | Domaine |
|----------|:-----------------:|---------|
| `auth` | 5 | Authentification |
| `config` | 7 | Configuration |
| `endpoint` | 8 | Gestion machines |
| `endpoint-group` | 8 | Groupes de machines |
| `remote-session` | 3 | Bureau à distance |
| `deployer` | 4 | Service deployer |
| `agent-deployment` | 2 | Config déploiement |
| `automation` | 13 | Schedules + Instances + Templates |
| `report` | 10 | Rapports + Données |
| `report-subscription` | 4 | Abonnements email |
| `software` | 14 | Dépôt logiciel + Versions + Upload |
| `update` | 3 | Patches manquants |
| `installed-software` | 5 | Inventaire logiciel |
| `vulnerability` | 10 | CVE + Remédiations |
| `data-source` | 5 | Sources de données |
| `script` | 5 | Bibliothèque scripts |
| `setting` | 7 | Paramètres avancés |
| `org` | 4 | Organisations |
| `user` | 8 | Utilisateurs |
| `role` | 10 | RBAC |
| `enterprise` | 4 | Config enterprise |
| `subscription` | 6 | Licences + Usage |
| `search` | 1 | Recherche |
| `log` | 1 | Logs diagnostiques |
| `audit` | 3 | Journal d'audit |
| `completion` | 4 | Autocomplétion |
| `version` | 1 | Version |
| **Total** | **~155** | |

---

## 10. Conventions de design

### Nommage
- Commandes en **kebab-case** : `endpoint-group`, `remote-session`, `installed-software`
- Flags en **kebab-case** : `--client-id`, `--output-file`, `--chunk-size`
- Flags courts : une lettre, quand non ambigu

### Patterns CRUD cohérents
Pour chaque ressource, les sous-commandes suivent le même pattern :
```
list        → GET collection
get <id>    → GET resource
create      → POST collection
update <id> → PATCH resource
delete <id> → DELETE resource
```

### Entrée de données complexes
Les payloads JSON peuvent être fournis de 3 manières :
```
--data '{"key": "value"}'      Inline JSON
--data @payload.json            Depuis un fichier
--data -                        Depuis stdin
```

### Opérations destructives
Les commandes `delete`, `enterprise close` demandent confirmation sauf avec `--yes, -y`.

### Pagination automatique
Le flag `--all` désactive la pagination et itère automatiquement sur toutes les pages.

### Sortie machine-friendly
Avec `--output json --quiet`, la CLI produit du JSON brut sans décoration, idéal pour le piping :
```bash
action1 endpoint list -O json -q | jq '.[].name'
```
