#!/usr/bin/env bash

PARENT_DIR=$(cd $(dirname $0);cd ..; pwd)

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

#source /Users/edy/opt/miniconda3/bin/activate
#conda activate base
echo "py env activated"
cd ${PARENT_DIR}
echo ${PARENT_DIR}
python3  scripts/import_yapi_data.py 
