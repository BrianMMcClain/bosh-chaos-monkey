BOSH Chaos Monkey
===

A tool developed to help test the resiliency of your BOSH releases

Usage
---

`./bosh-chaos-monkey -config config.json`

config.json
---

This file defines both what to test and how to test it. Refer to [config.json.go](https://github.com/BrianMMcClain/bosh-chaos-monkey/blob/master/config.json.example) for an example.

|Field|Description|Example|
|---|---|---|
|directorURL|URL for BOSH director|https://192.168.50.6:25555|
|username|Username to use when connecting to BOSH|admin|
|password|Password to use when connecting to BOSH|myPass|
|caPath|Path to Director CA cert|./ca.cert
|deploymentName|Name of deployment to test|zookeeper
|killInterval|Frequency to kill VMs (in seconds)|60|

Additional Flags
---
|Flag|Description|
|---|---|
|`-dry`|Dry run only, only logs which machine would be killed|

TODO
---
- Kill the VM at the IaaS level, currently all delete commands are going through BOSH, which is a good start, but could be better
- Ensure both basic and UAA auth work
- Add support for SOCKS5 proxy through the jumpbox that [bbl](https://github.com/cloudfoundry/bosh-bootloader) sets up