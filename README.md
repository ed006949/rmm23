# rmm23

Welcome to the Remote Monitoring and Management (episode 23).

# Development

## TODO

- [ ] QoD (quality of development)
    - [ ] error arraus
        - [ ] Implement error arrays instead of (wrapped) errors
            * Functions return an array of errors.
            * A dedicated checking function receives this error array.
            * The checking function evaluates errors against predefined syslog levels, logs them, and returns the
              highest severity level.
            * The calling function uses the returned severity level to determine subsequent actions.
            * A default error level can be specified for unexpected errors during function calls.
            * Developers must define an array of expected error severity levels; if this array is not provided (empty or
              nil), the default error level will be applied.
        - [ ] Or stop rambling and revise error handling methods.
    - [ ] implement automatic certificate chains building `@src/mod_crypto/method.go`
    - [ ] local chain verification `@src/mod_crypto/method.go`
- [ ] built-in AAA frontends:
    - [ ] `LDAP`
    - [ ] `Radius`

## Notes

* ASN: `uint32`
    * uidNumber == ASN
* ACL: `JunOS` view:
    * style: `JunOS`
        * security policy
        * policy options
        * firewall
    * application order:
        * infra (?)
            * entities (?)
        * domain (?)
            * entities (?)
        * ACL-groups
            * entities
        * user (?)
            * entities (?)
* `uidNumber`: `uint32` only important when interacting with FS
* `gidNumber`: `uint32` primary group, only important when interacting with FS
* `home`: only important when interacting with OS
* `memberOf`:
    * group membership, automatically issued by the server via a separate module, inconsistencies with current
      information have been noted, better not to rely on it, but rather calculate it from the member fields in groups,
      similar to the server module
* `ipHostNumber`:
    * user's IPv4 subnet (`/27`)
* `device` = special `user`

* data must be unique within the entire infrastructure:
    * UUID
        - [x] generate new null-based `UUID` against `DN` while loading legacy data from `LDAP`
    * `ipHostNumber`
    * `dn`
    * certificates (`fingerprint`) and certificates's `CN`
    * `uid`
    * `gid`
    * ASN
    * `uidNumber`
    * `gidNumber`
    * only show `notice`:
        * `uid` + `gid`
        * `uidNumber` + `gidNumber`
        * `uidNumber` + `gidNumber` + ASN


* static file generation:
    * go templates
    * backend: VFS
    * frontend: FS

* User accounting:
    * backend: redis DB
    * frontend: LDAP


* Entities:
    * `0x00` `a.uid`
    * `0x01` `b.uid`
    * `0x02` `c.uid`
    * `0x03` `d.uid`
    * `0x04` `e.uid`
    * `0x05` `f.uid`
    * `0x06` `g.uid`
    * `0x07` `h.uid`
    * `0x08` `i.uid`
    * `0x09` `j.uid`
    * `0x0a` `k.uid`
    * `0x0b` `l.uid`
    * `0x0c` `m.uid` (mobile)
    * `0x0d` `n.uid` (notebook)
    * `0x0e` `o.uid`
    * `0x0f` `p.uid`
    * `0x10` `q.uid`
    * `0x11` `r.uid`
    * `0x12` `s.uid`
    * `0x13` `t.uid` (tablet)
    * `0x14` `u.uid`
    * `0x15` `v.uid`
    * `0x16` `w.uid`
    * `0x17` `x.uid`
    * `0x18` `y.uid`
    * `0x19` `z.uid`
    * `0x1a` `special0x1a`
    * `0x1b` `special0x1b`
    * `0x1c` `special0x1c`
    * `0x1d` `special0x1d`
    * `0x1e` `special0x1e`
    * `0x1f` `special0x1f`

    - [ ] TODO: redo 0xXX (at this stage it won't lead to any terrible consequences):
        * `0x00` `subnet`
        * `0x01` `gateway`
        * `0x02` `a.uid`
        * `0x03` `b.uid`
        * `0x04` `c.uid`
        * `0x05` `d.uid`
        * `0x06` `e.uid`
        * `0x07` `f.uid`
        * `0x08` `g.uid`
        * `0x09` `h.uid`
        * `0x0a` `i.uid`
        * `0x0b` `j.uid`
        * `0x0c` `k.uid`
        * `0x0d` `l.uid`
        * `0x0e` `m.uid` (mobile)
        * `0x0f` `n.uid` (notebook)
        * `0x10` `o.uid`
        * `0x11` `p.uid`
        * `0x12` `q.uid`
        * `0x13` `r.uid`
        * `0x14` `s.uid`
        * `0x15` `t.uid` (tablet)
        * `0x16` `u.uid`
        * `0x17` `v.uid`
        * `0x18` `w.uid`
        * `0x19` `x.uid`
        * `0x1a` `y.uid`
        * `0x1b` `z.uid`
        * `0x1c` `reserved`
        * `0x1d` `reserved`
        * `0x1e` `reserved`
        * `0x1f` `broadcast`

## Core

- [ ] implement go routines (`context`)
- [x] FS i/o
- [ ] network i/o
    - [ ] ssh i/o
    - [ ] API i/o
- [ ] data processor:
    - [ ] go templates
    - [x] JSON
    - [x] XML
    - [ ] LDAP:
        - [x] load from LDAP
        - [ ] mirror changes to LDAP
- [ ] daemon
    - [ ] API
    - [ ] frontend
    - [ ] (?) cluster
        - [ ] slave
            - [ ] w/ caching (?)
            - [ ] readonly slave
        - [ ] multimaster
- [ ] DB
    - [x] redis as DB backend instead of memory
        - [x] db0: main db with indexing
        - [ ] db1: config
        - [ ] db2: `MQTT`
    - [ ] auth front-ends:
        - [ ] (?) built-in
        - [ ] `LDAP`
            - [ ] OpenLDAP
            - [ ] MS AD
- [ ] ACL
    - [ ] sanitize/normalize
    - [ ] weigh/prioritize
    - [ ] summarize/optimize
- [ ] PKI
    - [ ] sanitize/validate
        - [ ] key-cert validator
        - [ ] issuer-cert validator
    - [ ] implement ACME
- [ ] ~~implement `protobuf`~~

## Internal DB Structure

### Domain Management

* UUID
* `labeledURI`:
    * ACL
    * AAA
        * PKI
            * CA cert-key pair

### User Management

* UUID
* `uid`
* `uidNumber`
* `gidNumber`
* `home`
* `userPassword`
* `mail`
* `memberOf`
* `cn`
* `ipHostNumber`
* `labeledURI`:
    * AAA
        * SSH
        * PKI
            * cert-key pairs signed with domain's CA
                * entities
        * MFA
    * ACL
        * entities

### Group Management

* UUID
* `gid`
* `gidNumber`
* `member`
* `labeledURI`:
    * ACL
* special groups
* VPN-groups
    * user permissions
        * `vpn`: allow user to connect
        * `vpn-entity-[a-z]`: allow user entity to connect
        * `vpn-host-${uint32}`: allow user to connect to host
    * ACL
        * `vpn-acl-[a-z][a-z0-9]+`: group ACL
        * `vpn-acl-[a-z][a-z0-9]+-[a-z]`: group's entity ACL

* service groups
* groups `[a-z][a-z0-9]+` (for messaging, etc)
    * `cn`

### Device Management

TODO

* UUID (ASN?)
* `uid`
* `uidNumber`
* `gidNumber`
* `memberOf`
* `cn`
* `labeledURI`:

    * ~~host type:~~
        * ~~provider:~~
            * ~~API:~~
                * ~~host `address`~~
                * ~~AAA~~
            * ~~ASN~~
            * ~~AAA~~
        * ~~interim:~~
            * ~~ASN~~
            * ~~upstream host ASN~~
            * ~~hosting ASN~~
            * ~~AAA~~
        * ~~openvpn:~~
            * ~~URL~~
            * ~~listen `IPAddrPort`~~
            * ~~AAA~~
        * ~~ciscovpn:~~
            * ~~URL~~
            * ~~listen `IPAddrPort`~~
            * ~~AAA~~

    * type: `(provider|interim|openvpn|ciscovpn)`
    * ASN
    * upstream device ASN
    * hosting device ASN
    * URL
    * listen `IPAddrPort`
    * AAA
        * SSH
        * PKI
            * cert-key pairs signed with domain's CA
            * cert-key pairs signed with LE CA
        * MFA
    * ACL

## requirements

### build

* [Go][URL_Go]
    * [Rueidis][URL_Go_Rueidis]
    * [Rueidis OM][URL_Go_Rueidis_OM]
    * [LDAP][URL_Go_LDAP]

### run

* [Redis][URL_Redis] server with modules:
    * [RediSearch][URL_RediSearch]
    * [RedisJSON][URL_RedisJSON]

[URL_Redis]: https://github.com/redis/redis

[URL_RediSearch]: https://github.com/RediSearch/RediSearch

[URL_RedisJSON]: https://github.com/RedisJSON/RedisJSON

[URL_Go]: https://golang.org/

[URL_Go_Rueidis]: https://github.com/redis/rueidis

[URL_Go_Rueidis_OM]: https://github.com/redis/rueidis/om

[URL_Go_LDAP]: https://github.com/go-ldap/ldap/v3

## questions

### JSON

* pass string data surrounded with `"` ?

## MUST / SHOULD / LATER

### MUST

- Redis-backed storage with Rueidis + RedisJSON
- optimistic concurrency with `redis:",ver"`
- absolute expiration support with `redis:",exat"`
- stable Redis key for every entry with `redis:",key"`
- LDAP entry identity:
    - `entryUUID`
    - `entryDN`
- LDAP object model:
    - `objectClass`
    - `structuralObjectClass`
- LDAP operational metadata:
    - `creatorsName`
    - `createTimestamp`
    - `modifiersName`
    - `modifyTimestamp`
- schema model:
    - `attributeTypes`
    - `objectClasses`
    - `syntaxes`
    - `matchingRules`
- schema discovery support via `subschemaSubentry`

### SHOULD

- `schemaVersion` for future migrations
- `hasSubordinates` for tree navigation
- schema-driven validation before save/update
- index strategy for common LDAP search fields
- support for standard LDAP core schema concepts
- compatibility-focused schema responses for LDAP clients
- normalized comparisons via matching rules
- internal object type/kind discriminator

### LATER

- sync / replication support
- `entryCSN`
- `numSubordinates`
- multi-master or slave replication
- advanced LDAP operational attributes beyond the basic set
- conflict resolution policies for sync
- schema import/export tooling
- richer search/filter optimization
- password policy operational attributes 
