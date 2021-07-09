package domain

import (
	"../smallbridge/pkg/pb"
	"../smallbridge/wallet"
	"../smallbridge/crpto"
)


/*
Response Sample:
{
    “liq_id”:”1ddaeab7ee1363a3208f8790afd63”,
    “user_pubkey”: “ee92d3b109c714”,
“size”: “0.01000000”,
    “product”: “BTC-TUSD”,
    “address”: “15eotziaf1akHtNVo1ujZRY1yw2T9waSn5”,
    “sending_signature”: “304502201876f2ced75603ccbb9b5371de6cd772c382ef0e5267369f7b63952dc9572ac6022100c50b558daf5adf697d94ed00010afad2a6c61614fe4ef29d7f6f1cc3c410e1f8”,
    “created_at”: “2016-12-08T20:02:28.53864Z”
}
*/

//add by qiyanwei 2018/9/26
//This function take inputs to produce Liqudity proof 
func generateliquditiyproof(liqid string, userPubkey string, size BigInt, product string, address string, c *KeyManager) (liquidityProofV1 *pb.LiquidityProofV1){
	 liqudityProofWithoutSign := &pb.LiquidityProofWithoutSignature{
        id: liqid,
        userPubkey: userPubkey,
		size: size,
    	product: product,
		address: addrss,
		createdAt: proto.uint64(time.Now());
    }
    msg, err := proto.Marshal(liqudityProofWithoutSign) //序列化需要发送的消息
    hash := GethHash(msg)  //取消息的哈希
    crypto.Signature := c.SignData(hash) //对消息进行签名
    str := string(crypto.Signature) 
	liqudityProofV1 := &pb.LiquidityProofV1{
		lps: []&pb.LiquidityProofWithoutSignature{liqudityProofWithoutSign},
		signature: str
	}
	return liqudityProofV1
}


func Testgenerateliquditiyproof() bool {
	liqid ：= '1ddaeab7ee1363a3208f8790afd63'
	size := 3000
	product:= 'BTC-TUSD'
	key:= ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	user_pubkey := key.Public()
	c := &KeyManager{key, 111}
	address := crypto.Keccak256(user_pubkey[1:])[12:]    
	liqudityProofV1:=generateliquditiyproof(liqid,user_pubkey,size,product,address,c)
	liqudityProofWithoutSign := liqudityProofV1.liqudityProofWithoutSign
	msg, err := proto.Marshal(liqudityProofWithoutSign)
	return VerifySignature(msg, user_pubkey,liqudityProofV1.signature)
}