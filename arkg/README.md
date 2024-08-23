# Asynchronous Remote Key Generation (ARKG)

This is a simplified summary of ARKG using elliptic curves. There are many
interesting uses of ARKG (see the references) and generalisations and
(quantum-resistant) extensions of ARKG have been proposed.

ARKG has a simple, yet powerful property: Alice can create Bob's public keys
without knowing Bob's private key (that correspond to the generated public
keys).

Math behind a basic scheme is extremely simple: Bob creates his public/private key pair and shares the public key with Alice. Alice creates an ephemeral key pair, computes a shared secret and derives the public key from Bob's original public key and the shared secret - Alice can do this without Bob's involvement. Alice can then (at any later stage) share the ephemeral public key with Bob, and Bob can derive the corresponding private key from his private key and from the shared secret.

Let's draft the steps:

- Bob creates his private/public key pair

```pseudocode
(d, pub) = generateECDSAKeyPair()
```

- Bob shares his public key `pub` with Alice
- Alice then generates an ephemeral private/public key pair

```pseudocode
(d_e, pub_e) = generateECDSAKeyPair()
```

- Alice now computes a shared secret as

```pseudocode
sharedSecret = d_e * pub =
             = pub_e * d
```

Note that if Bob knows the `pub_e`, he can compute the same shared secret.

- Alice computes the derived public key as

```pseudocode
pub_derived = pub + HKDF(sharedSecret) * G
```

HKDF is a key derivation function.

- Alice shares the ephemeral public key `pub_e` with Bob
- Bob computes the shared secret as

```pseudocode
sharedSecret = d * pub_e
```

and derives the private key as:

```pseudocode
d_derived = d + HKDF(sharedSecret)
```

We can check that the public key is correct:

```pseudocode
pub_derived' = d_derived * G =
             = (d + HKDF(sharedSecret)) * G =
             = pub + HKDF(sharedSecret) * G =
             = pub_derived
```

Note that this is a simplistic summary of the math behind the scheme. For more
details check out the references.

## Potential use cases

- Key recovery framework
- Asynchronous Verifiable Credential issuance
- Messaging
  - Open question: can this be used to replace Signal's central key distribution server?
- Blockchain transactions
- other

## References

- <https://www.yubico.com/blog/yubico-proposes-webauthn-protocol-extension-to-simplify-backup-security-keys/>
- <https://hackernoon.com/blockchain-privacy-enhancing-technology-series-stealth-address-i-c8a3eb4e4e43>
- <https://github.com/w3c/webauthn/issues/1640>
- <https://github.com/Yubico/webauthn-recovery-extension>
- <https://www.ietf.org/archive/id/draft-bradleylundberg-cfrg-arkg-02.html>
- <https://github.com/Yubico/webauthn-recovery-extension/tree/master/benchmarks>
