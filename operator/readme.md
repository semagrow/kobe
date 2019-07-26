## Kobe Benchmark operator in Kubernetes ##
The kobe benchmark operator extends the kobe benchmarking tool so one can setup it easily 
in a cluster that runs kubernetes.
It is a kubernetes operator that allows the user to define the benchmarkexperiment by applying a set of yaml files 
that desrcibe new kubernetes custom resources which are explained below.The kobe-operator will use those resources to create 
and mantain the necessary components in kubernetes without the user having to specifically define them.

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

apiVersion: kobedataset.kobe.com/v1alpha1

kind: KobeDataset

metadata:

  name: dbpedia #the name of the
  
spec:

  image: kostbabis/virtuoso
  
  toDownload: true
  
  downloadFrom: http://users.iit.demokritos.gr/~gmouchakis/dumps/DBPedia-Subset.tar.gz
  
  count: 1
  
  port: 8890 



## KobeBenchmark ## 
A KobeBenchmark resource defines a benchmark. A benchmark consists of a set of datasets (chosen by name ) 
that must already be defined with KobeDataset resources. It also contains the definition of 1 or more sparql queries.In order for
the benchmark to be meaningfull the set of datasets should suffice for those queries.


## KobeFederator ##
A KobeFederator resource defines a federator. For semagrow the yaml is already supplied.
For other federators the following need to be provided in the yaml.

-The name of an image that deploys the federator.Also sparql port needs to be provided in the yaml .

-The name of an image that does the following job.It reads from a dataset dump and create a metadata config file .
User can set the input and output directories that his image might use f.e /sevod-scraper/input in the yaml and the operator will match those
to the real paths to the dumps and output folders.No further action needs to be taken by the user as to where those files are stored.
If the initialize happens from quering the dataset endpoint then the user should provide instead the name of
an image that does exaclty that, and the operator will provide it with the dataset sparql endpoint in an enviroment variable with the name DATASET_ENDPOINT.

-The name of an image that expects to find a set of config files and combine them to one config file .Again user can set input and output directories that
his image expects to read the small config files from and to write the combined file.

-The path that the federator needs its config files to be f.e etc/default/semagrow.The operator will make sure the above metadata file is reside at that path on the container
of the federeator
-The suffix of the federator's sparql endpoing f.e <endpoint>:<port>/SemaGrow/sparql .
The <endpoint>:<port> is provided by the operator based on the internal networking but the rest should be provided in the yaml since it is different based on the federator.

## KobeExperiment ##
A Kobe experiment resource defines the actual experiment. It consists of a federator (a KobeFederator resource) that will get benchmarked.
Also it requires the name of a benchmark that will be used (a KobeBenchmark resource).
The operator will first create a federation based on the set of datasets and the federator and will try to initialize everything
before creating a job that will run the standard kobe evaluation program (as many times as specified by the user in the yaml). If metadata files need to be created then this can take 
a lot based on datasets' sizes but operator makes sure the job will run after everything else is done.



     

