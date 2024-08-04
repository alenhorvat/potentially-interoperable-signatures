# Potentially Interoperable Signatures

This is a working draft about Potentially Interoperable Signature formats. This
work is motivated by the latest "Bring Your Own Format" movement. So, this is our
+1 contribution :D

## Deterministic JSON (dJSON)

Everybody loves JSON. If we could only have an algorithm to transform it into
the exactly same form, no matter how we permute the claims or order the array
elements.

Proposal is the following

- identify all the arrays
- transform the array into a map
- key in the map is SHA256 hash of the value

Of course we need to do this iteratively:

```javascript
// Function to walk through and transform JSON
function walkAndTransformJSON(input) {
    // Determine the type of the input
    switch typeOf(input) {
        // Case when input is a list
        case "list":
            newMap = {}                     // Initialize a new map to store transformed items
            indexList = {}                  // Initialize a map to keep track of indices and their counts

            // Iterate over each item in the list
            for each item in input {
                transformed = walkAndTransformJSON(item)   // Recursively transform each item
                transformedByte, err = jsonMarshal(transformed)   // Convert transformed item to JSON bytes
                index = computeHash(transformedByte)   // Compute a hash for the JSON bytes

                // Check if the index already exists
                if index in indexList {
                    _index = index
                    indexList[_index] += 1
                    index = format("%s_%d", _index, indexList[_index])   // Append a counter to the index
                } else {
                    indexList[index] = 0   // Initialize counter for new index
                }

                newMap[index] = transformed   // Add the transformed item to the new map with the computed index
            }
            return newMap   // Return the transformed map

        // Case when input is an object (map)
        case "object":
            // Recursively process each value in the map
            for each (key, value) in input {
                input[key] = walkAndTransformJSON(value)   // Transform each value
            }
    }
    return input   // Return the transformed input
}
```

Example:

Input JSON

```json
{
  "users": [
        {"name": "Charlie", "age": 35, "active": null},
        {"name": "Bob", "age": 25},
        {"name": "Alice", "age": 30, "active": [true, true, true]}
    ],
    "data": {
        "info": [
            {"id": 1, "value": "x"},
            {"id": 1, "value": "x"}
        ]
    },
    "emptyArray": [],
    "emptyObject": {},
    "stringValue": "example",
    "numberValue": 123.456,
    "booleanValue": true,
    "nullValue": null
}
```

result:

```json
{
  "booleanValue": true,
  "data": {
    "info": {
      "_id:699baea133190f7d31da6b58921fb8cb12f99e0810159393d6483f11182d37de": {
        "id": 1,
        "value": "x"
      },
      "_id:699baea133190f7d31da6b58921fb8cb12f99e0810159393d6483f11182d37de_1": {
        "id": 1,
        "value": "x"
      }
    }
  },
  "emptyArray": {},
  "emptyObject": {},
  "nullValue": null,
  "numberValue": 123.456,
  "stringValue": "example",
  "users": {
    "_id:bceb53a452e08b6d6388957ab80dfaadf8163f4adf2c253cf1db2b4592a79fbc": {
      "active": {
        "_id:b5bea41b6c623f7c09f1bf24dcae58ebab3c0cdd90ad966bc43a45b44867e12b": true,
        "_id:b5bea41b6c623f7c09f1bf24dcae58ebab3c0cdd90ad966bc43a45b44867e12b_1": true,
        "_id:b5bea41b6c623f7c09f1bf24dcae58ebab3c0cdd90ad966bc43a45b44867e12b_2": true
      },
      "age": 30,
      "name": "Alice"
    },
    "_id:f7cdf5ddc0e9fd5c708864e6d7c5d940b3503a0cb0d53c944a227d0b7f832e35": {
      "age": 25,
      "name": "Bob"
    },
    "_id:f82b65d160429a29d3d1f16ca2709c04799bdddfc843406691a27bc719309e45": {
      "active": null,
      "age": 35,
      "name": "Charlie"
    }
  }
}
```

In the last step, we sort all the elements and replace them with indexes 0, 1, 2, ...

```json
{
  "booleanValue": true,
  "data": {
    "info": {
      "0": {
        "id": 1,
        "value": "x"
      },
      "1": {
        "id": 1,
        "value": "x"
      }
    }
  },
  "emptyArray": {},
  "emptyObject": {},
  "nullValue": null,
  "numberValue": 123.456,
  "stringValue": "example",
  "users": {
    "0": {
      "active": {
        "0": true,
        "1": true,
        "2": true
      },
      "age": 30,
      "name": "Alice"
    },
    "1": {
      "age": 25,
      "name": "Bob"
    },
    "2": {
      "active": null,
      "age": 35,
      "name": "Charlie"
    }
  }
}
```

Following this algorithm we always end up with the same canonical structure. We call this Deterministic JSON (dJSON).

## JSON <-> CBOR

It is very easy to transform the canonical structure to CBOR, but it's also easy
to apply the same algorithm to a CBOR encoded structure.

Note: since CBOR supports more types than JSON and also has a support for custom
types, generic CBOR <-> JSON transformation is non-trivial.

## JWS and Data Integrity

If we operate with the dJSON, we can easily express a JWS-protected JSON in a
Data Integrity format in a way that the two signatures match exactly.
Note: We need to check if hex encoding is supported by Data Integrity.

How?

If dJSON is signed using JWS, content signed reads
`BASE64URL(UTF8(JWS Protected Header)) || '.' || BASE64URL(JWS Payload)`

If we take a VCDM v2, transform it into dJSON format, transform the "proof"
claim into the JWS header claim, we get exactly the same result. Note, if
"proof" has claims:

- protected
- header

the transformation is trivial since same claims appear in the JSON serialised JWS.

## BBS+ and other multi-message signatures

Luckily the story doesn't end with JWS and Data Integrity. We can also transform
the dJSON form using JSON pointers:

```json
{
  "/booleanValue": true,
  "/data/info/0/id": 1,
  "/data/info/0/value": "x",
  "/data/info/1/id": 1,
  "/data/info/1/value": "x",
  "/nullValue": null,
  "/numberValue": 123.456,
  "/stringValue": "example",
  "/users/0/active/0": true,
  "/users/0/active/1": true,
  "/users/0/active/2": true,
  "/users/0/age": 30,
  "/users/0/name": "Alice",
  "/users/1/age": 25,
  "/users/1/name": "Bob",
  "/users/2/active": null,
  "/users/2/age": 35,
  "/users/2/name": "Charlie"
}
```

We can also add a protected header to the structure:

```json
{
  "protected": {
    "foo": "bar"
  },
  "payload": {
    "/booleanValue": true,
    "/data/info/0/id": 1,
    "/data/info/0/value": "x",
    "/data/info/1/id": 1,
    "/data/info/1/value": "x",
    "/nullValue": null,
    "/numberValue": 123.456,
    "/stringValue": "example",
    "/users/0/active/0": true,
    "/users/0/active/1": true,
    "/users/0/active/2": true,
    "/users/0/age": 30,
    "/users/0/name": "Alice",
    "/users/1/age": 25,
    "/users/1/name": "Bob",
    "/users/2/active": null,
    "/users/2/age": 35,
    "/users/2/name": "Charlie"
  }
}
```

Note: we can express the full structure using JSON pointers.

We can now easily sign this as a multi-message ZKP-based signature, using BBS/BBS+ or any other.
We can selectively disclose claims and share with the verifier the corresponding dJSON structure, like:

```json
{
  "protected": {
    "foo": "bar"
  },
  "payload": {
    "booleanValue": true,
    "data": {
      "info": {
        "0": {
          "id": 1,
          "value": "x"
        },
        "1": {
          "id": 1,
          "value": "x"
        }
      }
    }
  }
}
```

We can encode the signature as JWS or data integrity and the verifier will easily verify the signature.

## JSON-LD and vocabularies

Vocabularies are extremely important as they help us defining and interpreting
the different claims. We can easily reference (link + hash of the external
object) JSON-LD context, SHACL schema or vocabulary/ontology in any other
format.

Preferred way (but not limited to) is using the JSON-LD context.

## The challenge

This is an early draft. Feel free to open an issue.
