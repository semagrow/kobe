## Kobe Benchmark operator in Kubernetes ##
The kobe benchmark operator extends the kobe benchmarking tool so one can setup it easily 
in a cluster that runs kubernetes.
It is a kubernetes operator that allows the user to define the benchmark experiment by applying a set of yaml files 
that desrcibe new kubernetes custom resources .The kobe-operator will use those resources to create 
and mantain the necessary components in kubernetes without the user having to worry about them.

## Deployment of the operator in kubernetes ## 
First clone this project and get in the kobe/operator directory and checkout the feat-k8s-operator branch. 

`git clone https://github.com/kostbabis/kobe` 

`git checkout feat-k8s-operator`

To build the operator go to operator/build and use the command `docker build -t <operator-image-name> . ` . 
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
You will get a confirmation message that each resource has successfully been created.

This will set the operator running in your kubernetes cluster and needs to be done only once.

**You might also need to install nfs-common to every node in your cluster if it doesn't already exists else the mounts to the nfs server used for caching will not work. For example in ubuntu use apt install nfs-common**

The general procedure of running an experiment is this. First you create a set of datasets by defining new **KobeDatasets** 
recources.Then you define one or more **KobeBenchmark** resources and one or more **KobeFederators** .At last you define a **KobeExperiment**. Everything is explained below.

## KobeDataset ##
The KobeDataset custom resource defines a dataset that could be used in an experiment.
The operator will create and mantain a pod that runs a virtuoso instance with that dataset. It will also cache the db file and dump files for future retrieval if the pod dies and restarts or if the user deletes the kobedataset and want to redefine it . The yaml archetype is the following:

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
`kubectl apply -f my-kobe-dataset.yaml`.
You will get a confirmation message. You can also check the progress of the dataset creation by using
`kubectl get pods ` and `kubectl logs <kobedataset-podname> `. 
Define many of these datasets depending on the experiments you want to run. One dataset can be used in many experiments and needs to only be defined once.
You can also find already made dataset definition under _operator/deploy/yamls/datasets/_  for a few datasets including a subset of dbpedia.


## KobeBenchmark ##
A KobeBenchmark custom resource defines a benchmark in kobe. 
A benchmark consists of a set of datasets that must be already  defined with the KobeDataset resources 
like its described above. It also contains the definition of one or more sparql queries that are gonna get tested against the datasets in the benchmark.
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
- Under _spec.datasets.name[*]_ you must write down the name of the datasets your benchmark will include. The names must
be the same as the metadata.name of the KobeDataset custom resources defined above.
- Under _spec.queries[*]_ you must write down the queries of your benchmark. Query name is the name of the query .Language 
for now should always be sparql and queryString should be the string that contains your query.
Apply the yaml again by issuing `kubectl apply -f my-kobe-benchmark.yaml`.
You will get a message that the resource has been created

## KobeFederator ##
A KobeFederator resource defines a federator. For semagrow the yaml is already supplied under _operator/deploy/yamls/federators_.
For a federator in general in order to be able to get benchmarked with kobe some things need to be supplied.Those will be described in detail in a bit.
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
Specifically define these:
- Under _spec.image_: here you must define an image that deploys your federator. For example in the above yaml the semagrow/semagrow image deploys semagrow on a tomcat server

- Under _spec.port_: here you must define the port that your federator's endpoint listens to.

- Under _spec.sparqlEnding_: you must provide the suffixe of your federators sparql endpoint .
 For example for semagrow which listens to `<internal-endpoint>:<port>/SemaGrow/sparql` then `sparqlEnding: /SemaGrow/sparql ` .The `<internal-endpoint>:<port>` will be provided by the operator to where its needed and you dont have to set this    anywhere.
 
 - Under the _spec.fedConfDir_ you must specify the directory your federator expects to find its metadata files 
 in order to operate properly. For semagrow that is _/etc/default/semagrow_ .
 
 - Under _spec.confFromFileImage_ you must provide the name of an image that does the following.
 It creates a container thatreads from _/kobe/input_ dump files of a dataset and writes at _/kobe/output_ metadata configuration files for that dataset.
 It can also instead query directly the database sparql endpoint to create the metadata file since we provide the init container with an environment variable called `END_POINT` which contains the full url of the sparql endpoint of the dataset
 
 The image should be oblivious of what dataset it makes the metadata for and incorporate only the necessary logic to make   that file. For example with semagrow we provide an image that uses the sevod-scraper (check it under semagrow in github) 
 to process the dump files of a dataset (f.e dbpedia) and return a dbpedia.ttl file for this specific set.
 The read and write directories of your image can be changed from the following 2 fields in the yaml
 _spec.inputDumpDir_ and _spec.outputDumpDir_ if its convenient.
 They automatically default to _/kobe/input_ , _/kobe/output_

 -Under _spec.ConfImage_ you must provide the name of an image that does the following.
 It reads from _/kobe/input_ a set of different metadata files and combines them to one big configuration file of metadata for the benchmark.
Your image should not care about what datasets the files belong to and only do the union of them .
For example, with semagrow we just need to turn each dataset metadata from .ttl to .nt then concatenate them and turn them back to .ttl. 
Again if you want to change the input and output directories your image expects to find the files and write to , you can with the following fields _.spec.inputDir_ and _.spec.outputDir_ .
 
 If the above are specified as described ,after the init process the federator will have the correct metadata file
 in the directory it expects to.
 
 Apply the federator yaml again with
 `kubectl apply -f my-kobe-federator.yaml`
 A confirmation message that it has been created should appear.

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
- Under _spec.benchark_: here you place the benchmark name. It must be the same as the name of a KobeBenchmark resource
you defined earlier.

- Under _spec.federator_: here you place the federator name. It must be the same as the name of a KobeFederation resource
you defined earlier.

- Under _spec.timesToRun_ : define the number of times you want the benchmark experiment to repeat.

- Under _spec.dryRun_ : if set to true the federation will be created and the federator initialized . The health checks 
will also happen but the experiment will hang there and no evaluation job will run till this flag is changed.

- Under _forceNewInit_ : if set to true it will always try to run the init image that create a metadata file from a dataset for this federator.
If set to false it will check and use preexisting metadata files if they exist for a pair of dataset - federator.
It can be used to save time since metadata extraction for big dataset take a long time and makes sense to not repeat this process.
This affects only the first init process with the image that makes a metadata file from a dataset dump or endpoint.
The second init process that combines many init files to one will always run again before init complete.

After your define the experiment apply it again with
`kubectl apply -f my-kobe-experimen.yaml ` 
To see the progress you can use `kubectl get pods` .The federation (that is the federator initialized with a set of datasets) will be the pod with a name same as the KobeExperiment.metadata.name .You can see there the stage of the init containers that run the init process .
You can also use `kubectl logs <federation-pod> -c initcontainer{0..x} ` to check the process of each one of the init containers as well.
Keep in mind that if forceNewInit is false only init containers that correspond to federator-dataset pairs that haven't initialized in the past will spawn.
To see the logs of the second init image of your federator use the 
`kubectl logs <federation-pod> -c initfinal `

After the init process is done a set of jobs will spawn sequentially based on timesToRun number.Those jobs  will run the evaluation program. The previous job needs to end before the next will start and the experiments will not run in parallel..
Currently to get the result of your benchmark you have to see the logs of these jobs using
`kubectl logs <federation-<job_number>-<job hash> > `. This will print the result to your screen.
To find the jobs pod name you can use `kubectl get pods`
