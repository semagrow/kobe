
cat $1  | grep voidDescription | 
     sed 's|.*voidDescription <\(.*\)>.*|\1|' | xargs cat | 
     sed -f ~/datasets/Bench1.mappings.sed > concat.out

cat concat.out | grep '^@prefix' | sort | uniq > prefixes.out
cat concat.out | grep -v '^@prefix' > void.out
cat prefixes.out void.out > void.n3


rm -f concat.out prefixes.out void.out
