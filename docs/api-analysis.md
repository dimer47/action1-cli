# Action1 API - Analyse complète

> Document de référence pour la conception de la CLI Action1
> API Version: 3.1.0 | OpenAPI Specification 3.1

---

## 1. Informations générales

| Élément | Valeur |
|---------|--------|
| **Nom** | Action1 API |
| **Version** | 3.1.0 |
| **Spec** | OpenAPI 3.1 |
| **Format** | REST / JSON |
| **Chemins** | 96 |
| **Opérations** | ~140 |
| **Schémas** | 171 |
| **Documentation** | https://app.action1.com/apidocs/ |

### Serveurs régionaux (Base URLs)

| Région | URL |
|--------|-----|
| Amérique du Nord (défaut) | `https://app.action1.com/api/3.0` |
| Europe | `https://app.eu.action1.com/api/3.0` |
| Australie | `https://app.au.action1.com/api/3.0` |

### Authentification

- **Méthode** : OAuth 2.0 (Client Credentials)
- **Endpoint** : `POST /oauth2/token`
- **Paramètres** : `client_id` + `client_secret`
- **Réponse** : `access_token` (expire après 3600s) + `refresh_token`
- **Usage** : Header `Authorization: Bearer <access_token>`

### Pagination

- Pagination par curseur : paramètres `limit`, `next_page`, `prev_page`

### Architecture multi-tenant

```
Enterprise
  └── Organizations (orgId)
        └── Endpoints (endpointId)
              └── Groups, Sessions, Updates, etc.
```

La plupart des endpoints utilisent `{orgId}` pour cibler une organisation.

---

## 2. Catalogue complet des endpoints (29 sections)

### 2.1 OAuth 2.0 (1 endpoint)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `POST` | `/oauth2/token` | Obtenir un token OAuth |

---

### 2.2 Search (1 endpoint)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/search/{orgId}` | Recherche rapide (rapports, endpoints, apps) |

---

### 2.3 Endpoints (8 endpoints)

Gestion des endpoints managés (serveurs, postes de travail, appareils).

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/endpoints/status/{orgId}` | Statut des endpoints |
| `GET` | `/endpoints/agent-installation/{orgId}/{installType}` | URL d'installation de l'agent |
| `GET` | `/endpoints/managed/{orgId}` | Lister tous les endpoints |
| `GET` | `/endpoints/managed/{orgId}/{endpointId}` | Détails d'un endpoint |
| `PATCH` | `/endpoints/managed/{orgId}/{endpointId}` | Modifier commentaire/nom |
| `DELETE` | `/endpoints/managed/{orgId}/{endpointId}` | Supprimer un endpoint |
| `POST` | `/endpoints/managed/{orgId}/{endpointId}/move` | Déplacer vers une autre org |
| `GET` | `/endpoints/managed/{orgId}/{endpointId}/missing-updates` | Mises à jour manquantes |

---

### 2.4 Endpoint Groups (7 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/endpoints/groups/{orgId}` | Lister les groupes |
| `POST` | `/endpoints/groups/{orgId}` | Créer un groupe |
| `GET` | `/endpoints/groups/{orgId}/{groupId}` | Détails d'un groupe |
| `PATCH` | `/endpoints/groups/{orgId}/{groupId}` | Modifier un groupe |
| `DELETE` | `/endpoints/groups/{orgId}/{groupId}` | Supprimer un groupe |
| `GET` | `/endpoints/groups/{orgId}/{groupId}/contents` | Lister les endpoints du groupe |
| `POST` | `/endpoints/groups/{orgId}/{groupId}/contents` | Ajouter/retirer des endpoints |

---

### 2.5 Remote Sessions (3 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `POST` | `/endpoints/managed/{orgId}/{endpointId}/remote-sessions` | Démarrer une session distante |
| `GET` | `/endpoints/managed/{orgId}/{endpointId}/remote-sessions/{sessionId}` | Détails d'une session |
| `PATCH` | `/endpoints/managed/{orgId}/{endpointId}/remote-sessions/{sessionId}` | Changer de moniteur |

---

### 2.6 Agent Deployment (2 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/endpoints/agent-deployment/{orgId}` | Paramètres de déploiement |
| `PATCH` | `/endpoints/agent-deployment/{orgId}` | Modifier les paramètres |

---

### 2.7 Deployers (4 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/endpoints/deployers/{orgId}` | Lister les deployers |
| `GET` | `/endpoints/deployer-installation/{orgId}/windowsEXE` | URL d'installation deployer |
| `GET` | `/endpoints/deployers/{orgId}/{deployerId}` | Détails d'un deployer |
| `DELETE` | `/endpoints/deployers/{orgId}/{deployerId}` | Supprimer un deployer |

---

### 2.8 Data Sources (5 endpoints)

Templates de scripts pour collecter des données depuis les endpoints.

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/data-sources/all` | Lister les sources de données |
| `POST` | `/data-sources/all` | Créer une source |
| `GET` | `/data-sources/all/{dataSourceId}` | Détails d'une source |
| `PATCH` | `/data-sources/all/{dataSourceId}` | Modifier une source custom |
| `DELETE` | `/data-sources/all/{dataSourceId}` | Supprimer une source custom |

---

### 2.9 Script Library (5 endpoints)

Bibliothèque de scripts PowerShell et CMD.

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/scripts/all` | Lister les scripts |
| `POST` | `/scripts/all` | Créer un script custom |
| `GET` | `/scripts/all/{scriptId}` | Détails d'un script |
| `PATCH` | `/scripts/all/{scriptId}` | Modifier un script custom |
| `DELETE` | `/scripts/all/{scriptId}` | Supprimer un script custom |

---

### 2.10 Advanced Settings (7 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/setting-templates/all` | Lister les templates de paramètres |
| `GET` | `/setting-templates/all/{templateId}` | Détails d'un template |
| `GET` | `/settings/all` | Lister les paramètres |
| `POST` | `/settings/all` | Créer un paramètre |
| `GET` | `/settings/all/{settingId}` | Détails d'un paramètre |
| `PATCH` | `/settings/all/{settingId}` | Modifier un paramètre |
| `DELETE` | `/settings/all/{settingId}` | Supprimer un paramètre |

---

### 2.11 Reports - Définitions (5 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/reports/all` | Lister les rapports |
| `GET` | `/reports/all/{reportOrCategoryId}` | Lister rapports d'une catégorie |
| `POST` | `/reports/all/custom` | Créer un rapport custom |
| `PATCH` | `/reports/all/custom/{reportId}` | Modifier un rapport custom |
| `DELETE` | `/reports/all/custom/{reportId}` | Supprimer un rapport custom |

---

### 2.12 Reports - Données (6 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/reportdata/{orgId}/{reportId}/data` | Obtenir les lignes du rapport |
| `GET` | `/reportdata/{orgId}/{reportId}/errors` | Erreurs du rapport |
| `GET` | `/reportdata/{orgId}/{reportId}/export` | Exporter un rapport |
| `POST` | `/reportdata/{orgId}/{reportId}/requery` | Re-exécuter un rapport |
| `GET` | `/reportdata/{orgId}/{reportId}/data/{reportRowId}/drilldown` | Détails d'une ligne |
| `GET` | `/reportdata/{orgId}/{reportId}/data/{reportRowId}/export` | Exporter les détails |

---

### 2.13 Software Repository (13 endpoints)

Gestion du dépôt logiciel (packages et versions).

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/software-repository/{orgId}` | Lister les packages |
| `POST` | `/software-repository/{orgId}` | Créer un package |
| `GET` | `/software-repository/{orgId}/{packageId}` | Détails d'un package |
| `PATCH` | `/software-repository/{orgId}/{packageId}` | Modifier un package |
| `DELETE` | `/software-repository/{orgId}/{packageId}` | Supprimer un package custom |
| `POST` | `/software-repository/{orgId}/{packageId}/clone` | Cloner un package |
| `POST` | `/software-repository/{orgId}/{packageId}/versions` | Créer une version |
| `GET` | `/software-repository/{orgId}/{packageId}/versions/{versionId}` | Détails d'une version |
| `PATCH` | `/software-repository/{orgId}/{packageId}/versions/{versionId}` | Modifier une version |
| `DELETE` | `/software-repository/{orgId}/{packageId}/versions/{versionId}` | Supprimer une version |
| `DELETE` | `.../{versionId}/additional-actions/{actionId}` | Supprimer une action additionnelle |
| `GET` | `/software-repository/{orgId}/match-conflicts` | Conflits de matching (nouveau) |
| `GET` | `/software-repository/{orgId}/{packageId}/match-conflicts` | Conflits de matching (existant) |

---

### 2.14 Software Repository - Upload (2 endpoints)

Upload multi-chunk resumable (jusqu'à 32 Go).

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `POST` | `.../{versionId}/upload` | Initialiser l'upload |
| `PUT` | `.../{versionId}/upload` | Envoyer les chunks |

---

### 2.15 Updates / Patches (3 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/updates/{orgId}` | Lister les mises à jour manquantes |
| `GET` | `/updates/{orgId}/{packageId}` | MAJ d'un package spécifique |
| `GET` | `/updates/{orgId}/{packageId}/versions/{versionId}/endpoints` | Endpoints manquant une MAJ |

---

### 2.16 Installed Software Inventory (5 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/installed-software/{orgId}/data` | Apps installées |
| `POST` | `/installed-software/{orgId}/requery` | Re-interroger toutes les apps |
| `GET` | `/installed-software/{orgId}/errors` | Erreurs de collecte |
| `GET` | `/installed-software/{orgId}/data/{endpointId}` | Apps d'un endpoint |
| `POST` | `/installed-software/{orgId}/requery/{endpointId}` | Re-interroger un endpoint |

---

### 2.17 Automations - Schedules (7 endpoints)

Planification d'automations (actions + scope + calendrier).

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/automations/schedules/{orgId}` | Lister les schedules |
| `POST` | `/automations/schedules/{orgId}` | Créer un schedule |
| `GET` | `/automations/schedules/{orgId}/{automationId}` | Détails d'un schedule |
| `PATCH` | `/automations/schedules/{orgId}/{automationId}` | Modifier un schedule |
| `DELETE` | `/automations/schedules/{orgId}/{automationId}` | Supprimer un schedule |
| `GET` | `.../{automationId}/deployment-statuses` | Statuts de déploiement |
| `DELETE` | `.../{automationId}/actions/{actionId}` | Supprimer une action |

---

### 2.18 Automations - Instances (6 endpoints)

Instances d'exécution d'automations.

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/automations/instances/{orgId}` | Lister les instances |
| `POST` | `/automations/instances/{orgId}` | Appliquer une automation |
| `GET` | `/automations/instances/{orgId}/{automationId}` | Détails d'une instance |
| `GET` | `.../{instanceId}/endpoint-results` | Résultats par endpoint |
| `GET` | `.../{instanceId}/endpoint-results/{endpointId}/details` | Détails d'un endpoint |
| `POST` | `.../{instanceId}/stop` | Arrêter une automation |

---

### 2.19 Automations - Action Templates (2 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/automations/action-templates` | Lister les templates d'actions |
| `GET` | `/automations/action-templates/{templateId}` | Détails d'un template |

---

### 2.20 Diagnostic Logging (1 endpoint)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/logs/{orgId}` | Obtenir les logs diagnostiques |

---

### 2.21 Subscription (3 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/subscription/enterprise` | Info licence enterprise |
| `POST` | `/subscription/enterprise/trial` | Essai gratuit |
| `POST` | `/subscription/enterprise/quote` | Demander un devis |

---

### 2.22 Usage Statistics (3 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/subscription/usage/enterprise` | Stats enterprise |
| `GET` | `/subscription/usage/organizations` | Stats par organisation |
| `GET` | `/subscription/usage/organizations/{orgId}` | Stats d'une organisation |

---

### 2.23 Users Management (8 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/me` | Utilisateur courant |
| `PATCH` | `/me` | Modifier l'utilisateur courant |
| `GET` | `/users` | Lister les utilisateurs |
| `POST` | `/users` | Créer un utilisateur |
| `GET` | `/users/{userId}` | Détails d'un utilisateur |
| `PATCH` | `/users/{userId}` | Modifier un utilisateur |
| `DELETE` | `/users/{userId}` | Supprimer un utilisateur |
| `GET` | `/users/{userId}/roles` | Rôles d'un utilisateur |

---

### 2.24 Enterprise (4 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/enterprise` | Paramètres enterprise |
| `PATCH` | `/enterprise` | Modifier les paramètres |
| `POST` | `/enterprise/request-closure` | Fermer le compte |
| `POST` | `/enterprise/revoke-closure` | Révoquer la fermeture |

---

### 2.25 Organizations (4 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/organizations` | Lister les organisations |
| `POST` | `/organizations` | Créer une organisation |
| `PATCH` | `/organizations/{orgId}` | Modifier une organisation |
| `DELETE` | `/organizations/{orgId}` | Supprimer une organisation |

---

### 2.26 Role-Based Access Control (10 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/roles` | Lister les rôles |
| `POST` | `/roles` | Créer un rôle |
| `GET` | `/roles/{roleId}` | Détails d'un rôle |
| `PATCH` | `/roles/{roleId}` | Modifier un rôle |
| `DELETE` | `/roles/{roleId}` | Supprimer un rôle |
| `POST` | `/roles/{roleId}/clone` | Cloner un rôle |
| `GET` | `/roles/{roleId}/users` | Utilisateurs d'un rôle |
| `POST` | `/roles/{roleId}/users/{userId}` | Assigner un utilisateur |
| `DELETE` | `/roles/{roleId}/users/{userId}` | Retirer un utilisateur |
| `GET` | `/permissions` | Lister les templates de permissions |

---

### 2.27 Report Subscriptions (4 endpoints)

Abonnements email aux rapports (weekly/daily patch statistics).

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/me/report-subscriptions` | Lister les abonnements |
| `POST` | `/me/report-subscriptions` | Créer un abonnement |
| `PATCH` | `/me/report-subscriptions/{subscriptionId}` | Modifier |
| `DELETE` | `/me/report-subscriptions/{subscriptionId}` | Supprimer |

---

### 2.28 Vulnerability Management (8 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/vulnerabilities/{orgId}` | Logiciels vulnérables |
| `GET` | `/vulnerabilities/{orgId}/{cveId}` | Détails d'une CVE (org) |
| `GET` | `/vulnerabilities/{orgId}/{cveId}/endpoints` | Endpoints affectés |
| `GET` | `/CVE-descriptions/{cveId}` | Détails d'une CVE (global) |
| `GET` | `/vulnerabilities/{orgId}/{cveId}/remediations` | Remédiations passées |
| `POST` | `/vulnerabilities/{orgId}/{cveId}/remediations` | Documenter un contrôle compensatoire |
| `PATCH` | `.../{cveId}/remediations/{remediationId}` | Modifier une remédiation |
| `DELETE` | `.../{cveId}/remediations/{remediationId}` | Supprimer une remédiation |

---

### 2.29 Audit Trail (3 endpoints)

| Méthode | Chemin | Description |
|---------|--------|-------------|
| `GET` | `/audit/events` | Événements d'audit |
| `GET` | `/audit/events/{id}` | Détails d'un événement |
| `GET` | `/audit/export` | Exporter le journal d'audit |

---

## 3. Résumé par domaine fonctionnel

| Domaine | Sections | Nb endpoints | Cas d'usage principaux |
|---------|----------|:------------:|------------------------|
| **Gestion des endpoints** | 2.3 - 2.7 | 24 | Inventaire machines, groupes, déploiement agents, sessions distantes |
| **Configuration** | 2.8 - 2.10 | 17 | Data sources, scripts, paramètres avancés |
| **Rapports** | 2.11 - 2.12 | 11 | Définitions, exécution, export de rapports |
| **Dépôt logiciel** | 2.13 - 2.14 | 15 | Packages, versions, upload de fichiers |
| **Déploiement logiciel** | 2.15 - 2.16 | 8 | Patches manquants, inventaire logiciel installé |
| **Automations** | 2.17 - 2.19 | 15 | Schedules, instances d'exécution, templates |
| **Sécurité / IAM** | 2.23 - 2.26 | 26 | Utilisateurs, organisations, rôles, permissions |
| **Vulnérabilités** | 2.28 | 8 | CVE, endpoints affectés, remédiations |
| **Abonnements** | 2.21 - 2.22, 2.27 | 10 | Licences, usage, abonnements rapports |
| **Audit & Logs** | 2.20, 2.29 | 4 | Logs diagnostiques, journal d'audit |
| **Divers** | 2.1, 2.2, 2.24 | 6 | Auth, recherche, enterprise |

---

## 4. Ressources externes identifiées

| Ressource | URL |
|-----------|-----|
| Documentation interactive (Swagger UI) | https://app.action1.com/apidocs/ |
| Guide API | https://www.action1.com/api-documentation/ |
| REST API Overview | https://www.action1.com/action1-rest-api/ |
| PSAction1 (module PowerShell officiel) | https://www.action1.com/psaction1/ |
| PSAction1 sur GitHub | https://github.com/Action1Corp/PSAction1 |
| Action1 MCP Server (communautaire) | https://github.com/ghively/Action1MCP |

---

## 5. Points techniques importants pour la CLI

1. **Multi-région** : La CLI devra supporter la sélection de région (NA/EU/AU) pour construire la base URL
2. **OAuth flow** : Gestion automatique du token (obtention, refresh, stockage sécurisé)
3. **orgId omniprésent** : La plupart des endpoints nécessitent un `orgId` — la CLI devra supporter un org par défaut configurable
4. **Pagination curseur** : Implémentation de l'itération automatique sur les pages
5. **Upload chunked** : Support de l'upload multi-chunk pour les fichiers volumineux (jusqu'à 32 Go)
6. **Formats de sortie** : JSON, table, CSV pour l'export des données
7. **Hiérarchie** : Enterprise > Organizations > Endpoints — les commandes doivent refléter cette structure
