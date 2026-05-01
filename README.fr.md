# action1-cli

CLI pour l'API Action1 — gerez vos endpoints, automations, patches et bien plus depuis le terminal.

## Fonctionnalites

- **~140 endpoints API** couverts : endpoints, automations, rapports, deploiement logiciel, vulnerabilites, RBAC, etc.
- **Stockage securise des credentials** : secrets client stockes dans le keychain OS (Keychain macOS, Windows Credential Manager, Linux Secret Service)
- **Multi-profils** : gerez plusieurs comptes et organisations Action1
- **Multi-region** : Amerique du Nord, Europe, Australie
- **Sortie flexible** : table, JSON, CSV, YAML
- **Multi-plateforme** : macOS, Linux, Windows (amd64 et arm64)
- **Pagination automatique** : recuperez tous les resultats avec `--all`
- **Integration MCP** : 50 tools pour Claude Code, VS Code, JetBrains

## Prerequis

- Un compte [Action1](https://www.action1.com) avec un acces API
- Des **credentials API** (Client ID et Client Secret) generes depuis la console Action1

## Installation

### Methode 1 : Telecharger le binaire (recommande)

Rendez-vous sur la page [Releases](https://github.com/dimer47/action1-cli/releases/latest) et telechargez l'archive correspondant a votre plateforme.

Ou en une commande :

**macOS (Apple Silicon — M1/M2/M3/M4) :**

```bash
curl -sL https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_darwin_arm64.tar.gz | tar xz
sudo mv action1 /usr/local/bin/
```

**macOS (Intel) :**

```bash
curl -sL https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_darwin_amd64.tar.gz | tar xz
sudo mv action1 /usr/local/bin/
```

**Linux (amd64) :**

```bash
curl -sL https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_linux_amd64.tar.gz | tar xz
sudo mv action1 /usr/local/bin/
```

**Linux (arm64 — Raspberry Pi, etc.) :**

```bash
curl -sL https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_linux_arm64.tar.gz | tar xz
sudo mv action1 /usr/local/bin/
```

**Windows :**

Telechargez `action1-cli_windows_amd64.zip` depuis les [Releases](https://github.com/dimer47/action1-cli/releases/latest), decompressez et ajoutez le dossier au `PATH`.

### Methode 2 : Depuis les sources (necessite Go 1.21+)

```bash
go install github.com/dimer47/action1-cli@latest
```

Le binaire sera installe dans `$GOPATH/bin/` (generalement `~/go/bin/`). Assurez-vous que ce repertoire est dans votre `PATH`.

### Methode 3 : Compiler localement

```bash
git clone https://github.com/dimer47/action1-cli.git
cd action1-cli
go build -o action1 .
./action1 version
```

### Verifier l'installation

```bash
action1 version
# action1 version 0.1.0
```

## Mise a jour

La CLI verifie automatiquement les nouvelles versions au demarrage et vous notifie quand une mise a jour est disponible.

```bash
# Mettre a jour vers la derniere version
action1 self-update

# Verifier les mises a jour sans installer
action1 self-update --check
```

La mise a jour est telechargee depuis GitHub Releases et remplace le binaire actuel. Si le binaire est dans un repertoire protege (ex: `/usr/local/bin/`), `sudo` sera demande automatiquement.

## Demarrage rapide

### 1. Generer des credentials API

1. Connectez-vous a la [console Action1](https://app.action1.com)
2. Allez dans **Configuration > API Credentials**
3. Cliquez **Create API Credentials**
4. Copiez le **Client ID** et le **Client Secret**
5. Notez votre region (Amerique du Nord, Europe ou Australie)

### 2. S'authentifier

```bash
action1 auth login --region eu
```

Repondez aux questions :
```
Region: eu (https://app.eu.action1.com/api/3.0)
Client ID: api-key-xxxxx@action1.com     # Collez votre Client ID
Client Secret:                             # Collez votre secret (masque)
Successfully authenticated.
Token expires in 3600 seconds.
```

Les credentials client sont stockes dans le **keychain OS** (chiffre). Les tokens sont stockes localement avec des permissions restreintes (0600).

### 3. Configurer l'organisation par defaut

```bash
# Lister les organisations pour trouver l'ID
action1 org list

# Definir l'organisation par defaut
action1 config set org <orgId>
```

### 4. Commencer a utiliser

```bash
# Lister vos endpoints
action1 endpoint list

# Verifier le statut des endpoints
action1 endpoint status

# Lister les vulnerabilites
action1 vulnerability list

# Executer un rapport
action1 report data <reportId>
```

## Configuration

### Fichier de configuration

Situe dans :
- **macOS** : `~/Library/Application Support/action1/config.yaml`
- **Linux** : `~/.config/action1/config.yaml`
- **Windows** : `%APPDATA%\action1\config.yaml`

### Multi-profils

```bash
# Creer un profil "production"
action1 config use-profile production
action1 auth login --region eu

# Creer un profil "staging"
action1 config use-profile staging
action1 auth login --region na

# Voir tous les profils (* = actif)
action1 config list-profiles
#   default
# * production
#   staging

# Changer de profil
action1 config use-profile production

# Utiliser un profil ponctuellement
action1 endpoint list --profile staging
```

### Stockage des credentials

| OS | Backend | Ce qui est stocke |
|----|---------|-------------------|
| macOS | Keychain | Client ID + Secret |
| Windows | Credential Manager | Client ID + Secret |
| Linux | Secret Service (GNOME Keyring / KWallet) | Client ID + Secret |
| Fallback | Fichier config (0600) | Tout |

Les tokens (access + refresh) sont toujours stockes dans le fichier de configuration avec des permissions restreintes, car ils sont trop volumineux pour certaines implementations de keychain.

Utilisez `--no-keychain` pour forcer le stockage fichier pour tout.

## Utilisation

### Endpoints

```bash
action1 endpoint list                              # Lister tous les endpoints
action1 endpoint list --all                        # Tout recuperer (pagination auto)
action1 endpoint list --filter "name eq 'SRV01'"   # Filtre OData
action1 endpoint get <id>                          # Details d'un endpoint
action1 endpoint status                            # Compteurs online/offline
action1 endpoint update <id> --name "nouveau-nom"  # Renommer
action1 endpoint delete <id>                       # Supprimer (demande confirmation)
action1 endpoint move <id> --to-org <orgId>        # Deplacer vers une autre org
action1 endpoint missing-updates <id>              # Patches manquants
action1 endpoint install-url windowsEXE            # URL d'installation de l'agent
```

### Groupes d'endpoints

```bash
action1 endpoint-group list                        # Lister les groupes
action1 endpoint-group create --name "Serveurs"    # Creer
action1 endpoint-group members <id>                # Lister les membres
action1 endpoint-group add <id> --endpoints a,b,c  # Ajouter des endpoints
action1 endpoint-group remove <id> --endpoints a,b # Retirer des endpoints
```

### Automations

```bash
# Planifications
action1 automation schedule list
action1 automation schedule create --data @schedule.json
action1 automation schedule get <id>
action1 automation schedule delete <id>

# Execution immediate
action1 automation instance run --data @automation.json
action1 automation instance results <instanceId>
action1 automation instance stop <instanceId>

# Templates d'actions
action1 automation template list
action1 automation template get <templateId>
```

### Rapports

```bash
action1 report list                                # Lister les rapports
action1 report data <reportId>                     # Lignes du rapport
action1 report data <reportId> --all               # Toutes les lignes
action1 report export <reportId> --output-file r.csv
action1 report requery <reportId>                  # Re-executer
action1 report drilldown <reportId> <rowId>        # Details d'une ligne
```

### Depot logiciel

```bash
action1 software list                              # Lister les packages
action1 software get <id>                          # Details d'un package
action1 software clone <id>                        # Cloner un package
action1 software version create <pkgId> --data @v.json
action1 software upload <pkgId> <verId> ./installer.exe
```

### Mises a jour (Patches)

```bash
action1 update list                                # Toutes les MAJ manquantes
action1 update get <packageId>                     # MAJ d'un package
action1 update endpoints <pkgId> <verId>           # Endpoints manquant une MAJ
```

### Inventaire logiciel installe

```bash
action1 installed-software list                    # Toutes les apps installees
action1 installed-software get <endpointId>        # Apps sur un endpoint
action1 installed-software requery                 # Rafraichir les donnees
```

### Vulnerabilites

```bash
action1 vulnerability list                         # Logiciels vulnerables
action1 vulnerability get <cveId>                  # Details CVE (org)
action1 vulnerability cve <cveId>                  # Details CVE (global)
action1 vulnerability endpoints <cveId>            # Endpoints affectes
action1 vulnerability remediation list <cveId>     # Remediations passees
action1 vulnerability remediation create <cveId> --data @controle.json
```

### Scripts

```bash
action1 script list                                # Lister les scripts
action1 script create --name "Nettoyage" --type powershell --file cleanup.ps1
action1 script get <id>
action1 script update <id> --file cleanup_v2.ps1
action1 script delete <id>
```

### Utilisateurs & RBAC

```bash
action1 user me                                    # Utilisateur courant
action1 user list                                  # Tous les utilisateurs
action1 user create --email admin@co.com --name "Admin"
action1 role list                                  # Tous les roles
action1 role assign <roleId> <userId>              # Assigner un utilisateur
action1 role permissions                           # Templates de permissions
```

### Organisations

```bash
action1 org list                                   # Lister les orgs
action1 org create --name "Production"             # Creer
action1 org update <id> --data '{"name":"Prod"}'   # Modifier
action1 org delete <id>                            # Supprimer
```

### Journal d'audit

```bash
action1 audit list                                 # Evenements d'audit
action1 audit list --from 2026-01-01 --to 2026-01-31
action1 audit get <id>                             # Details d'un evenement
action1 audit export --output-file audit.json      # Exporter
```

### Autres commandes

```bash
action1 search "requete"                           # Recherche rapide
action1 log get                                    # Logs diagnostiques
action1 enterprise get                             # Parametres enterprise
action1 subscription info                          # Info licence
action1 subscription usage                         # Statistiques d'utilisation
```

## Flags globaux

| Flag | Court | Description | Defaut |
|------|-------|-------------|--------|
| `--org` | `-o` | ID de l'organisation | depuis la config |
| `--region` | `-r` | Region serveur : `na`, `eu`, `au` | depuis la config |
| `--output` | `-O` | Format de sortie : `table`, `json`, `csv`, `yaml` | `table` |
| `--profile` | `-p` | Profil de configuration | depuis la config |
| `--quiet` | `-q` | Supprimer les en-tetes et decorations | `false` |
| `--verbose` | `-v` | Afficher les requetes HTTP | `false` |
| `--no-color` | | Desactiver les couleurs | `false` |
| `--no-keychain` | | Forcer le stockage fichier | `false` |
| `--config` | | Chemin du fichier de config | auto-detecte |

## Autocompletion shell

```bash
# Bash
action1 completion bash > /etc/bash_completion.d/action1

# Zsh (ajoutez a votre .zshrc)
action1 completion zsh > "${fpath[1]}/_action1"

# Fish
action1 completion fish > ~/.config/fish/completions/action1.fish

# PowerShell
action1 completion powershell > action1.ps1
```

## Aliases de commandes

Pour une saisie plus rapide :

| Commande | Alias |
|----------|-------|
| `endpoint` | `ep` |
| `endpoint-group` | `epg` |
| `automation` | `auto` |
| `vulnerability` | `vuln` |
| `software` | `sw` |
| `installed-software` | `isw` |
| `data-source` | `ds` |
| `subscription` | `sub` |
| `report-subscription` | `report-sub` |

## Integration MCP (Claude Code, VS Code, JetBrains)

La CLI integre un serveur [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) qui expose **50 tools** directement utilisables par les assistants IA.

### Configuration

Ajoutez dans vos settings Claude Code (ou VS Code / JetBrains avec l'extension Claude) :

```json
{
  "mcpServers": {
    "action1": {
      "command": "action1",
      "args": ["mcp-serve"]
    }
  }
}
```

> Le serveur MCP utilise automatiquement vos credentials depuis `action1 auth login`. Aucune configuration supplementaire necessaire.

### Tools MCP disponibles

| Tool | Description |
|------|-------------|
| `org-list` | Lister les organisations |
| `org-create` | Creer une organisation |
| `endpoint-list` | Lister les endpoints manages |
| `endpoint-get` | Details d'un endpoint |
| `endpoint-status` | Statut (online/offline) |
| `endpoint-update` | Modifier nom/commentaire |
| `endpoint-delete` | Supprimer un endpoint |
| `endpoint-missing-updates` | Patches manquants |
| `endpoint-group-list/get/members` | Groupes d'endpoints |
| `automation-schedule-list/get/delete` | Planifications |
| `automation-instance-list/get/results/stop` | Instances |
| `automation-template-list` | Templates d'actions |
| `report-list/data/export/requery` | Rapports |
| `software-list/get` | Depot logiciel |
| `update-list/get` | Mises a jour manquantes |
| `installed-software-list/get` | Inventaire logiciel |
| `vulnerability-list/get/endpoints/cve` | Vulnerabilites |
| `script-list/get` | Scripts |
| `data-source-list` | Sources de donnees |
| `user-me/list/get/roles` | Utilisateurs |
| `role-list/get/users/permissions` | Roles RBAC |
| `enterprise-get` | Parametres enterprise |
| `subscription-info/usage` | Licences et usage |
| `search` | Recherche |
| `audit-list/get` | Journal d'audit |
| `log-get` | Logs diagnostiques |

### Exemple d'utilisation dans Claude Code

Une fois configure, vous pouvez simplement dire :

- *"Liste mes endpoints Action1"*
- *"Quelles vulnerabilites affectent mon organisation ?"*
- *"Montre-moi les patches manquants pour l'endpoint X"*
- *"Quelles automations sont en cours ?"*

Claude appellera automatiquement les bons tools MCP.

## Developpement

```bash
# Cloner
git clone https://github.com/dimer47/action1-cli.git
cd action1-cli

# Compiler
go build -o action1 .

# Lancer les tests
go test ./...

# Linter
go vet ./...
```

### Creer une nouvelle release

```bash
git tag v0.1.0
git push origin v0.1.0
# GitHub Actions compile et publie automatiquement
```

## Licence

MIT
