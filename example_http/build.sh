#!/bin/zsh
### MacOS ####
# You have to adjust the different paths according to your setup/system. This is for working locally on macos
export PY_DIR='/Library/Frameworks/Python.framework/Versions/3.8'
export PY_VERSION='3.8'
CGO_LDFLAGS='-Wl,-rpath /Library/Frameworks/Python.framework/Versions/3.8 -L/Library/Frameworks/Python.framework/Versions/3.8/lib -lpython3.8 -lhelheim_cffi' CGO_CFLAGS='-I/Library/Frameworks/Python.framework/Versions/3.8/include/python3.8' go build main.go

# Specify where application finds the cffi lib file on runtime when not copied in applications working directory
# DYLD_LIBRARY_PATH="/Library/Frameworks/Python.framework/Versions/3.8/lib" ./main

# We assume here that you copied the cffi lib file into the applications working directory and we do not need to define the path for the .dylib
# ./main


### Ubuntu ####
# You have to adjust the different paths according to your setup/system. This is for working locally on ubuntu
# CGO_LDFLAGS='-Wl,-rpath /usr/include/python3.8 -L/usr/include/python3.8 -lpython3.8 -lhelheim_cffi' CGO_CFLAGS="-I/usr/include/python3.8" go build main.go

# Specify where application finds the cffi lib file on runtime when not copied in applications working directory
# LD_LIBRARY_PATH="/usr/include/python3.8" ./main

# We assume here that you copied the cffi lib file into the applications working directory and we do not need to define the path for the .so
# ./main

### Windows ####
# You have to adjust the different paths according to your setup/system. This is for working locally on windows
# go env -w CGO_LDFLAGS='-Wl,-rpath C:\Users\Administrator\AppData\Local\Programs\Python\Python38 -LC:\Users\Administrator\AppData\Local\Programs\Python\Python38 -lpython38 -lhelheim_cffi'
# go env -w CGO_CFLAGS="-IC:\Users\Administrator\AppData\Local\Programs\Python\Python38\include"
#
# Compile the application
# go build main.go

# We assume here that you copied the cffi lib file into the applications working directory and we do not need to define the path for the .dll
# ./main.exe
