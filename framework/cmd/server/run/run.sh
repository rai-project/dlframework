#!/bin/bash
# declare -a arr=("NO_TRACE")
declare -a arr=("MODEL_TRACE" "NO_TRACE" "FRAMEWORK_TRACE" "FULL_TRACE")

for i in "${arr[@]}"
do
	echo "$i"
	go run evaluate_url.go --tracer_address=3.95.28.134:16686 --database_address=3.95.28.134:27017 --trace_level=$i
done
