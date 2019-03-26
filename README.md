# Nvidiabeat

Welcome to Nvidiabeat.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/Eifoen/nvidiabeat`

## Getting Started with Nvidiabeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Build

To build the binary for Nvidiabeat run the command below. This will generate a binary
in the same directory with the name nvidiabeat.

```
make
```


### Run

To run Nvidiabeat with debugging output enabled, run:

```
./nvidiabeat -c nvidiabeat.yml -e -d "*"
```

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  Nvidiabeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Nvidiabeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/Eifoen/nvidiabeat
git clone https://github.com/Eifoen/nvidiabeat ${GOPATH}/src/github.com/Eifoen/nvidiabeat
```
## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
