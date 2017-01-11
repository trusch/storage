#!/bin/bash

docker tag trusch/storaged:latest trusch/storaged:$(git describe)
docker push trusch/storaged:latest
docker push trusch/storaged:$(git describe)

exit $?
