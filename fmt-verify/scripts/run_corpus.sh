#!/bin/bash
# Verify sqlc fmt across the endtoend corpus.
# For each config dir touching postgresql/mysql:
#   - generate in a pristine tree (control)
#   - fmt (+ capture diff) then generate in a second tree
#   - a case "needs fixing" only if fmt crashes/self-reports a bug,
#     generate breaks after fmt, or generated output changes beyond the
#     embedded SQL text.
set -u
SP=/tmp/claude-0/-home-user-sqlc/283d5543-837a-5dae-b903-ff6af9400f86/scratchpad
SQLC=$SP/sqlc
TD=/home/user/sqlc/internal/endtoend/testdata
WORK=$SP/corpus
RES=$SP/results
rm -rf "$WORK" "$RES"; mkdir -p "$WORK" "$RES"
export PATH="$SP:$PATH"
export SQLCCACHE=$SP/sqlccache
export POSTGRESQL_SERVER_URI="postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"
export MYSQL_SERVER_URI="root:mysecretpassword@tcp(127.0.0.1:3306)/mysql?multiStatements=true&parseTime=true"

echo "copying trees..."
cp -r "$TD" "$WORK/orig"
cp -r "$TD" "$WORK/fmt"

strip_file() { # $1 = file; prints normalized content
  case "$1" in
    *.go)   awk -f "$SP/strip_sql.awk" "$1" ;;
    *.json) grep -v '"text"' "$1" ;;
    *)      cat "$1" ;;
  esac
}

configs=$(cd "$TD" && find . \( -name sqlc.json -o -name sqlc.yaml -o -name sqlc.yml \) | sort)
total=0
for cfg in $configs; do
  rel=$(dirname "$cfg"); rel=${rel#./}
  grep -qE 'postgresql|mysql' "$TD/$cfg" || { echo "$rel" >> "$RES/skipped_engine.txt"; continue; }
  total=$((total+1))
  id=$(echo "$rel" | tr / _)
  A="$WORK/orig/$rel"; B="$WORK/fmt/$rel"

  (cd "$A" && timeout 30 "$SQLC" generate > "$RES/$id.gen_a.log" 2>&1)
  aok=$?
  (cd "$B" && timeout 30 "$SQLC" fmt --diff > "$RES/$id.fmtdiff.txt" 2> "$RES/$id.fmt.err")
  fexit=$?

  problems=""
  [ $fexit -ne 0 ] && problems="fmt exited $fexit"
  grep -q "this is a bug in sqlc fmt" "$RES/$id.fmt.err" && problems="${problems:+$problems; }fmt self-reported a bug"

  if [ ! -s "$RES/$id.fmtdiff.txt" ] && [ -z "$problems" ]; then
    echo "$rel" >> "$RES/unchanged.txt"
    rm -f "$RES/$id."*; continue
  fi

  (cd "$B" && timeout 30 "$SQLC" fmt >/dev/null 2>&1)
  (cd "$B" && timeout 30 "$SQLC" generate > "$RES/$id.gen_b.log" 2>&1)
  bok=$?

  if [ $aok -ne 0 ]; then
    if [ -f "$A/stderr.txt" ] || [ -d "$A/stderr" ]; then
      echo "$rel" >> "$RES/expected_failure.txt"
    else
      echo "$rel: $(head -1 "$RES/$id.gen_a.log")" >> "$RES/unverifiable.txt"
    fi
    [ -n "$problems" ] && { echo "### $rel"; echo "reason: $problems"; cat "$RES/$id.fmt.err"; echo; } >> "$RES/needs_fix.txt"
    rm -f "$RES/$id."*; continue
  fi

  if [ $bok -ne 0 ]; then
    problems="${problems:+$problems; }generate failed after fmt"
  else
    structdiff=""
    while IFS= read -r f; do
      frel=${f#"$A"/}
      case "$frel" in *.sql|sqlc.json|sqlc.yaml|sqlc.yml|exec.json) continue ;; esac
      d=$(diff -u <(strip_file "$A/$frel") <(strip_file "$B/$frel") 2>&1) || \
        structdiff="$structdiff--- $frel
$d
"
    done < <(find "$A" -type f)
    [ -n "$structdiff" ] && problems="${problems:+$problems; }generated output changed beyond SQL text"
  fi

  if [ -n "$problems" ]; then
    { echo "### $rel"
      echo "reason: $problems"
      echo "--- query diff:"; cat "$RES/$id.fmtdiff.txt"
      [ -s "$RES/$id.fmt.err" ] && { echo "--- fmt stderr:"; cat "$RES/$id.fmt.err"; }
      [ $bok -ne 0 ] && { echo "--- generate log:"; head -20 "$RES/$id.gen_b.log"; }
      [ -n "${structdiff:-}" ] && { echo "--- structural diff:"; echo "$structdiff"; }
      echo
    } >> "$RES/needs_fix.txt"
    echo "$rel" >> "$RES/needs_fix_list.txt"
  else
    echo "$rel" >> "$RES/benign.txt"
  fi
  rm -f "$RES/$id."*
done
echo "done: $total configs processed" > "$RES/DONE"
