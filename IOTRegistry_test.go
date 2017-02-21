package main

// Notes from Testing popcode
// Public Key: 02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc
// Private Key: 94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20
// Hyperledger address hex 74ded2036e988fc56e3cff77a40c58239591e921
// Hyperledger address Base58: 8sDMfw2Ti7YumfTkbf7RHMgSSSxuAmMFd2GS9wnjkUoX

// Notes from Testing popcode2
// Public Key: 02cb6d65b04c4b84502015f918fe549e95cad4f3b899359a170d4d7d438363c0ce
// Private Key: 60977f22a920c9aa18d58d12cb5e90594152d7aa724bcce21484dfd0f4490b58
// Hyperledger address hex 10734390011641497f489cb475743b8e50d429bb
// Hyperledger address Base58: EHxhLN3Ft4p9jPkR31MJMEMee9G

//Owner1 key
// Public Key: 0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7
// Private Key: 7142c92e6eba38de08980eeb55b8c98bb19f8d417795adb56b6c4d25da6b26c5

// Owner2 key
// Public Key: 02e138b25db2e74c54f8ca1a5cf79e2d1ed6af5bd1904646e7dc08b6d7b0d12bfd
// Private Key: b18b7d3082b3ff9438a7bf9f5f019f8a52fb64647ea879548b3ca7b551eefd65
// func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
// 	_, err := stub.MockInit("1", "", args)
// 	if err != nil {
// 		fmt.Println("INIT", args, "failed", err)
// 		t.FailNow()
// 	}
// }

// func checkInvoke(t *testing.T, stub *shim.MockStub, args []string) {
// 	_, err := stub.MockInvoke("1", "invoke", args)
// 	if err != nil {
// 		fmt.Println("invoke", args, "failed", err)
// 		t.FailNow()
// 	}
// }

// func checkQuery(t *testing.T, stub *shim.MockStub, name string, value string) {
// 	bytes, err := stub.MockQuery("balance", []string{name})

// 	if err != nil {
// 		fmt.Println("Query", name, "failed", err)
// 		t.FailNow()
// 	}
// 	if bytes == nil {
// 		fmt.Println("Query", name, "failed to get value")
// 		t.FailNow()
// 	}
// 	if string(bytes) != value {
// 		fmt.Println("Query value", name, "was not", value, "as expected instead", string(bytes))
// 		t.FailNow()
// 	}
// }

// func mint(t *testing.T, stub *shim.MockStub, counterSeed string) {
// 	createArgs := TuxedoPopsTX.CreateTX{}
// 	createArgs.Address = "74ded2036e988fc56e3cff77a40c58239591e921"
// 	createArgs.Amount = 10
// 	pubKeyBytes, err := hex.DecodeString("03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	createArgs.CreatorPubKey = pubKeyBytes
// 	hexCreatorSig := generateCreateSig(counterSeed, 10, "Test Asset", "Test Data", "74ded2036e988fc56e3cff77a40c58239591e921", "7ff1ac3d9dfc56315ee610d0a15609d13c399cf9c92ba2e32e7b1d25ea5c9494")

// 	createArgs.CreatorSig, err = hex.DecodeString(hexCreatorSig)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	createArgs.Data = "Test Data"
// 	createArgs.Type = "Test Asset"
// 	createArgBytes, err := proto.Marshal(&createArgs)
// 	createArgBytesStr := hex.EncodeToString(createArgBytes)
// 	_, err = stub.MockInvoke("3", "create", []string{createArgBytesStr})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func generateCreateSig(CounterSeedStr string, amount int, assetType string, data string, addr string, privateKeyStr string) string {

// 	privKeyByte, _ := hex.DecodeString(privateKeyStr)

// 	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

// 	message := CounterSeedStr + ":" + addr + ":" + strconv.FormatInt(int64(amount), 10) + ":" + assetType + ":" + data
// 	fmt.Println("Signed Message")
// 	fmt.Println(message)
// 	messageBytes := sha256.Sum256([]byte(message))
// 	sig, _ := privKey.Sign(messageBytes[:])
// 	return hex.EncodeToString(sig.Serialize())
// }

// func possess(t *testing.T, stub *shim.MockStub, counterSeed string, idx int) {
// 	transferArgs := TuxedoPopsTX.TransferOwners{}
// 	transferArgs.Address = "74ded2036e988fc56e3cff77a40c58239591e921"
// 	transferArgs.Data = "Test possess"
// 	transferArgs.PopcodePubKey, _ = hex.DecodeString("02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc")
// 	ownerBytes, _ := hex.DecodeString("0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7")
// 	transferArgs.Owners = [][]byte{ownerBytes}
// 	transferArgs.Output = int32(idx)

// 	//could use a discussion of what's going on with ownerBytes... why encoded to [][]byte? see comment on 137
// 	ownerHex := hex.EncodeToString(ownerBytes)
// 	hexPossessSig := generatePossessSig(counterSeed, idx, "Test possess", ownerHex, "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20")
// 	var err error
// 	transferArgs.PopcodeSig, err = hex.DecodeString(hexPossessSig)
// 	transferArgsBytes, _ := proto.Marshal(&transferArgs)
// 	transferArgsBytesStr := hex.EncodeToString(transferArgsBytes)

// 	_, err = stub.MockInvoke("4", "transfer", []string{transferArgsBytesStr})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func generatePossessSig(CounterSeedStr string, outputIdx int, data string, newOwnersHex string, privateKeyStr string) string {
// 	privKeyByte, _ := hex.DecodeString(privateKeyStr)

// 	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

// 	message := CounterSeedStr + ":" + strconv.FormatInt(int64(outputIdx), 10) + ":" + data

// 	//from comment on 116...what is newOwnersHex?
// 	newOwnersTmp, err := hex.DecodeString(newOwnersHex)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	newOwners := [][]byte{newOwnersTmp}

// 	for _, newO := range newOwners {
// 		message += ":"
// 		message += hex.EncodeToString(newO)
// 	}
// 	// fmt.Printf("Signed message %s \n", message)
// 	mDigest := sha256.Sum256([]byte(message))
// 	sig, _ := privKey.Sign(mDigest[:])
// 	return hex.EncodeToString(sig.Serialize())
// }

// // func generatePossessSig(CounterSeedStr string, outputIdx int, data string, newOwners [][]byte, privateKeyStr string) []byte {
// // 	privKeyByte, _ := hex.DecodeString(privateKeyStr)

// // 	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

// // 	message := CounterSeedStr + ":" + strconv.FormatInt(int64(outputIdx), 10) + ":" + data

// // 	for _, newO := range newOwners {
// // 		message += ":"
// // 		message += hex.EncodeToString(newO)
// // 	}
// // 	// fmt.Printf("Signed message %s \n", message)
// // 	mDigest := sha256.Sum256([]byte(message))
// // 	sig, _ := privKey.Sign(mDigest[:])
// // 	return sig.Serialize()
// // }

// func unitize(t *testing.T, stub *shim.MockStub, counterSeed string) {
// 	unitizeArgs := TuxedoPopsTX.Unitize{}
// 	unitizeArgs.Data = "Test Unitize"
// 	unitizeArgs.DestAddress = "10734390011641497f489cb475743b8e50d429bb"
// 	unitizeArgs.DestAmounts = []int32{10}
// 	unitizeArgs.SourceAddress = "74ded2036e988fc56e3cff77a40c58239591e921"
// 	unitizeArgs.SourceOutput = 0
// 	unitizeArgs.PopcodePubKey, _ = hex.DecodeString("02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc")
// 	ownerSig := generateUnitizeSig(counterSeed, unitizeArgs.DestAddress, 0, []int{10}, "Test Unitize", "7142c92e6eba38de08980eeb55b8c98bb19f8d417795adb56b6c4d25da6b26c5")
// 	unitizeArgs.OwnerSigs = [][]byte{ownerSig}
// 	unitizeArgs.PopcodeSig = generateUnitizeSig(counterSeed, unitizeArgs.DestAddress, 0, []int{10}, "Test Unitize", "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20")
// 	unitizeArgsBytes, _ := proto.Marshal(&unitizeArgs)
// 	unitizeArgsBytesStr := hex.EncodeToString(unitizeArgsBytes)

// 	_, err := stub.MockInvoke("4", "unitize", []string{unitizeArgsBytesStr})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func generateUnitizeSig(CounterSeedStr string, destAddr string, outputIdx int, amounts []int, data string, privateKeyStr string) []byte {
// 	privKeyByte, _ := hex.DecodeString(privateKeyStr)

// 	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

// 	message := CounterSeedStr + ":" + destAddr + ":" + data + ":" + strconv.FormatInt(int64(outputIdx), 10)

// 	for _, amount := range amounts {
// 		message += ":" + strconv.FormatInt(int64(amount), 10)
// 	}
// 	mDigest := sha256.Sum256([]byte(message))
// 	sig, _ := privKey.Sign(mDigest[:])
// 	return sig.Serialize()
// }

// func combine(t *testing.T, stub *shim.MockStub, counterSeed string) {
// 	combineArgs := TuxedoPopsTX.Combine{}
// 	combineArgs.Type = "Test Combined Asset"
// 	combineArgs.PopcodePubKey, _ = hex.DecodeString("02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc")
// 	combineArgs.CreatorPubKey, _ = hex.DecodeString("03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26")
// 	// combineArgs.CreatorSig =

// }

// func TestPopcodeChaincode(t *testing.T) {
// 	bst := new(IOTRegistry)
// 	stub := shim.NewMockStub("tuxedoPops", bst)
// 	checkInit(t, stub, []string{"Hello World"})
// 	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"af5eef44907ccdcc33051d035f32f42de0d093fac2fd9d15923448f6af46bc43","Outputs":null}`)
// 	mint(t, stub, "af5eef44907ccdcc33051d035f32f42de0d093fac2fd9d15923448f6af46bc43")
// 	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"e91d1eab53d597e8e18bb9ebbbaec66d08187d7e14a4a58c8782610ce7c7a74b","Outputs":["{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Data\",\"Type\":\"Test Asset\",\"PrevCounter\":\"af5eef44907ccdcc33051d035f32f42de0d093fac2fd9d15923448f6af46bc43\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
// 	mint(t, stub, "e91d1eab53d597e8e18bb9ebbbaec66d08187d7e14a4a58c8782610ce7c7a74b")
// 	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254","Outputs":["{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Data\",\"Type\":\"Test Asset\",\"PrevCounter\":\"af5eef44907ccdcc33051d035f32f42de0d093fac2fd9d15923448f6af46bc43\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}","{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Data\",\"Type\":\"Test Asset\",\"PrevCounter\":\"e91d1eab53d597e8e18bb9ebbbaec66d08187d7e14a4a58c8782610ce7c7a74b\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
// 	possess(t, stub, "1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254", 1)
// 	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"d3e41e748a7094cc520319623479f97dfb6aae0ea915940b72926384fe8d0e8c","Outputs":["{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Data\",\"Type\":\"Test Asset\",\"PrevCounter\":\"af5eef44907ccdcc33051d035f32f42de0d093fac2fd9d15923448f6af46bc43\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}","{\"Owners\":[\"0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7\"],\"Threshold\":1,\"Data\":\"Test possess\",\"Type\":\"Test Asset\",\"PrevCounter\":\"1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
// 	unitize(t, stub, "d3e41e748a7094cc520319623479f97dfb6aae0ea915940b72926384fe8d0e8c")
// 	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"afab4e267a433fe306d1da4608629ce9a280bde98f7004ff883383d65b9f5948","Outputs":["{\"Owners\":[\"0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7\"],\"Threshold\":1,\"Data\":\"Test possess\",\"Type\":\"Test Asset\",\"PrevCounter\":\"1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
// 	checkQuery(t, stub, "10734390011641497f489cb475743b8e50d429bb", `{"Address":"10734390011641497f489cb475743b8e50d429bb","Counter":"83b298acdf5d7231597ffb776c8f027877ca89cbafa7675a3f177619b0a9ad74","Outputs":["{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Unitize\",\"Type\":\"Test Asset\",\"PrevCounter\":\"d3e41e748a7094cc520319623479f97dfb6aae0ea915940b72926384fe8d0e8c\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
// }
