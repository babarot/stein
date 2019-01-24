#!/bin/bash

go build || exit 1

# stein loads the HCL files located on .policy directory by default
# in addition, .policy directory can be overridden by each directory of given arguments
#
# in this case,
#   stein applies rules located in these default directory to _examples/manifests/microservices/*/*/*/*
#   * _examples/.policy/
#   * _examples/manifests/.policy/
#   stein doesn't apply this rules to them
#   * _examples/spinnaker/.policy/
#
# Regardless of the default directory placed under the given path,
# the following environment variables can be specified for the policy applied to all paths.
# this variables can take multiple values separated by a comma, also can take directories and files
#
# export STEIN_POLICY=root-policy/,another-policy/special.hcl

./stein apply \
    _examples/manifests/microservices/*/*/*/* \
    _examples/spinnaker/*/*/*
