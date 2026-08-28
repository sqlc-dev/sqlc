#!/bin/bash
# Re-verify every config with a fmt diff using the AST-based Go stripper
# and the marino semantic check for MySQL. Produces final classification.
set -u
SP=/tmp/claude-0/-home-user-sqlc/283d5543-837a-5dae-b903-ff6af9400f86/scratchpad
W=$SP/corpus
RES=$SP/results
OUT=$SP/final
rm -rf "$OUT"; mkdir -p "$OUT"

strip_file() {
  case "$1" in
    *.go)   "$SP/stripgo" "$1" ;;
    *.json) grep -v '"text"' "$1" ;;
    *)      cat "$1" ;;
  esac
}

# All configs that had any fmt diff: benign + needs_fix + (unverifiable/expected with sql diffs)
all=$(cat "$RES/benign.txt" "$RES/needs_fix_list.txt" 2>/dev/null; \
      cut -d: -f1 "$RES/unverifiable.txt" 2>/dev/null; cat "$RES/expected_failure.txt" 2>/dev/null)

for rel in $all; do
  A="$W/orig/$rel"; B="$W/fmt/$rel"
  [ -d "$A" ] || continue
  # which .sql files changed
  changed=$(cd "$A" && find . -maxdepth 1 -name '*.sql' | while read -r f; do
    cmp -s "$A/$f" "$B/$f" || echo "${f#./}"
  done)
  [ -z "$changed" ] && { echo "$rel" >> "$OUT/nodiff.txt"; continue; }

  problems=""
  detail=""

  # semantic check for mysql configs
  cfgfile=$(ls "$A"/sqlc.json "$A"/sqlc.yaml "$A"/sqlc.yml 2>/dev/null | head -1)
  if grep -q mysql "$cfgfile"; then
    for f in $changed; do
      fp=$("$SP/mysqlfp" "$A/$f" "$B/$f" 2>&1); rc=$?
      if [ $rc -eq 1 ]; then
        problems="${problems:+$problems; }semantic drift in $f"
        detail="$detail--- semantic check ($f):
$fp
"
      elif [ $rc -eq 2 ]; then
        echo "$rel: $f: $(echo "$fp" | head -1)" >> "$OUT/fp_unparseable.txt"
      fi
    done
  fi

  # structural codegen check (only where the pristine generate succeeded)
  if grep -qxF "$rel" "$RES/benign.txt" "$RES/needs_fix_list.txt" 2>/dev/null; then
    structdiff=""
    while IFS= read -r f; do
      frel=${f#"$A"/}
      case "$frel" in *.sql|sqlc.json|sqlc.yaml|sqlc.yml|exec.json) continue ;; esac
      d=$(diff -u <(strip_file "$A/$frel") <(strip_file "$B/$frel") 2>&1) || \
        structdiff="$structdiff--- $frel
$d
"
    done < <(find "$A" -type f)
    if [ -n "$structdiff" ]; then
      problems="${problems:+$problems; }generated code changed beyond SQL text"
      detail="$detail--- structural diff:
$structdiff
"
    fi
  fi

  if [ -n "$problems" ]; then
    { echo "### $rel"
      echo "reason: $problems"
      echo "--- query diff:"
      for f in $changed; do diff -u "$A/$f" "$B/$f" | tail -n +3; done
      echo "$detail"
      echo
    } >> "$OUT/needs_fix.md"
    echo "$rel" >> "$OUT/needs_fix_list.txt"
  else
    echo "$rel" >> "$OUT/clean.txt"
  fi
done
echo done > "$OUT/DONE"
