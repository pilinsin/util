
[liboqs install](https://github.com/open-quantum-safe/liboqs)

On Ubuntu:
```bash
sudo apt install astyle cmake gcc ninja-build libssl-dev python3-pytest python3-pytest-xdist unzip xsltproc doxygen graphviz python3-yaml
```
```bash
cd $(ANYPATH)
git clone -b main https://github.com/open-quantum-safe/liboqs.git
cd liboqs

mkdir build && cd build
cmake -GNinja -DBUILD_SHARED_LIBS=ON ..
ninja

sudo ninja install
```

[liboqs-go install](https://github.com/open-quantum-safe/liboqs-go)
```bash
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
```
```bash
cd $(ANYPATH2) && git clone https://github.com/open-quantum-safe/liboqs-go
```

Edit $(ANYPATH2)/liboqs-go/.config/liboqs.pc:
```bash
LIBOQS_INCLUDE_DIR=/usr/local/include
LIBOQS_LIB_DIR=/usr/local/lib
```
```bash
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$(ANYPATH2)/liboqs-go/.config
```
```bash
go get github.com/open-quantum-safe/liboqs-go/oqs
```
