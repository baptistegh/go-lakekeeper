# Examples

This example follows the same principles as those used in the official Lakekeeper examples.

In this case, only the [`access-control-simple`](https://github.com/lakekeeper/lakekeeper/blob/main/examples/access-control-simple) component is used.

The objective is to configure the various users and the warehouse using the Lakekeeper CLI.

## Prerequisites

Before proceeding, you must start the Lakekeeper example environment. This setup does **not** require running the notebooks at this stage.

To do so, execute the following commands:

```sh
git clone --depth 1 --branch main https://github.com/lakekeeper/lakekeeper.git
cd lakekeeper/examples/access-control-simple
docker compose up -d
cd ../../../
```

## Initializing the Environment

Once the Lakekeeper example environment is running, you can initialize the configuration using the `init.sh` script. This script details the steps executed by the CLI when the container starts.

To launch it, return to the original folder `(go-lakekeeper/examples)` and run:

```sh
docker compose up
```

If you prefer to build the image manually for testing purposes, you can use:

```sh
docker compose -f docker-compose.yaml -f docker-compose-build.yaml up --build
```

## Available Notebooks

Once your environment is correctly set up, you can use the notebooks in the [Jupyter UI](http://localhost:8888/):

* `03-01-Spark.ipynb`
* `03-02-Trino.ipynb`
* `03-03-Starrocks.ipynb`
* `03-04-PyIceberg.ipynb`
* `03-02-DuckDB.ipynb`
