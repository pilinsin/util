
[liboqs install](https://github.com/open-quantum-safe/liboqs)

On Ubuntu

sudo apt install astyle cmake gcc ninja-build libssl-dev python3-pytest python3-pytest-xdist unzip xsltproc doxygen graphviz python3-yaml

git clone -b main https://github.com/open-quantum-safe/liboqs.git
cd liboqs

mkdir build && cd build
cmake -GNinja -DBUILD_SHARED_LIBS=ON ..
ninja

sudo ninja install


[liboqs-go install](https://github.com/open-quantum-safe/liboqs-go)

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
cd $HOME && git clone https://github.com/open-quantum-safe/liboqs-go

we assume that liboqs directory is at $OQSDIR
in $HOME/liboqs-go/.config/liboqs.pc:
LIBOQS_INCLUDE_DIR=$OQSDIR/build/include
LIBOQS_LIB_DIR=$OQSDIR/build/lib

export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$HOME/liboqs-go/.config

go get github.com/open-quantum-safe/liboqs-go/oqs
