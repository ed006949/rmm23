# rmm23

RMM, episode 23

Welcome to the Remote Monitoring and Management (episode 23) WiKi.

* должно быть уникальным в масштабах всей инфры:
	* UUID
	* ipHostNumber
	* все сертификаты (отпечатки) и CN в них
	* dn
	* uid
	* uidNumber
	* gid
	* gidNumber


* применение ACL:
	* вся инфра
	* домен
	* сущность домена
	* группа
	* сущность группы
	* пользователь
	* сущность пользователя


* uidNumber: важно только в случае взаимодействия с FS
* gidNumber: основная группа, важно только в случае взаимодействия с FS
* домашний каталог: важно только в случае взаимодействия с OS
* memberOf: членство в группах, выдаётся сервером автоматически отдельным модулем, были замечены несоответствия с актуальной информацией, лучше не опираться, а как и модуль сервера высчитывать из полей member в группах


* основа для генерации конфигурации: синтаксис JunOS
* основа учёта пользователей: LDAP
	* рассматривается возможность вести собственную БД пользователей
		* хоть в NoSQL
		* реализаця встроенных LDAP-демонов для различных "миров":
			* OpenLDAP
			* AD
		* это позволит:
			* увеличить гибкость применения
			* упростит реализацию
			* избавит от зависимости от сторонних LDAP-серверов
			* ускорит обработку
*


- [ ] spec
	- [ ] core
		- [ ] API
		- [ ] read-only slave
			- [ ] w/ caching?
		- [ ] built-in DB
		- [ ] built-in LDAP
			- [ ] frontend:
				- [ ] OpenLDAP
				- [ ] MS AD
		- [ ] frontend
	- [ ] управление доменами:
		- [ ] ACL
		- [ ] управление пользователями
			- [ ] UUID
			- [ ] сущности
				- [ ] ipHostNumber (выдаваемый диапозон IPv4 (/27))
				- [ ] AAA
					- [ ] PKI
				- [ ] ACL
				- [ ] 0x00 a.uid
				- [ ] 0x01 b.uid
				- [ ] 0x02 c.uid
				- [ ] 0x03 d.uid
				- [ ] 0x04 e.uid
				- [ ] 0x05 f.uid
				- [ ] 0x06 g.uid
				- [ ] 0x07 h.uid
				- [ ] 0x08 i.uid
				- [ ] 0x09 j.uid
				- [ ] 0x0a k.uid
				- [ ] 0x0b l.uid
				- [ ] 0x0c m.uid (mobile)
				- [ ] 0x0d n.uid (notebook)
				- [ ] 0x0e o.uid
				- [ ] 0x0f p.uid
				- [ ] 0x10 q.uid
				- [ ] 0x11 r.uid
				- [ ] 0x12 s.uid
				- [ ] 0x13 t.uid (tablet)
				- [ ] 0x14 u.uid
				- [ ] 0x15 v.uid
				- [ ] 0x16 w.uid
				- [ ] 0x17 x.uid
				- [ ] 0x18 y.uid
				- [ ] 0x19 z.uid
				- [ ] 0x1a special0x1a
				- [ ] 0x1b special0x1b
				- [ ] 0x1c special0x1c
				- [ ] 0x1d special0x1d
				- [ ] 0x1e special0x1e
				- [ ] 0x1f special0x1f
				- [ ] TODO: передалать 0xXX (на данном этапе это не повлечёт каких-то жутких последствий):
					- [ ] 0x00 subnet
					- [ ] 0x01 gateway
					- [ ] 0x02 a.uid
					- [ ] 0x03 b.uid
					- [ ] 0x04 c.uid
					- [ ] 0x05 d.uid
					- [ ] 0x06 e.uid
					- [ ] 0x07 f.uid
					- [ ] 0x08 g.uid
					- [ ] 0x09 h.uid
					- [ ] 0x0a i.uid
					- [ ] 0x0b j.uid
					- [ ] 0x0c k.uid
					- [ ] 0x0d l.uid
					- [ ] 0x0e m.uid (mobile)
					- [ ] 0x0f n.uid (notebook)
					- [ ] 0x10 o.uid
					- [ ] 0x11 p.uid
					- [ ] 0x12 q.uid
					- [ ] 0x13 r.uid
					- [ ] 0x14 s.uid
					- [ ] 0x15 t.uid (tablet)
					- [ ] 0x16 u.uid
					- [ ] 0x17 v.uid
					- [ ] 0x18 w.uid
					- [ ] 0x19 x.uid
					- [ ] 0x1a y.uid
					- [ ] 0x1b z.uid
					- [ ] 0x1c reserved
					- [ ] 0x1d reserved
					- [ ] 0x1e reserved
					- [ ] 0x1f broadcast
			- [ ] uid
			- [ ] uidNumber
			- [ ] gidNumber
			- [ ] домашний каталог
			- [ ] userPassword
			- [ ] mail
			- [ ] memberOf
			- [ ] cn
			- [ ] AAA
				- [ ] SSH
				- [ ] PKI
				- [ ] MFA
		- [ ] управление группами:
			- [ ] UUID
			- [ ] gid
			- [ ] gidNumber
			- [ ] member
			- [ ] специальные группы
			- [ ] служебные группы
			- [ ] группы подразделений (обмен сообщениями, пр.)
				- [ ] cn
			- [ ] сущности пользователей
				- [ ] ACL (ACL для каждой сущности)
				- [ ] type a..z
			- [ ] VPN-группы подразделений
				- [ ] ACL (ACL для каждой группы)
				- [ ] type a..z
		- [ ] управление хостами:
			- [ ] UUID
			- [ ] uid
			- [ ] uidNumber
			- [ ] gidNumber
			- [ ] memberOf
			- [ ] cn
			- [ ] AAA
				- [ ] SSH
				- [ ] PKI
				- [ ] MFA
			- [ ] ACL
			- [ ] тип хоста:
				- [ ] interim
					- [ ] ASN
					- [ ] upstream ASN
					- [ ] hostname
				- [ ] OpenVPN
					- [ ] listen
					- [ ] hostname
				- [ ] CiscoVPN
					- [ ] listen
					- [ ] hostname
