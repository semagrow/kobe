# kobectl

kobectl controls the KOBE open benchmarking engine.

## Installation

kobectl is a `sh` script and can be found in the
[bin](https://github.com/semagrow/kobe/tree/devel/bin). Essentially, kobectl is
a wrapper of Kubernetes commands so [kubectl] must be installed and available.

You can make kobectl accessible to your path by
```sh
export PATH="$(pwd)/bin:$PATH"
```

## Commands

| Command    | Explanation                                       |
|:-----------|:--------------------------------------------------|
|  `apply`   | apply a resource using a .yaml configuration file |
|  `get`     | display all resources of specific type            |
|  `show`    | show the state of a benchmark or an experiment    |
|  `delete`  | delete a resource of specific type                |
|  `install` | install KOBE components                           |
|  `purge`   | uninstall KOBE                                    |
|  `help`    | print a help message                              |


## Usage

*  `kobectl apply [configuration_file]`
*  `kobectl get [resource_type]`
*  `kobectl show [resource_type] [resource]`
*  `kobectl delete [resource_type] [resource]`
*  `kobectl install [component] [kobe-directory]`
*  `kobectl purge [kobe-directory]`

`[resource_type]` can be any of:
  `benchmark(s)`,
  `experiment(s)`,
  `federatortemplate(s)`,
  `datasettemplate(s)`.

`[component]` can be any of:
  `operator`, `operator-v1`, `operator-v1beta1`, `istio`, `efk`, `full`

## Other

For more advanced control options for KOBE, use [kubectl].

[kubectl]: https://kubernetes.io/docs/reference/kubectl/overview/