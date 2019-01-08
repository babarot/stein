#!/bin/bash

go build || exit 1

export STEIN_POLICY=_examples/.policy,_examples/manifests/.policy/
./stein apply _examples/manifests/microservices/*/*/*/*

echo

export STEIN_POLICY=_examples/.policy,_examples/spinnaker/.policy/
./stein apply _examples/spinnaker/*/*/*
