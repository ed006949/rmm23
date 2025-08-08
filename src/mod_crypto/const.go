package mod_crypto

const (
	__                   = " "
	_EC                  = "EC"
	_ECDH                = "ECDH"
	_ED                  = "ED"
	_RSA                 = "RSA"
	_DSA                 = "DSA"
	_CERTIFICATE         = "CERTIFICATE"
	_CERTIFICATE_        = _CERTIFICATE + __
	_PRIVATE             = "PRIVATE"
	_KEY                 = "KEY"
	_REQUEST             = "REQUEST"
	_PRIVATE_KEY         = _PRIVATE + __ + _KEY
	__PRIVATE_KEY        = __ + _PRIVATE_KEY
	_RSA_PRIVATE_KEY     = _RSA + __ + _PRIVATE_KEY
	_EC_PRIVATE_KEY      = _EC + __ + _PRIVATE_KEY
	_CRL                 = "CRL"
	_X500                = "X500"
	_X509                = "X509"
	_X509_CRL            = _X509 + __ + _CRL
	_CERTIFICATE_REQUEST = _CERTIFICATE + __ + _REQUEST
)
