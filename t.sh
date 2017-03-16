# maybe more powerful
# for mac (sed for linux is different)
dir=`echo ${PWD##*/}`
grep "robot-train" * -R | grep -v Godeps | awk -F: '{print $1}' | sort | uniq | xargs sed -i '' "s#robot-train#$dir#g"
mv robot-train.ini $dir.ini

