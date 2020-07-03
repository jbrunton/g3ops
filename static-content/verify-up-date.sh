#!/bin/bash

make statik

if [[ `git status --porcelain` ]]; then
  echo "Static content is out of date. Run \"make statik\" and commit the changes."
  exit 1
else
  echo "Static content is up to date."
fi
