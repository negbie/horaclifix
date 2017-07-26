echo "Replace IP ${@:1:1} with ${@:2:1}"
find ./ -name '*.yml' -print0 | xargs -0 sed -i "s/${@:1:1}/${@:2:1}/g"
