FillOwlimDbpedia Demo

This demo fills a BigOWLim sesame store with parts of DBpedia rdf data.


PREREQUISITES:
	DBpedia rdf datasets (part) must be downloaded and! unzipped to "%baseDir%\data\rdf\dbpedia351\":
	
	Ontology Infobox Types: http://downloads.dbpedia.org/3.5.1/en/instance_types_en.nt.bz2
	Ontology Infobox Properties: http://downloads.dbpedia.org/3.5.1/en/mappingbased_properties_en.nt.bz2
	
NOTES:
    - data sets are rather huge, hence filling takes quite some time!!!
	- evaluation results can be seen in stdout and result/*.csv
	- note: for filling config property "fill" must be set to true
	- no queries are evaluated for fill mode
	- filling of local repositories must be done only once per type