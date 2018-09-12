# vault-iam-request

[![CircleCI](https://circleci.com/gh/dvianello/vault-iam-request.svg?style=svg)](https://circleci.com/gh/dvianello/vault-iam-request) 
[![Go Reports card](https://goreportcard.com/badge/github.com/dvianello/vault-iam-request)](https://goreportcard.com/report/github.com/dvianello/vault-iam-request)
[![Coverage Status](https://coveralls.io/repos/github/dvianello/vault-iam-request/badge.svg?branch=master)](https://coveralls.io/github/dvianello/vault-iam-request?branch=master)


A small golang program to build the STS request used to perform Vault IAM-based auth in AWS. 

---

#### Usage
```
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


#### Concourse
`vault-iam-request` was developed to allow the integration of Vault as
 a [Concourse secret backend](https://concourse-ci.org/creds.html) without hardcoding a long-lived token. The output
 of `vault-iam-request` can be directly fed into `vault_remote_auth_param` and Concourse will use it to auth against
 a Vault role and obtain a token. 
 
 As STS call are timestamped, Concourse won't be able to use the same call again to re-auth if the token expires. For 
 this reason, it's highly recommended to configure the Vault role to issue a  
 [periodic token](https://www.vaultproject.io/docs/concepts/tokens.html#periodic-tokens) instead of a normal token. As 
 long as Concourse will renew the token within the `period`, the ttl of the token will be reset and it will keep 
 working. Should the token expire and not being renewable, ATC should be restarted with a new STS call.