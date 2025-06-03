# name

desc

### to implement

- [ ] create prj structure for "config"
	- [ ] implement logic
	- [ ] implement funcs

- [ ] add built-in ldap server (to be able to receive online updates from upstream)

- [ ] add config generation for:
	- [ ] openvpn
	- [ ] ocserv
	- [ ] amnezia
	- [ ] local server FW (?)
	- [ ] site FW

- [ ] add config apply engine

- [ ] implement queueing

### to fix

- [ ] логика взаимодействия с LDAP (скорость обработки)
- [ ] отдельные модули генерации файлов для разных сервисов:
	- [ ] доделать go templates
	- [x] openvpn
	- [ ] ssl vpn
	- [ ] amnezia
		- [ ] wg
		- [ ] awg
- [ ] добавить работу с vfs
	- [ ] большое кол-во файлов
	- [ ] синхронизация с fs 
