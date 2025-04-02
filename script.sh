# export CMAKE_BIN="/cmake/bin"
# export PATH="${CMAKE_BIN}:$PATH"
# export VCPKG_ROOT="/workspaces/shell-tester/vcpkg"
# export PATH="${VCPKG_ROOT}:$PATH"

# # sudo apt-get update && sudo apt-get install --no-install-recommends -y zip=3.* && sudo apt-get install --no-install-recommends -y g++=4:* && sudo apt-get install --no-install-recommends -y build-essential=12.* && sudo apt-get clean && rm -rf /var/lib/apt/lists/*

# # # cmake is required by vcpkg
# # wget --progress=dot:giga https://github.com/Kitware/CMake/releases/download/v3.30.5/cmake-3.30.5-Linux-x86_64.tar.gz && tar -xzvf cmake-3.30.5-Linux-x86_64.tar.gz
# # sudo mv cmake-3.30.5-linux-x86_64/ /cmake

# # git clone https://github.com/microsoft/vcpkg.git && ./vcpkg/bootstrap-vcpkg.sh -disableMetrics

# git clone https://git.codecrafters.io/36b33a214956339f debug
# cd debug
# vcpkg install --no-print-usage
# sed -i '1s/^/set(VCPKG_INSTALL_OPTIONS --no-print-usage)\n/' ${VCPKG_ROOT}/scripts/buildsystems/vcpkg.cmake

#!/bin/bash

for i in $(seq 1 100)
do
    echo "Running iteration $i"
    make test_debug > /tmp/test
    if [ $? -ne 0 ]; then
        echo "make test_debug failed on iteration $i"
        exit 1
    fi
done