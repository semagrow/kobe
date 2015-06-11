FillNativeNytimes Demo

This demo fills a native sesame store with nytimes rdf data.


PREREQUISITES:
	Nytimes rdf datasets must be downloaded to "%baseDir%\data\rdf\nytimes\":
	
	Locations: http://data.nytimes.com/locations.rdf
	Organizations: http://data.nytimes.com/organizations.rdf
	People: http://data.nytimes.com/people.rdf
	
NOTES:
    - data sets are rather small, hence filling goes fast
	- evaluation results can be seen in stdout and result/*.csv
	- note: for filling config property "fill" must be set to true
	- no queries are evaluated for fill mode
	- filling of local repositories must be done only once per type