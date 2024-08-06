# ECDSA experiments

What if we can use the our (HW protected?) private key d to create a signature that looks like it was created by a different private key d'?

Let's do an exercise with ECDSA.

## Algorithm

### Generating an EC Private/Public Key Pair

First, we generate an ECDSA private and public key pair. This will be the
foundation of our cryptographic operations.

```pseudocode
(d, pub) = generateECDSAKeyPair()
```

- `d`: Our private key, kept secret.
- `pub`: The public key, shared with others.

### Message to be Signed

Now, let's define the message we intend to sign. For this example, we'll use a simple "Hello World" message.

```pseudocode
msg = "Hello World"
```

### Hash of the Message

We need to hash the message to ensure integrity and create a fixed-length representation of our message.

```pseudocode
z = Hash(msg, "SHA256")
```

- `z`: The hash of our message.

### Create a Random Integer for "Blinding"

Next, we introduce a blinding factor. This adds an additional layer of security to our cryptographic process.

```pseudocode
alpha = RandomInteger(1, n - 1)
```

- `alpha`: A random integer between 1 and \( n-1 \), used for blinding.

### New Hash to be Signed

We then combine our original hash with the blinding factor to create a new hash, `zPrime`.

```pseudocode
zPrime = z + alpha
```

- `zPrime`: The blinded hash of our message.

### Sign the zPrime

Using our private key, we sign this new hash `zPrime`.

```pseudocode
(r, s) = signECDSA(zPrime, d)
```

- `{r, s}`: The components of our digital signature.

Hypothesis: above signature equals

```pseudocode
(r, s) = signECDSA(zPrime, d) = signECDSA(z, alpha * r^-1 + d)
```

Looks strange, but it works, namely:

```pseudocode
  z + alpha + r*d mod n =
= z + alpha * r * r^-1 + r*d mod n =
= z + r ( alpha * r^-1 + d) mod n  =
= z + r dPrime mod n 

dPrime = alpha * r^-1 + d
pubPrime = ECPublicKey(alpha * r^-1 + d)
```

Let's continue.

### Compute the Inverse of r

To compute our new public key, we first need the modular inverse of `r`.

```pseudocode
rInv = ModularInverse(r, n)
```

- `rInv`: The modular inverse of `r`.

### Compute the new Public Key

We "blind" our public key to account for the blinding factor. This new public key, `pubPrime`, reflects the modified signing process.

```pseudocode
pubPrime = ECPointAddition(a, b, PublicKeyFromScalar(alpha * rInv), pub, p)
```

- `pubPrime`: The adjusted public key.

Note: Not sure if blinding is the best term.

### Verification Process

In the verification step, we check if the original message hash `z` was correctly signed using the modified private key. The verifier will use `pubPrime` for this purpose.

### Verify the Signature

We verify the signature with the original hash `z` and the adjusted public key `pubPrime`.

```pseudocode
result = verifySignECDSA(z, pubPrime, (r, s))
```

- `result`: The outcome of the signature verification. It returns true if the signature is valid and false otherwise.

## The challenge

This is an early draft. Feel free to open an issue.

TODO: check the existing work on the topic!
