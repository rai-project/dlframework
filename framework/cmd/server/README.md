# MLModelScope Agent Sever Commands

**Note** that evaluations currently only run on datasets known by [DLDataset](https://github.com/rai-project/dldataset)

## Running Evaluations

One can run evaluations across different frameworks and models or on a single framework and model.

### Running Evaluations on all Frameworks / Models

[evaluate.go](https://github.com/rai-project/dlframework/blob/master/framework/cmd/server/run/evaluate.go) is a wrapper tool exists to make it easier to run evaluations across frameworks and models.
One can specify the [frameworks, models, and batch sizes](https://github.com/rai-project/dlframework/blob/master/framework/cmd/server/evaluate.go#L31-L72) to use within the file and then run evaluate.go.

- [ ]: TODO: allow one to specify the frameworks, models, and batch sizes from the command line

### Running Evaluations on a single Framework / Model

#### Example Usage

```{sh}
./tensorflow_agent dataset --debug --verbose --publish=true --fail_on_error=true --gpu=true --batch_size=320 --model_name=BVLC-Reference-CaffeNet --model_version=1.0 --database_name=tx2_carml_model_trace --database_address=minsky1-1.csl.illinois.edu --publish_predictions=false --num_file_parts=8 --trace_level=FULL_TRACE
```

### Command line options

## Available Models

```
agent info models
```

## Checking Divergence

- [ ]: TODO

To compare a single prediction's divergence you use

```
agent database divergence --database_address=$DATABASE_ADDR --database_name=carml --source=$SOURCE_ID --target=$TARGET_ID
```

## Analysing / Summarizing Results

- [ ]: TODO

```
agent evaluation --help
```
