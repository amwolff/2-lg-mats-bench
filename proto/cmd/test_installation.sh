#!/bin/bash
docker run --rm -v "$(pwd):/work" uber/prototool:latest \
prototool grpc matmult.proto \
--address $1 \
--method matmult.Performer/MultiplyMatrices \
--data \
'{
    "multiplier": {
        "columns": [
            {
                "coefficients": [
                    1,
                    0,
                    0
                ]
            },
            {
                "coefficients": [
                    0,
                    1,
                    0
                ]
            },
            {
                "coefficients": [
                    0,
                    0,
                    1
                ]
            }
        ]
    },
    "multiplicand": {
        "columns": [
            {
                "coefficients": [
                    1,
                    0,
                    0
                ]
            },
            {
                "coefficients": [
                    0,
                    1,
                    0
                ]
            },
            {
                "coefficients": [
                    0,
                    0,
                    1
                ]
            }
        ]
    }
}'
