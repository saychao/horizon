#!/usr/bin/env bash

set -a
horizon --config ${CONFIG:-/config.yaml} "${@:1}"
