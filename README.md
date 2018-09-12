# vault-iam-request

[![CircleCI](https://circleci.com/gh/dvianello/vault-iam-request.svg?style=svg)](https://circleci.com/gh/dvianello/vault-iam-request) 
[![Go Reports card](https://goreportcard.com/badge/github.com/dvianello/vault-iam-request)](https://goreportcard.com/report/github.com/dvianello/vault-iam-request)
[![Coverage Status](https://coveralls.io/repos/github/dvianello/vault-iam-request/badge.svg?branch=master)](https://coveralls.io/github/dvianello/vault-iam-request?branch=master)


A small golang script to build the STS request used to perform Vault IAM-based auth in AWS. 

---

#### Usage
```bash
vault-iam-request [OPTIONS]

Application Options:
  -r, --role= The Vault role to authenticate against
  -j, --json  Output data in JSON format
  -f, --file= Write output to file instead of stdout

Help Options:
  -h, --help  Show this help message
```


##### Credentials
`vault-iam-request` will need valid AWS credentials to be able to talk to STS. As we're re-using the 
Vault cli codebase, we _automagically_ support authentication via the standard environment variables (`AWS_*`), 
credentials stored in `~/.aws/credentials` as well instance profiles. 
