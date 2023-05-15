![](https://drive.google.com/uc?export=view&id=1-QlQX0ThujDRQK7eP_CeT2wTjJprkbeQ)
![GitHub](https://img.shields.io/github/license/wf001/rouph)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/wf001/rouph)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/wf001/rouph?filename=src%2Frouphc%2Fgo.mod)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/wf001/rouph/go.yaml?branch=main)

Rouph is an open source programming language, made with reference to [rui314/chibicc](https://github.com/rui314/chibicc/tree/reference) and [here](https://www.sigbus.info/compilerbook).

# Installation
There are three way to install rouph; using a precompiled binary, install from source or using docker. Using docker is easiest and recommended. Please note that Rouph is available to run only on x86-64 architecture, GNU/Linux distribution.

To install a precompiled binary, download the zip package from [here](https://github.com/wf001/rouph/releases/tag/v1.0.0). Make sure to put it on `$PATH`.

To compile from source, you need Go1.15 or later runtime. After Go environment setup, clone the source code, fetch dependencies and build by running the following command,

```
git clone https://github.com/wf001/rouph.git
```
```
cd rouph
```
```
make
```

To use docker, run the following command.

```
git clone https://github.com/wf001/rouph.git
```
```
cd rouph
```
```
docker build -t rouph:latest .
```
```
docker run -ti rouph:latest
```

# Usage
To compile `rouph` source code, you must provide the command as following.

```
# build only
rouph build main.rouph
# build and execute immediately
rouph run main.rouph
```
Sample Rouph source code are also placed on `sample`. If you want to build sample/hello-world.rouph, provide the command as following,

```
rouph build sample/hello-world.rouph
```
```
./hello-world.rouph
```

If you want to build and run sample/hello-world.rouph, provide the command as following,

```
rouph run sample/hello-world.rouph
```
# Specification
Language specification is defined on [SPEC.md](https://github.com/wf001/rouph/blob/main/SPEC.md).

# License
Licensed under the MIT License.
