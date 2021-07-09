package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"crypto/sha1"
	"math/big"
)




//ecc签名--私钥
func EccSignature(sourceData []byte, privateKeyFilePath string) ([]byte, []byte) {
	//1，打开私钥文件，读出内容
	file, err := os.Open(privateKeyFilePath)
	if err != nil {
		panic(err)
	}
	info, err := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//2,pem解密
	block, _ := pem.Decode(buf)
	//x509解密
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//哈希运算
	hashText := sha1.Sum(sourceData)
	//数字签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashText[:])
	if err != nil {
		panic(err)
	}
	rText, err := r.MarshalText()
	if err != nil {
		panic(err)
	}
	sText, err := s.MarshalText()
	if err != nil {
		panic(err)
	}
	defer file.Close()
	return rText, sText
}

//ecc认证

func EccVerify(rText, sText, sourceData []byte, publicKeyFilePath string) bool {
	//读取公钥文件
	file, err := os.Open(publicKeyFilePath)
	if err != nil {
		panic(err)
	}
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)

	//x509
	publicStream, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//接口转换成公钥
	publicKey := publicStream.(*ecdsa.PublicKey)
	hashText := sha1.Sum(sourceData)
	var r, s big.Int
	r.UnmarshalText(rText)
	s.UnmarshalText(sText)
	//认证
	res := ecdsa.Verify(publicKey, hashText[:], &r, &s)
	defer file.Close()
	return res
}
