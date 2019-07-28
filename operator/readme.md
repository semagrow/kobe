## Kobe Benchmark operator in Kubernetes ##
The kobe benchmark operator extends the kobe benchmarking tool so one can setup it easily 
in a cluster that runs kubernetes.
It is a kubernetes operator that allows the user to define the benchmark experiment by applying a set of yaml files 
that desrcibe new kubernetes custom resources .The kobe-operator will use those resources to create 
and mantain the necessary components in kubernetes without the user having to worry about them.

## Deployment of the operator in kubernetes## 
To build the operator go to operator/build and use the command 'docker build -t <operator-image-name>' . 
Push that image to a public registry.Alternative just use the already made image kostbabis/kobe-operator.

To deploy the operator first go to *operator/deploy/init/cluster* and use 
`kubectl create -f clusterrole.yaml`

`kubectl create -f cluster_role_binding.yaml`

`kubectl create -f service_account.yaml `

Then go to *operator/deploy/init/crds* and use 
`kubectl create -f kobedataset_v1alpha1_kobedataset_crd.yaml`

`kubectl create -f kobebenchmark_v1alpha1_kobebenchmark_crd.yaml `

`kubectl create -f kobefederator_v1alpha1_kobefederator_crd.yaml `

`kubectl create -f kobefederation_v1alpha1_kobefederation_crd.yaml`

`kubectl create -f kobeexperiment_v1alpha1_kobeexperiment_crd.yaml  `

`kubectl create -f kobeutil_v1alpha1_kobeutil_crd.yaml  `


Finally go to *operator/deploy/init/operator-deploy* and use 
`kubectl create -f operator.yaml`

This will set the operator running in your kubernetes cluster and needs to be done only once.

## KobeDataset ##
The KobeDataset custom resource defines a dataset that could be used in an experiment.
The operator will create and mantain a pod that runs a virtuoso instance with that dataset. It will also cache the db file and dump files for future retrieval if the pod dies and restarts or if the user deletes the kobedataset and want to redefine it. The yaml archetype is the following:

```apiVersion: kobedataset.kobe.com/v1alpha1
kind: KobeDataset
metadata:
  name: dbpedia               #the name of the dataset (must be small letters and <15 chars)
spec:
  image: kostbabis/virtuoso   
  forceLoad: true             #if false it will skip downloading and data loading and fetch the db/dump files
                              #if they exist (f.e the dataset was loaded earlier)
  downloadFrom:  http://...   #the dump location 
  count: 1                    #how many instances of this database you want in your cluster (under same service)
  port: 8890                  #the port it listens to (it defaults to 8890)
```
After writing the yaml of the dataset in the above format apply it 
`kubectl apply -f my-kobe-dataset.yaml`
Define many of these datasets depending on the experiments you want to run.One dataset can be used in many experiments and needs to only be defined once.

## KobeBenchmark ## 
A KobeBenchmark custom resource defines a benchmark in kobe. 
A benchmark consists of a set of datasets (chosen by name) that must be already  defined with the KobeDataset resources 
like its described above. It also contains the definition of one or more sparql queries.In order for the benchmark to be meaningfull the set of datasets should suffice for those queries.
The yaml is the following:
```
apiVersion: kobebenchmark.kobe.com/v1alpha1
kind: KobeBenchmark
metadata:
  name: cross-domain
spec:
  datasets:
    - name: swdfood           #these are kobedataset resources!
    - name: dbpedia
    - name: jamendo
    - name: nyt
    - name: geonames
    - name: lmdb
  queries:
    - name: q1
      language: sparql
      queryString: "SELECT ?predicate ?object WHERE { ............"
```


## KobeFederator ##
A KobeFederator resource defines a federator. For semagrow the yaml is already supplied here.
For a federator in general in order to be "kobe ready"  some things need to supplied that will be described in detail.
The yaml archetype is the following. 

```
apiVersion: kobefederator.kobe.com/v1alpha1
kind: KobeFederator
metadata:
  name: semagrow
spec:
  image: semagrow/semagrow
  port: 8080
  sparqlEnding: /SemaGrow/sparql
  imagePullPolicy: Always
  fedConfDir: /etc/default/semagrow
  
  confFromFileImage: kostbabis/semagrow-init #matadata file image from dump or endpoint
  inputDumpDir: /sevod-scraper/input
  outputDumpDir: /sevod-scraper/output
  
  confImage: kostbabis/semagrow-init-all #init image metadata from many
  inputDir: /kobe/input
  outputDir: /kobe/output
 
```

-Under _spec.image_ you must define an image that deploys your federator on a server.  

-Under _spec.port_ you must define the port your federator's endpoint listens to 

-Under _spec.sparqlEnding_ you must provide the suffixe of your federators sparql endpoint .
 For example for semagrow which listens to `<internal-endpoint>:<port>/SemaGrow/sparql` 
 then `sparqlEnding: /SemaGrow/sparql `
 
 -Under the _spec.fedConfDir_ you must specify the directory your federator expects to find its metadata files 
 in order to operate properly.For semagrow that is _/etc/default/semagrow_
 
 -Under _spec.confFromFileImage_ you must provide the name of an image that does the following.
 It reads from /kobe/input dump files of a dataset and writes at /kobe/output metadata configuration files for that dataset.
 It can also instead query the database endpoint to create the metadata file since we provide the init container 
 with thes environment variable END_POINT=".."
 The read and write directories of your image can be changed from the following 2 fields in the yaml
 _spec.inputDumpDir_ and _spec.outputDumpDir_ if its convenient.They automatically default to /kobe/input , /kobe/output

 -Under _spec.ConfImage_ you must provide the name of an image that does the following.
 It read from /kobe/input a set of different metadata files and combines them to one big configuration file of metadata for
 the experiment. For semagrow we just need to turn each dataset metadata from ttl to nt then cat them and turn them back to   .ttl. Again you can change the input and output directories your image expects to find the files and write to ,with the 
 following fields _.spec.inputDir_ and _.spec.outputDir_ .
 
 If the above are specified as described after the init process the federator will have the correct metadata file
 in the directory it expects it.

## KobeExperiment ##
A Kobe experiment resource defines the actual experiment. It consists of a federator (a KobeFederator resource) that will get benchmarked.
Also it requires the name of a benchmark that will be used (a KobeBenchmark resource).
The yaml archetype is the following:
```
apiVersion: kobeexperiment.kobe.com/v1alpha1
kind: KobeExperiment
metadata:
  name: kobeexp1
spec:
  benchmark: cross-domain  # a kobe benchmark resource
  federator: semagrow      # a kobe federator resource
  timesToRun: 1            
  dryRun: true
  forceNewInit: false 
  evalImage: kostbabis/kobe-evaluator  #the eval image for kobe-operator
 
```
-Under _spec.timesToRun_ : define the number of times you want the benchmark experiment to repeat.

-Under _spec.dryRun_ : if set to true the federation will be created and the federator initialized and the health checks 
will also happen but the experiment will hang there and no eval job will run till this flag is changed.

-Under _forceNewInit_ : if set to true it will always try to run the init image that create a metadata file from a dataset.
If set to false it will check and use preexisting metadata files if they exist for a pair of dataset -federator.
It can be used to save time since metadata extraction for big dataset take a long time and makes sense to not repeat this process.
This affects only the first init process with the image that makes a metadata file from a dataset dump or endpoint.
The second init process that combines many init files to one will always run again before init complete.
