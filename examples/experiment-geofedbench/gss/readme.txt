Because there is no way to scrape the metadata for each dataset, we created experiment specification,
which uses an init container that initializes each semagrow federator with metadata created beforehand.

semagrow-templates-27 contain the templates for each experiment for the experiment with 27 datasets.
semagrow-templates-19 contain the templates for each experiment for the experiment with 19 datasets.
As in the relative benchmark specifications,
apply one semagrowtemplate specification at a time.

