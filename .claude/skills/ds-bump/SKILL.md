---
name: ds-bump
description: >-
  Propage une version publiée (tag) de go-kiban-design-system vers ses repos
  consommateurs (rekon, crm, link, klin ; plus tard workfloo, kiban-cloud) :
  go get @tag + go mod tidy + go build, sans committer. À utiliser quand
  l'utilisateur veut diffuser/propager une nouvelle version du design system
  après avoir publié un tag — déclencheurs : "bump le design system en vX",
  "actualizar design system en vX.Y.Z", "mettre à jour go-kiban-design-system
  dans les repos", "propager la version vX aux consommateurs".
---

# ds-bump — propager une version du design system

La logique vit dans le script versionné **`scripts/ds-bump.sh`** (à la racine de ce repo). Tu ne la réimplémentes pas : tu exécutes le script et tu interprètes sa sortie.

## Procédure

1. **Récupère le tag** depuis la demande de l'utilisateur (ex. `v0.0.25`), et la liste optionnelle de repos à cibler (noms courts : `rekon`, `crm`, `link`, `klin`). Si aucun tag n'est exprimé, demande-le et arrête-toi.

2. **Lance le script** depuis la racine du design system :

   ```bash
   ./scripts/ds-bump.sh <tag> [repo ...]
   ```

   - Sans repo → tous les consommateurs par défaut. Avec → seulement ceux listés (`./scripts/ds-bump.sh v0.0.25 crm link`).
   - Les repos cibles sont des dossiers frères, hors du workspace, et le script fait du réseau (`go get`) + des écritures cross-repo. **Exécute-le avec le sandbox désactivé** (`dangerouslyDisableSandbox: true`), sinon le process est tué.

3. **Interprète le tableau « RÉSUMÉ »** affiché en fin et restitue-le à l'utilisateur :
   - Le script s'arrête de lui-même si le tag n'existe pas sur le remote → dis à l'utilisateur de `git push origin <tag>` d'abord.
   - Repo `SKIP (replace actif)` → l'utilisateur doit commenter son `replace` local avant de re-bumper.
   - `build FAIL` → remonte l'erreur exacte (imprimée juste au-dessus du résumé).

4. **Ne committe ni ne push rien.** Le script laisse les `go.mod`/`go.sum` modifiés ; l'utilisateur review et committe lui-même.

## Notes

- Un `build FAIL` dû à un changement d'API du design system est un signal réel : rapporte-le, ne « répare » pas le code du consommateur sans demande explicite. Un échec dans un paquet sans rapport (ex. un `package main` utilitaire sans `func main`) est pré-existant — signale-le comme tel.
- **Piège shell** : dans ce harness, `grep`/`sed`/`awk` peuvent être des fonctions injectées qui `exec` et tuent un script multi-commandes. `ds-bump.sh` les évite (détection via `go list -m`). Si tu improvises du shell autour, préfixe par `command grep` et capture les sorties dans des variables.
- **Ajouter un consommateur** (workfloo, kiban-cloud, …) : décommente la ligne dans le tableau `REPOS` de `scripts/ds-bump.sh`. Le nom court devient utilisable comme argument de ciblage.
