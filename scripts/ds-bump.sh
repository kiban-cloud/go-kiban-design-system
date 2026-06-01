#!/usr/bin/env bash
#
# ds-bump.sh — propage une version publiée de go-kiban-design-system vers les
# repos consommateurs (go get + go mod tidy + go build).
#
# Usage:
#   scripts/ds-bump.sh <tag> [repo ...]
#
#   <tag>     ex. v0.0.25  (obligatoire ; doit déjà être poussé sur le remote)
#   [repo]    noms courts à cibler (rekon crm link klin). Si absent : tous.
#
# Ne committe ni ne push rien : laisse go.mod/go.sum modifiés dans chaque repo.
#
# Notes d'implémentation :
#   - Détection version + replace via `go list -m` (jamais `grep` : dans certains
#     harness `grep` est une fonction qui exec et tue le script).
#   - Un build cassé dans un repo n'empêche pas de traiter les suivants.

set -uo pipefail

readonly MOD="github.com/kiban-cloud/go-kiban-design-system"

# Racine du design system = parent du dossier scripts/.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly DS_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
readonly SIBLINGS="$(cd "$DS_ROOT/.." && pwd)"

# Repos consommateurs : "nom_court=chemin" (chemins frères du design system).
# Pour activer workfloo / kiban-cloud plus tard, décommente les lignes.
REPOS=(
  "rekon=$SIBLINGS/rekon-backend"
  "crm=$SIBLINGS/crm"
  "link=$SIBLINGS/link"
  "klin=$SIBLINGS/klin-backend"
  # "workfloo=$SIBLINGS/workfloo-backend"
  # "kiban-cloud=$SIBLINGS/kiban-cloud-backend"
)

die() { echo "erreur: $*" >&2; exit 1; }

path_for() {
  local want="$1" entry
  for entry in "${REPOS[@]}"; do
    [ "${entry%%=*}" = "$want" ] && { echo "${entry#*=}"; return 0; }
  done
  return 1
}

ds_version() {
  go list -m -f '{{.Version}}' "$MOD" 2>/dev/null
}

ds_replace_path() {
  # Vide si pas de replace ; sinon le chemin/cible du replace.
  go list -m -f '{{if .Replace}}{{.Replace.Path}}{{end}}' "$MOD" 2>/dev/null
}

# --- args ---------------------------------------------------------------
TAG="${1:-}"
[ -n "$TAG" ] || die "tag manquant. Usage: scripts/ds-bump.sh <tag> [repo ...]"
shift || true

# Liste cible : noms passés en argument, ou tous.
TARGETS=()
if [ "$#" -gt 0 ]; then
  for name in "$@"; do
    p="$(path_for "$name")" || die "repo inconnu: '$name' (connus: rekon crm link klin)"
    TARGETS+=("$name=$p")
  done
else
  TARGETS=("${REPOS[@]}")
fi

# --- 1. le tag existe-t-il sur le remote ? ------------------------------
echo ">> Vérification du tag $TAG sur le remote…"
remote_hit="$(git -C "$DS_ROOT" ls-remote --tags origin "refs/tags/$TAG")"
[ -n "$remote_hit" ] || die "le tag $TAG n'existe pas sur le remote. Pousse-le d'abord: git push origin $TAG"

# --- 2. boucle sur les repos --------------------------------------------
declare -a SUMMARY=()
for entry in "${TARGETS[@]}"; do
  name="${entry%%=*}"; path="${entry#*=}"
  echo
  echo "==================== $name ($path) ===================="

  if [ ! -d "$path" ]; then
    echo "  dossier introuvable -> SKIP"
    SUMMARY+=("$name|—|—|SKIP (dossier introuvable)")
    continue
  fi
  cd "$path" || { SUMMARY+=("$name|—|—|SKIP (cd impossible)"); continue; }

  repl="$(ds_replace_path)"
  if [ -n "$repl" ]; then
    echo "  REPLACE ACTIF ($repl) -> SKIP (commente le replace d'abord)"
    SUMMARY+=("$name|—|—|SKIP (replace actif)")
    continue
  fi

  before="$(ds_version)"
  if ! go get "$MOD@$TAG" 2>&1 | sed 's/^/    [get] /'; then
    echo "  go get a échoué"; SUMMARY+=("$name|$before|?|go get FAIL"); continue
  fi
  if ! go mod tidy 2>&1 | sed 's/^/    [tidy] /'; then
    echo "  go mod tidy a échoué (on continue pour voir le build)"
  fi
  after="$(ds_version)"
  echo "  version: $before -> $after"

  if go build ./... 2>/tmp/ds-bump-build.log; then
    echo "  build: OK"
    SUMMARY+=("$name|$before|$after|build OK")
  else
    echo "  build: FAIL"
    sed 's/^/    /' /tmp/ds-bump-build.log
    SUMMARY+=("$name|$before|$after|build FAIL")
  fi
done

# --- 3. résumé ----------------------------------------------------------
echo
echo "======================= RÉSUMÉ ======================="
printf '%-12s %-10s %-10s %s\n' "repo" "avant" "après" "statut"
for row in "${SUMMARY[@]}"; do
  IFS='|' read -r r b a s <<< "$row"
  printf '%-12s %-10s %-10s %s\n' "$r" "$b" "$a" "$s"
done
echo
echo "Aucun commit effectué : review les go.mod/go.sum puis committe toi-même."
