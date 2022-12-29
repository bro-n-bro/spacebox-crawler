#!/bin/bash
cd $( dirname $0 )

cd ..

#echo "bash: $( find . -name '*.sh' | xargs wc -l | tail -n 1 | awk '{ print $1 }' )"
echo "golang: $( find . -name '*.go' -not -path './vendor/*' -not -path '*_models.go' -not -path 'test_*.go' | xargs wc -l | tail -n 1 | awk '{ print $1 }' )"
