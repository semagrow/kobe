if [ "$1" == "semagrow" ]
then
	cat suites/$1-$2/$3/result/result.csv
	grep tion\ time run.txt | awk '{ print $8 }'
fi
