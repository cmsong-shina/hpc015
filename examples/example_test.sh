#! /bin/sh

# Set current occupants as 3
curl localhost:8888/cs/count -X POST -d 3


# Get current occupants
curl localhost:8888/cs/count