if [ -z $1 ]
then
  echo "usage: $0 [experiment]"
  exit
fi

EXPERIMENTS_DATASETS=`dirname $0`"/../resources/experiments_datasets.csv"
DATASETS=`dirname $0`"/../resources/datasets.csv"
DOWNLOAD=`dirname $0`"/download.sh"

csvsql --query "SELECT url, directory 
                FROM experiments_datasets, datasets 
                WHERE experiments_datasets.dataset=datasets.dataset
                AND experiments_datasets.experiment=\"$1\"" \
      ../resources/experiments_datasets.csv ../resources/datasets.csv \
 | tail -n +2 | xargs $DOWNLOAD

