module golang.zx2c4.com/wireguard

go 1.18

require (
<<<<<<< HEAD
	github.com/garnoth/pkclient v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20220131195533-30dcbda58838
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	golang.org/x/sys v0.0.0-20220204135822-1c1b9b1eba6a
	golang.zx2c4.com/go118/netip v0.0.0-20211111135330-a4a02eeacf9d
=======
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f
	golang.org/x/sys v0.0.0-20220315194320-039c03cc5b86
>>>>>>> 6a08d81f6bc465a2276c61093d96e567d00beb24
	golang.zx2c4.com/wintun v0.0.0-20211104114900-415007cec224
)

require (
	github.com/miekg/pkcs11 v1.1.1 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
)

replace github.com/miekg/pkcs11 => ../../pkcs11

replace github.com/garnoth/pkclient => ../pkclient
