#!/bin/bash
# declare -a arr=("NO_TRACE")
declare -a arr=("MODEL_TRACE" "FRAMEWORK_TRACE" "FULL_TRACE")

go build evaluate_url.go

for i in "${arr[@]}"
do
	echo "$i"
	./evaluate_url --tracer_address=3.95.28.134:16686 --database_address=3.95.28.134:27017 --trace_level=$i
done
