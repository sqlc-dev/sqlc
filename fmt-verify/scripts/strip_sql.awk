# Blank out the contents of multi-line Go raw strings (the embedded SQL),
# leaving everything else intact. Single-line raw strings (struct tags)
# are kept as-is.
{
  if (in_raw) {
    i = index($0, "`")
    if (i > 0) { print substr($0, i); in_raw = 0 }
    next
  }
  n = gsub(/`/, "`")
  if (n == 1) { i = index($0, "`"); print substr($0, 1, i); in_raw = 1; next }
  print
}
