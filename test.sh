#!/bin/bash

abc="abc"

if echo abc | grep -v "b"; then
  echo "yes"
else
  echo "no"
fi