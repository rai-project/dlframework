#!/bin/bash
# declare -a arr=("NO_TRACE")
declare -a arr=("MODEL_TRACE" "FRAMEWORK_TRACE" "FULL_TRACE")

go build evaluate_url.go

for i in "${arr[@]}"
do
	echo "$i"
	./evaluate_url --database_address=3.95.28.134:27017 --trace_level=$i --tracer_address=`curl http://169.254.169.254/latest/meta-data/public-ipv4`
done
