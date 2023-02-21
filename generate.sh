#!/usr/bin/env bash

GENERATOR_IMAGE=registry.gitlab.com/tokend/openapi-go-generator:acee5bc606997864a49bc213e8dd164bfb077d18 # latest generator commit in master

GENERATED="${GOPATH}/src/gitlab.com/tokend/regources/generated"
OPENAPI_DIR="${GOPATH}/src/github.com/SafeRE-IT/horizon/docs/build"
PACKAGE_NAME=regources

function printHelp {
    echo "usage: ./generate.sh [<flags>] [<args> ...]
            script to generate regources for horizon

            Flags:
                  --package PACKAGE        package name of generated stuff (first line of file)
                  --image IMAGE            name of generator docker image (default is openapi-generator)

              -h, --help                   Show this help.
              -v, --to-vendor              put generated output to vendor dir
              -p, --path-to-generate PATH  path to put generated things
              -i, --input OPENAPI_DIR      path to dir where openapi.yaml is stored (default horizon/docs/build for horizon)"
}

function parseArgs {
    while [[ -n "$1" ]]
    do
        case "$1" in
            -h | --help)
                printHelp && exit 0
                ;;
            -v | --to-vendor)
                GENERATED="${GOPATH}/src/github.com/SafeRE-IT/horizon/vendor/gitlab.com/tokend/regources/generated"
                ;;
            -p | --path-to-generate) shift
                [[ ! -d $1 ]] && echo "path $1 does not exist or not a dir" && exit 1
                GENERATED=$1
                ;;
            --package) shift
                [[ ! -z "$1" ]] && echo "package name not specified" && exit 1
                PACKAGE_NAME=$1
                ;;
            -i | --input) shift
                [[ ! -f "$1/openapi.yaml" ]] && echo "file openapi.yaml does not exist in $1 or not a file" && exit 1
                OPENAPI_DIR=$1
                ;;
            --image) shift
                [[ "$(docker images -q $1)" == "" ]] && echo "image $1 does not exist locally" && exit 1
                GENERATOR_IMAGE=$1
                ;;
        esac
        shift
    done
}

function generate {
    docker run -v ${OPENAPI_DIR}:/openapi -v ${GENERATED}:/generated ${GENERATOR_IMAGE} generate --generate-horizon-stuff --meta-for-lists
}

parseArgs $@
generate
