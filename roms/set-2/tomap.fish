#!/usr/bin/env fish

set files (ls *.ch8)


for f in $files
  set name (echo "$f" | sed 's@.ch8@@')
  set path $f
  echo '{' name: \"$name\", path: \"$f\" '}',
end




