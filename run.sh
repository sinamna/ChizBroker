#!/bin/bash

go build -o build/server/broker
docker build . -t broker
docker run broker