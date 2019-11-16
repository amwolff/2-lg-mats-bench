#!/bin/bash

docker run --rm -v "$(pwd):/work" uber/prototool:latest \
prototool grpc \
--address $1 \
--method amwolff.matmult.v1.MatrixProductAPI/Multiply \
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

# TODO(amwolff): add result ('{"result":{"columns":[{"coefficients":[1,0,0]},{"c
# oefficients":[0,1,0]},{"coefficients":[0,0,1]}]}}') validation.
