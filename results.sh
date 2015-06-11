if [ "$1" == "semagrow" ]
then
	cat suites/semagrow-reactive/$2/result/result.csv
	grep tion\ time run.txt | awk '{ print $8 }'
fi
