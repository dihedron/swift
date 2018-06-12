# A minimalistic OpenStack Object Storage (aka Swift) Client

Tis project provides a minimalistic implementation of an openStack Swift (Object Storage v1) API client; it provides a very small subset of the official OpenStack CLI functionalities, and only those pertaining to ```objects```. 

## Purpose

The purpose of this client is to providea self-contained, easily deployable OpenStack Swift CLI with no dependencies (no Python, non PIP) for use in VM and Docker container deployments, where there is the need to donload or upload stuff from an Object Storage and to keep the footprint to a minimum. For better functionality to file size ratio, use UPX compressor after building the CLI.


