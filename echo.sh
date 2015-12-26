#!/bin/bash

>&2 echo "stderr test"
echo "repo is $1"
echo "project is $2"
echo "branch is $3"
echo "type is $4"
echo "ref is $5"
echo "sleeping 2 seconds"
sleep 2
echo "exiting"
exit 1