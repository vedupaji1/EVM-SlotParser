#!/bin/bash

solc --storage-layout --pretty-json -o $PWD/tempDirForSolc --overwrite ./contracts/OneInch_TestContract.sol
cp ./tempDirForSolc/* ./storageLayout.json
rm -rf tempDirForSolc