#!/bin/bash

set -e

COMMAND=$1
ENV_FILE=$2

eval "$COMMAND"

env >>"$ENV_FILE"
