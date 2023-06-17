#!/bin/sh

mangalcli run "$(cat ./run.lua)" --provider "$1" --vars="title=$2"