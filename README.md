# rmm23

Welcome to the Remote Monitoring and Management (episode 23).

# Development

## TODO

- [ ] QoD (quality of development)
	- [ ] Implement error arrays instead of (wrapped) errors
		* Functions return an array of errors.
		* A dedicated checking function receives this error array.
		* The checking function evaluates errors against predefined syslog levels, logs them, and returns the highest severity level.
		* The calling function uses the returned severity level to determine subsequent actions.
		* A default error level can be specified for unexpected errors during function calls.
		* Developers must define an array of expected error severity levels; if this array is not provided (empty or nil), the default error level will be applied.
	- [ ] Or stop rambling and revise error handling methods.
- [ ] local chain verification @src/mod_crypto/method.go

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
* `uidNumber`: `uint32` важно только в случае взаимодействия с FS
* `gidNumber`: `uint32` основная группа, важно только в случае взаимодействия с FS
* `home`: важно только в случае взаимодействия с OS
* `memberOf`:
	* членство в группах, выдаётся сервером автоматически отдельным модулем, были замечены несоответствия с актуальной информацией, лучше не опираться, а как и модуль сервера высчитывать из полей member в группах
* `ipHostNumber`:
	* user's IPv4 subnet (`/27`)
* `device` = special `user`

* data must be unique within the entire infrastructure:
	* UUID
	* `ipHostNumber`
	* `dn`
	* все сертификаты (`fingerprint`) и `CN` в них
	* `uid`
	* `gid`
	* ASN
	* `uidNumber`
	* `gidNumber`
	* only show `notice`:
		* `uid` + `gid`
		* `uidNumber` + `gidNumber`
		* `uidNumber` + `gidNumber` + ASN


* формирование файлов для создания на FS: go templates

* Учёт пользователей: LDAP
	* Рассматривается возможность вести собственную БД пользователей:
		* SQL/NoSQL для хранения дерева
		* NoSQL для обмена diff
		* Реализация встроенных LDAP-демонов для различных "миров":
			* OpenLDAP
			* MS AD
	* Это позволит:
		* увеличить гибкость применения
		* упростит реализацию
		* избавит от зависимости от сторонних LDAP-серверов
		* ускорит обработку


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

	- [ ] TODO: передалать 0xXX (на данном этапе это не повлечёт каких-то жутких последствий):
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
	- [ ] JSON
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
		- [x] db0: main db with `rediSearch` for indexing
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
* [Redis][URL_Redigo]
* [RediSearch][URL_redisearch-go]

### run

* [Redis][URL_Redis] server with modules:
	* [RediSearch][URL_RediSearch]
	* [~~RedisJSON~~][URL_RedisJSON]

[URL_Redis]: https://github.com/redis/redis

[URL_RediSearch]: https://github.com/RediSearch/RediSearch

[URL_RedisJSON]: https://github.com/RedisJSON/RedisJSON

[URL_Go]: https://golang.org/

[URL_redisearch-go]: https://github.com/RediSearch/redisearch-go

[URL_redigo]: https://github.com/gomodule/redigo
