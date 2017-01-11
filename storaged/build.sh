#!/bin/bash

go build -ldflags '-linkmode external -extldflags -static' || exit $?
docker build -t trusch/storaged . || exit $?

exit $?

