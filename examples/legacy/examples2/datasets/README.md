# README

A simple script for downloading dataset dumps locally.

### Usage

```
./scripts/staging.sh [experiment_name]
```
To add a custom experiment add the new datasets in the datasets.csv and the experiment name with the corresponding dataset name in the experiments_datasets.csv file.  

* datasets.csv schema:
     ```
    [dataset_name],[url/to/download],[path/to/unzip]
    ```
* experiments_datasets.csv schema:
    ```
    [experiment_name],[dataset_name]
    ```

### Requirements

The scriipt requires the [csvkit](https://csvkit.readthedocs.io/en/1.0.1/) tool.
