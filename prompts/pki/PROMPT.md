- A small bash-script `./pki.sh`:
	* Works in the local directory only
	* with `index.txt`, `serial` and `crlnumber` support
	* Run options prepended with `-`:
		* verbose
			* report about every `op` (`op`, `op` result)
		* verify
			* `name`
				* find `$1.crt.pem`
					* is not `CA`
						* check if exist
						* verify `key` against `cert`
						* verify `cert` against `CA`
					* is `CA`
						* verify `key` against `cert`
		* create
			* `name`
				* is not `CA`
					* load `CA`
					* check if not exist
					* create `key`
					* create `cert` against `key` and `CA`
				* is `CA`
					* check if not exist
					* create `key`
					* create `cert` against `key` and `CA`
		* revoke
			* is not `CA`
				* load `CA`
				* check if exist
				* verify `cert` against `CA`
		* delete
			* is not `CA`
				* revoke `name`
				* delete `name`

# Examples

## Run

### Make it executable

chmod +x ./pki.sh

### Be verbose verification

./pki.sh -verbose ....

### Create CA

./pki.sh -create CA

#### Create cert for "alice"

./pki.sh -create alice

### Verify alice

./pki.sh -verify alice

### Revoke alice

./pki.sh -revoke alice

### Delete alice

./pki.sh -delete alice

### List all

./pki.sh -list

```
STATUS    SERIAL       CN                            EXPIRY/REVOCATION         NOTES
VALID     1001         bob                           2025-08-10 16:00:00Z     -
REVOKED   1000         alice                         2025-08-09 16:00:00Z     revoked on expiry=2025-08-10 16:00:00Z
EXPIRED   1002         charlie                       2024-08-10 16:00:00Z     -

```
