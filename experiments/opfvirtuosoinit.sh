DATASETS="
https://data.openphacts.org/free/2.0/rdf/aers.tar
https://data.openphacts.org/free/2.0/rdf/bao.tar
https://data.openphacts.org/free/2.0/rdf/caloha.tar
https://data.openphacts.org/free/2.0/rdf/chebi.tar
https://data.openphacts.org/free/2.0/rdf/chembl.tar
https://data.openphacts.org/free/2.0/rdf/conceptwiki.tar
https://data.openphacts.org/free/2.0/rdf/disgenet.tar
https://data.openphacts.org/free/2.0/rdf/doid.tar
https://data.openphacts.org/free/2.0/rdf/drugbank.tar
https://data.openphacts.org/free/2.0/rdf/enzyme.tar
https://data.openphacts.org/free/2.0/rdf/go.tar
https://data.openphacts.org/free/2.0/rdf/goa.tar
https://data.openphacts.org/free/2.0/rdf/ncats.tar
https://data.openphacts.org/free/2.0/rdf/nextprot.tar
https://data.openphacts.org/free/2.0/rdf/ocrs.tar
https://data.openphacts.org/free/2.0/rdf/uniprot.tar
https://data.openphacts.org/free/2.0/rdf/wikipathways.tar
"

for DATASET in $DATASETS
do
  wget $DATASET
  TAR=`echo $DATASET | colrm 1 41`
  tar xvf $TAR
done

wget https://data.openphacts.org/free/2.0/ims/ims-linksets-2.0.tar.gz
tar xzvf ims-linksets-2.0.tar.gz
echo "http://ims.openphacts.org/" > linksets/default.graph
