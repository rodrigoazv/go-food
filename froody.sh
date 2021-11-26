#!/usr/bin/env bash

if [ ! -d node_modules ]; then
  yarn install
fi

yarn start -- --host=0.0.0.0 --port 42