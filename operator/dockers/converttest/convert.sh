for file in *.rdf
do
  java -jar rdf2rdf-1.0.1-2.3.1.jar ${file} /output/${file}.nt
done
