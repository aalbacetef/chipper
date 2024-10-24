#!/usr/bin/env fish


set DIR (dirname (status --current-filename))
set files (ls $DIR/**/*.ch8)

set json_fpath $DIR/manifest.json 


echo '{' > "$json_fpath"

set n (count $files)
set want (math "$n - 1")

function is_last
  set -l current (grep ':' $json_fpath | wc -l)

  if test "$current" -eq "$want"
    return 0
  end 

  return 1
end 

for f in $files
  set -l name (echo (basename $f) | sed 's@.ch8@@')
  set -l fpath (echo $f | sed 's@./@@')
  
  set -l trailing ","
  
  if  is_last
    set trailing ""
  end 
  

  echo "  \"$name\":  \"$fpath\"$trailing " >> $json_fpath
end

echo -n "}" >> "$json_fpath" 



