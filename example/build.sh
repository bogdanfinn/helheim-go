#!/bin/zsh
# You have to adjust the different paths according to your setup/system
export PY_DIR='/Library/Frameworks/Python.framework/Versions/3.8'
export PY_VERSION='3.8'
CGO_LDFLAGS='-Wl,-rpath /Library/Frameworks/Python.framework/Versions/3.8 -L/Library/Frameworks/Python.framework/Versions/3.8/lib -lpython3.8 -lhelheim_cffi' CGO_CFLAGS="-I/Library/Frameworks/Python.framework/Versions/3.8/include/python3.8" go build main.go

# Specify where application finds the cffi lib file on runtime
# DYLD_LIBRARY_PATH="/Library/Frameworks/Python.framework/Versions/3.8/lib" ./main

# We assume here that you copied the cffi lib file into the applications working directory and we do not need to define the path for the dyld
# ./main