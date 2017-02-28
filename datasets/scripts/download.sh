
URL=${1%,*}
DIR=${1#*,}
EXT=${URL##*.}

cd $DIR
wget $URL

case $EXT in
  tar)
    tar xvf *.tar
  ;;
  7z)
    7z x *.7z
  ;;
  zip)
    unzip .zip
  ;;
esac

    
