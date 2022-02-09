![](https://drive.google.com/uc?export=view&id=1-QlQX0ThujDRQK7eP_CeT2wTjJprkbeQ)

Rouph is an open source programming language, made with reference to [rui314/chibicc](https://github.com/rui314/chibicc/tree/reference) and [here](https://www.sigbus.info/compilerbook).

# Installation
There are three way to install rouph; using a precompiled binary, install from source or using docker. Using docker is easiest and recommended. Please note that Rouph is available to run only on intel x86-64 architecture, GNU/Linux distribution.

To install a precompiled binary, download the zip package from [here](https://github.com/wf001/crude-lang-go/releases). Make sure to put it on `$PATH`.

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
rouph run test.rouph
```
Sample Rouph source code are also placed on `sample`. If you want to build sample/add1.rouph, provide the command as following,

```
rouph run add1.rouph
```

# Specification
Language specification is defined on [SPEC.md](https://github.com/wf001/crude-lang-go/blob/main/SPEC.md)

# License
Licensed under the MIT License.
