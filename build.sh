#!/bin/bash

cd view
npm run build
cd ..
go build -i -v
echo 'build rp done'
