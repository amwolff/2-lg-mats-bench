#!/bin/bash

docker run --rm -v "$(pwd):/work" uber/prototool:latest \
prototool format -w
