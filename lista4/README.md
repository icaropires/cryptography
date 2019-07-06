# Lista 4

## Executing applications

### (Question 1) RSA

``` bash
go run rsa.go [p] [q] [key]
```

### (Question 2) List points of elliptic curve with order

``` bash
# Example: a=1 b=1 p=23
go run ecc.go list_points.go [a] [b] [p]
```

### (Question 3) Ciphering and Deciphering using elliptic curve

``` bash
# Example: a=1 b=1 p=23 Gx=0 Gy=1 Px=19 Py=5
# Warning: Point P must be on curve
go run ecc.go ecc_cipher_decipher.go [a] [b] [p] [Gx] [Gy] [Px] [Py]
```

### (Question 4) Digital Signature using elliptic curve

#### Signing

**Warning:** File will be changed without asking !!!

``` bash
# Using curve secp256k1
go run ecc.go sha_256.go sign.go [filename]
```

#### Verifying

``` bash
# Using curve secp256k1
# Public key is printed when signing
go run ecc.go sha_256.go verify.go [filename] [PublicKeyX] [PublicKeyY]
```
