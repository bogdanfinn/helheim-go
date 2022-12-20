#!/bin/zsh
### MacOS ####
# You have to adjust the different paths according to your setup/system. This is for working locally on macos
export PY_DIR='/Library/Frameworks/Python.framework/Versions/3.10'
export PY_VERSION='3.10'
CGO_LDFLAGS='-Wl,-rpath /Library/Frameworks/Python.framework/Versions/3.10 -L/Library/Frameworks/Python.framework/Versions/3.10/lib -lpython3.10 -lhelheim_cffi' CGO_CFLAGS='-I/Library/Frameworks/Python.framework/Versions/3.10/include/python3.10' go build main.go

# Specify where application finds the cffi lib file on runtime when not copied in applications working directory
# DYLD_LIBRARY_PATH="/Library/Frameworks/Python.framework/Versions/3.10/lib" ./main

# We assume here that you copied the cffi lib file into the applications working directory and we do not need to define the path for the .dylib
# ./main


### Ubuntu ####
# You have to adjust the different paths according to your setup/system. This is for working locally on ubuntu
# CGO_LDFLAGS='-Wl,-rpath /usr/include/python3.10 -L/usr/include/python3.10 -lpython3.10 -lhelheim_cffi' CGO_CFLAGS="-I/usr/include/python3.10" go build main.go

# Specify where application finds the cffi lib file on runtime when not copied in applications working directory
# LD_LIBRARY_PATH="/usr/include/python3.10" ./main

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
