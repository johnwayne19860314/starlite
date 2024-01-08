#!/usr/bin/env bash


cd grpc
cd proto
echo "current directory "${PWD}

for n in {1..1}; do
    echo looping $n
    if buf generate
    then
        echo "builded grpc-gateway"
        break;
    fi
done

