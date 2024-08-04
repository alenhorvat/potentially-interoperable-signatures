package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/signatures/bbs"
	"github.com/coinbase/kryptology/pkg/signatures/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/gtank/merlin"
)

func main() {
	// curve := curves.BLS12381(&curves.PointBls12381G2{})
	// Input JSON object
	jsonInput := []byte(`{
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
		     }`)

	var jsonData interface{}
	err := json.Unmarshal(jsonInput, &jsonData)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// JSON to canonical form
	transformedData := walkAndTransformJSON(jsonData)

	transformedJSON, err := json.MarshalIndent(transformedData, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling transformed JSON: %v", err)
	}

	fmt.Println(string(transformedJSON))

	// Canonical form
	transformedData = walkAndReplace(transformedData)
	transformedJSON2, err := json.MarshalIndent(transformedData, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling transformed JSON: %v", err)
	}
	fmt.Println(string(transformedJSON2))

	// JSON canonical to CBOR
	resultCbor, err := cbor.Marshal(transformedData)
	if err != nil {
		panic(err)
	}
	fmt.Println("cbor:", hex.EncodeToString(resultCbor))

	// JSON to pointer
	jsonPointer := make(map[string]interface{})

	_ = walkAndTransformToJSONPointer(transformedData, "", jsonPointer)
	output, err := json.MarshalIndent(jsonPointer, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling transformed JSON: %v", err)
	}

	fmt.Println(string(output))

	// JSON pointer to array
	jsonArray := jsonPointerToArray(jsonPointer)

	fmt.Println(jsonArray)

	msgs := arrayToScalar(jsonArray)
	// msgs := []curves.Scalar{
	// 	curve.Scalar.New(2),
	// 	curve.Scalar.New(3),
	// 	curve.Scalar.New(4),
	// 	curve.Scalar.New(5),
	// }
	fmt.Println("msgs", msgs)

	revealAll(msgs)
	// test(msgs)

	//
	//		pk, sk, err := bbs.NewKeys(curve)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		generators, err := new(bbs.MessageGenerators).Init(pk, len(jsonArray))
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		sig, err := sk.Sign(generators, msgs)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		err = pk.Verify(sig, generators, msgs)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		fmt.Println("pk.Verify", err)
	//
	//		// msg to proof message
	//		var proofMsgs []common.ProofMessage
	//
	//		for _, element := range msgs {
	//			proofMsgs = append(proofMsgs, common.ProofSpecificMessage{
	//				Message: element,
	//			})
	//		}
	//
	//		println("proofMsgs", proofMsgs)
	//
	//		// Create a Proof of Knowledge signature
	//		pok, err := bbs.NewPokSignature(sig, generators, proofMsgs, rand.Reader)
	//		println("pok, err: ", pok, err)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		nonce := curve.Scalar.Random(rand.Reader)
	//		transcript := merlin.NewTranscript("TestPokSignatureProofWorks")
	//		transcript.AppendMessage([]byte("nonce"), nonce.Bytes())
	//		pok.GetChallengeContribution(transcript)
	//		okm := transcript.ExtractBytes([]byte("signature proof of knowledge"), 64)
	//		challenge, err := curve.Scalar.SetBytesWide(okm)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		pokSig, err := pok.GenerateProof(challenge)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		revealedMsgs := map[int]curves.Scalar{
	//			0: msgs[0],
	//		}
	//		// Manual verify to show how when used in conjunction with other ZKPs
	//		transcript = merlin.NewTranscript("TestPokSignatureProofWorks")
	//		pokSig.GetChallengeContribution(generators, revealedMsgs, challenge, transcript)
	//		transcript.AppendMessage([]byte("nonce"), nonce.Bytes())
	//		// okm = transcript.ExtractBytes([]byte("signature proof of knowledge"), 64)
	//		// vChallenge, err := curve.Scalar.SetBytesWide(okm)
	//		// if err != nil {
	//		// 	panic(err)
	//		// }
	//
	//		// Use the all-inclusive method
	//		transcript = merlin.NewTranscript("TestPokSignatureProofWorks")
	//		println(pokSig.Verify(revealedMsgs, pk, generators, nonce, challenge, transcript))
	//
	//		// Create a nonce
	//		// nonce := curve.Scalar.Random(rand.Reader)
	//		// transcript := merlin.NewTranscript("TestPokSignatureProofWorks")
	//		// pok.GetChallengeContribution(transcript)
	//		// transcript.AppendMessage([]byte("nonce"), nonce.Bytes())
	//		// okm := transcript.ExtractBytes([]byte("signature proof of knowledge"), 64)
	//		// challenge, err := curve.Scalar.SetBytesWide(okm)
	//		// if err != nil {
	//		// 	panic(err)
	//		// }
	//
	//		// pokSig, err := pok.GenerateProof(challenge)
	//		// println("pokSig, err: ", pokSig, err)
	//		// if err != nil {
	//		// 	panic(err)
	//		// }
	//		// println("pokSig.Verify: ", pokSig.VerifySigPok(pk))
	//
	//		// revealedMsgs := map[int]curves.Scalar{
	//		// 	2: msgs[2],
	//		// 	3: msgs[3],
	//		// }
	//
	//		// // Manual verify to show how when used in conjunction with other ZKPs
	//		// transcript = merlin.NewTranscript("TestPokSignatureProofWorks")
	//		// pokSig.GetChallengeContribution(generators, revealedMsgs, challenge, transcript)
	//		// transcript.AppendMessage([]byte("nonce"), nonce.Bytes())
	//
	//		// okm = transcript.ExtractBytes([]byte("signature proof of knowledge"), 64)
	//		// vChallenge, err := curve.Scalar.SetBytesWide(okm)
	//		// println("vChallenge, err: ", vChallenge, err)
	//		// if err != nil {
	//		// 	panic(err)
	//		// }
	//		// println("challenge vs vChallenge", challenge, challenge.Cmp(vChallenge))
	//
	//		// // Use the all-inclusive method
	//		// transcript = merlin.NewTranscript("TestPokSignatureProofWorks")
	//		// result := pokSig.Verify(revealedMsgs, pk, generators, nonce, challenge, transcript)
	//
	//		// println(result)
}

func test(msgs []curves.Scalar) {
	curve := curves.BLS12381(&curves.PointBls12381G2{})
	pk, sk, err := bbs.NewKeys(curve)
	if err != nil {
		panic(err)
	}
	// _, ok := pk.value.(*curves.PointBls12381G2)
	generators, err := new(bbs.MessageGenerators).Init(pk, 4)
	if err != nil {
		panic(err)
	}
	sig, err := sk.Sign(generators, msgs[:4])
	if err != nil {
		panic(err)
	}

	// Here we need to set the messages that are revealed
	proofMsgs := []common.ProofMessage{
		&common.ProofSpecificMessage{
			Message: msgs[1],
		},
		&common.RevealedMessage{
			Message: msgs[2],
		},
		&common.ProofSpecificMessage{
			Message: msgs[0],
		},
		&common.RevealedMessage{
			Message: msgs[3],
		},
	}
	revealedMsgs := map[int]curves.Scalar{
		2: msgs[2],
		3: msgs[3],
	}
	// msg to proof message
	// var proofMsgs []common.ProofMessage

	// for _, element := range msgs {
	// 	proofMsgs = append(proofMsgs, common.ProofSpecificMessage{
	// 		Message: element,
	// 	})
	// }
	// for i := 0; i < len(msgs); i += 1 {
	// 	proofMsgs = append(proofMsgs, &common.RevealedMessage{
	// 		Message: msgs[i],
	// 	})
	// }

	println(len(proofMsgs), len(msgs))

	println("proofMsgs", proofMsgs)

	pok, err := bbs.NewPokSignature(sig, generators, proofMsgs, rand.Reader)
	if err != nil {
		panic(err)
	}
	nonce := curve.Scalar.Random(rand.Reader)
	transcript := merlin.NewTranscript("TestPokSignatureProofWorks")
	pok.GetChallengeContribution(transcript)
	transcript.AppendMessage([]byte("nonce"), nonce.Bytes())
	okm := transcript.ExtractBytes([]byte("signature proof of knowledge"), 64)
	challenge, err := curve.Scalar.SetBytesWide(okm)
	if err != nil {
		panic(err)
	}

	pokSig, err := pok.GenerateProof(challenge)
	if err != nil {
		panic(err)
	}
	fmt.Println(pokSig.VerifySigPok(pk))

	// Manual verify to show how when used in conjunction with other ZKPs
	transcript = merlin.NewTranscript("TestPokSignatureProofWorks")
	pokSig.GetChallengeContribution(generators, revealedMsgs, challenge, transcript)
	transcript.AppendMessage([]byte("nonce"), nonce.Bytes())
	okm = transcript.ExtractBytes([]byte("signature proof of knowledge"), 64)
	// vChallenge, err := curve.Scalar.SetBytesWide(okm)

	// Use the all-inclusive method
	transcript = merlin.NewTranscript("TestPokSignatureProofWorks")
	fmt.Println(pokSig.Verify(revealedMsgs, pk, generators, nonce, challenge, transcript))

}

// Create a proof were all messages are revealed
func revealAll(msgs []curves.Scalar) {
	curve := curves.BLS12381(&curves.PointBls12381G2{})
	pk, sk, err := bbs.NewKeys(curve)
	if err != nil {
		panic(err)
	}
	// Generators are created by the issuer
	generators, err := new(bbs.MessageGenerators).Init(pk, len(msgs))
	if err != nil {
		panic(err)
	}
	// Generators are created by the issuer
	generators2, err := new(bbs.MessageGenerators).Init(pk, len(msgs))
	if err != nil {
		panic(err)
	}
	sig, err := sk.Sign(generators, msgs)
	if err != nil {
		panic(err)
	}

	// msg to proof message
	proofMsgs := []common.ProofMessage{}
	revealedMsgs := make(map[int]curves.Scalar)

	for i := 0; i < len(msgs); i += 1 {
		// Note: seems that order doesn't matter
		proofMsgs = append(proofMsgs, &common.RevealedMessage{
			Message: msgs[i],
		})
		revealedMsgs[i] = msgs[i]
	}

	// PoK signature
	pok, err := bbs.NewPokSignature(sig, generators2, proofMsgs, rand.Reader)
	if err != nil {
		panic(err)
	}

	nonce := curve.Scalar.Random(rand.Reader)
	transcript := merlin.NewTranscript("TestPokSignatureProofWorks")
	pok.GetChallengeContribution(transcript)
	transcript.AppendMessage([]byte("nonce"), nonce.Bytes())
	okm := transcript.ExtractBytes([]byte("signature proof of knowledge"), 64)
	challenge, err := curve.Scalar.SetBytesWide(okm)
	if err != nil {
		panic(err)
	}

	pokSig, err := pok.GenerateProof(challenge)
	if err != nil {
		panic(err)
	}
	fmt.Println(pokSig.VerifySigPok(pk))

	// Manual verify to show how when used in conjunction with other ZKPs
	transcript = merlin.NewTranscript("TestPokSignatureProofWorks")
	pokSig.GetChallengeContribution(generators, revealedMsgs, challenge, transcript)
	transcript.AppendMessage([]byte("nonce"), nonce.Bytes())
	// okm = transcript.ExtractBytes([]byte("signature proof of knowledge"), 64)
	// vChallenge, err := curve.Scalar.SetBytesWide(okm)

	// Generators are created by the issuer
	generators3, err := new(bbs.MessageGenerators).Init(pk, len(msgs))
	if err != nil {
		panic(err)
	}
	// Use the all-inclusive method
	transcript = merlin.NewTranscript("TestPokSignatureProofWorks")
	fmt.Println("pokSig.Verify:", pokSig.Verify(revealedMsgs, pk, generators3, nonce, challenge, transcript))

}
