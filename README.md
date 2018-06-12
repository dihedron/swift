# A minimalistic OpenStack Object Storage (aka Swift) Client

Tis project provides a minimalistic implementation of an openStack Swift (Object Storage v1) API client; it provides a very small subset of the official OpenStack CLI functionalities, and only those pertaining to ```objects```. 

## Purpose

The purpose of this client is to provide a self-contained, easily deployable OpenStack Swift CLI with no dependencies (no Python, non PIP) for use in VM and Docker container deployments, where there is the need to download from or upload to an Object Storage while keeping the footprint to a minimum. 
For better functionality-to-file-size ratio, use UPX compressor after building the CLI: the executable may shrink from something less than 12 MB.

## How to use

Simply run the cmmand with no arguments to get a help screen.

```bash
$ swift

This program provides a minimalistic OpenStack Swift v1 client with the ability
to list all objects in a bucket, optionally filter the list, put a new file into 
an existing bucket (upload), retrieve a file from a bucket (download), and delete 
it.

Usage:
  swift [command]

Examples:
swift [command] [args...]

Available Commands:
  about       Retrieve information about an object in a bucket
  get         Retrieve an object from a bucket
  help        Help about any command
  list        List objects in a bucket
  put         Store an object into a bucket
  remove      Remove an object from a bucket

Flags:
  -c, --config string   configuration file (default is $HOME/.swift.yaml)
  -h, --help            help for swift

Use "swift [command] --help" for more information about a command.

```

### The ```list``` command

To get a lexicographically sorted list of objects in an existing bucket (aka container), use the ```list``` command as follows:

```bash
$ swift list my-bucket
```

The ```list``` command accepts an optional additional parameter specifying a regular expression that will be applied to all oject names to filter thos of interest; for instance, running the following:

```bash
$ swift list my-bucket "^exam.*\.gz$"
```

will return all obejcts beginning with ```exam``` and ending with the ```.gz``` extension, e.g. ```example.tar.gz```. Names that do not match the regular expression are omitted.

### The ```get``` command

To download an existing object and write it to STDOUT (which can then be redirected to any file or piped into another command), use the following:

```bash
$ swift get my-bucket my-object > /tmp/myfile.out
```

Otherwise, to write the object directly to file, specify the destination as an additional parameter:

```bash
$ swift get my-bucket my-object myfile.out
```

### The ```put``` command

To upload a local file or some data read from a stream (e.g. STDIN) to the object store, use the following command:

```bash
$ swift put my-bucket my-object < myfile.out
```

To specify the name of the input file directly on the CLI, use:

```bash
$ swift put my-bucket my-object myfile.out
```

