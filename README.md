# rmm23

Welcome to the Remote Monitoring and Management (episode 23).

# development

## notes

* работа с ACL: синтаксис JunOS
	* JunOS:
		* security policy
		* policy options
		* firewall
	* применение ACL:
		* вся инфра
		* домен
		* сущности в домене (?)
		* группа
		* сущности в группе (?)
		* пользователь
		* сущности пользователя
* создание файлов на FS: go templates
* учёт пользователей: LDAP
	* рассматривается возможность вести собственную БД пользователей
		* SQL/NoSQL для хранения дерева
		* NoSQL для обмена diff
		* реализаця встроенных LDAP-демонов для различных "миров":
			* OpenLDAP
			* MS AD
		* это позволит:
			* увеличить гибкость применения
			* упростит реализацию
			* избавит от зависимости от сторонних LDAP-серверов
			* ускорит обработку
* уникальность данных в масштабах всей инфры:
	* UUID
	* ipHostNumber
	* dn
	* все сертификаты (отпечатки) и CN в них
	* uid
	* uidNumber
	* gid
	* gidNumber
* uidNumber: важно только в случае взаимодействия с FS
* gidNumber: основная группа, важно только в случае взаимодействия с FS
* home: важно только в случае взаимодействия с OS
* memberOf: членство в группах, выдаётся сервером автоматически отдельным модулем, были замечены несоответствия с актуальной информацией, лучше не опираться, а как и модуль сервера высчитывать из полей member в группах
* ipHostNumber: user's IPv4 subnet (/27)
* entities:
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

## core

- [ ] implement go routines
- [x] FS i/o
- [ ] network i/o
	- [ ] ssh i/o
	- [ ] API i/o
- [ ] data processor:
	- [ ] go templates
	- [ ] JSON
	- [x] XML
	- [ ] LDAP
- [ ] daemon
	- [ ] API
	- [ ] frontend
	- [ ] cluster
		- [ ] slave
			- [ ] w/ caching (?)
			- [ ] readonly slave
		- [ ] multimaster
- [ ] built-in DB (?)
	- [ ] storage: SQL/NoSQL
	- [ ] sync: NoSQL (diff)
	- [ ] implement MQTT (?)
	- [ ] built-in LDAP
		- [ ] frontend
			- [ ] OpenLDAP
			- [ ] MS AD
- [ ] ACL
	- [ ] sanitize/normalize
	- [ ] weigh/prioritize
	- [ ] summarize/optimize
- [ ] PKI
	- [ ] key-cert validator
	- [ ] issuer-cert validator
	- [ ] implement ACME
- [ ] implement protobuf

## internal DB structure

### domain management

* UUID
* ACL
* AAA
	* PKI

### user management

* UUID
* entities
	* ipHostNumber
	* AAA
		* PKI
	* ACL
* uid
* uidNumber
* gidNumber
* home
* userPassword
* mail
* memberOf
* cn
* AAA
	* SSH
	* PKI
	* MFA

### group management

* UUID
* gid
* gidNumber
* member
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
	* cn

### host management

* UUID (ASN?)
* uid
* uidNumber
* gidNumber
* memberOf
* cn
* AAA
	* SSH
	* PKI
	* MFA
* ACL
* upstream host ASN `${uint32}`
* host `address`
* API `address`
* ASN `${uint32}`
* listen `IPAddrPort`
* host type `(provider|interim|openvpn|ciscovpn)`
