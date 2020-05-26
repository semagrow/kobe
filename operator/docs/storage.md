# Storage

The configuration and deployment of a federation of massive datasets from
scratch may take a considerable amount of time. There are, however, several
chunks of data that are produced and can be cached during an experiment to
drastically speed up subsequent experiments. Kobe needs persistent storage
volumes to cache and reuse data.

The data that are subject to storing beyond the lifecycle of an experiment 
are:
* The downloaded files that comprise datasets along with the checksums. It is
  frequent that the same datasets can be used in different benchmarks. By
  caching locally the downloaded files alleviates the stress on the network
  bandwidth.
* Each dataset can be loaded and served by various database systems. Importing
  massive datasets into a database is a time-consuming process. It is reasonable
  to backup and restore the already-imported database files when the same system
  for the same dataset. Restoring the database files is a much faster process
  than reloading it from the downloaded dataset files. It should be noted that
  the backed up files are determined by the version of tha database system.
* The last type of files that can be cached are the files used by the federator.
  Some federators may used metadata for each dataset they federate. This process
  can also be time-consuming so backing those files may also speed up the
  experiments initialization phase.