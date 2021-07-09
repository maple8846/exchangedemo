package storage

import (
	"database/sql"
	"strings"
	"errors"
	_ "github.com/lib/pq"
	"github.com/kyokan/smallbridge/pkg/domain"
	log "github.com/inconshreveable/log15"
	"github.com/google/uuid"
	"time"
	"github.com/kyokan/smallbridge/pkg/conv"
	"github.com/kyokan/smallbridge/pkg/pb"
	"strconv"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var logger = log.New("module", "storage")

type PGStorage struct {
db *sql.DB
}

type RawOrder struct {
	ID             string
	UserPubkey     string
	MakerAsset     string
	TakerAsset     string
	Quantity	   float32
	Price          float32
	CancelAt       string
	Memo           string
	CreatedAt      string
}


type HTLCHashInfo struct {
	HTLCHash       string
	PHash     	   string
	Timeout        string
	MakerAddress   string
	TakerAddress   string
}


const(         
	posted = iota  //自增值
	intentfill
	done
	cancel
)
func NewPGStorage(dbUrl string) (Storage, error) {
	parts := strings.Split(dbUrl, "://")
	if len(parts) != 2 {
		return nil, errors.New("mal-formed database URL")
	}
	if parts[0] != "postgres" {
		return nil, errors.New("must be a postgres database")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &PGStorage{
		db: db,
	}, nil
}

func (p *PGStorage) Start() error {
	err := p.db.Ping()
	if err != nil {
		return err
	}

	logger.Info("connected to PG")
	return nil
}

func (p *PGStorage) Stop() error {
	return p.db.Close()
}

func (p *PGStorage) CreateOrder(order *domain.Order) (*domain.Order, error) {
	var idStr string

	if order.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		idStr = id.String()
	} else {
		idStr = order.ID
	}


	tempcreated := uint64(time.Now().Unix())
	created := strconv.FormatUint(tempcreated,10)
	err := WithTransaction(p.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO orders (orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdAt, signature, status) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ,$11)`,
			idStr,
			conv.PubkeyToStr(order.UserPubkey),
			order.MakerAsset,
			order.TakerAsset,
			order.Quantity,
			order.Price,
			order.CancelAt,
			order.Memo,
			created,
			"nil",
			posted,	
		)
	


		return err
	})

	if err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:             idStr,
		UserPubkey:     order.UserPubkey,
		MakerAsset:     order.MakerAsset,
		TakerAsset: 	order.TakerAsset,
		Quantity:       order.Quantity,
		Price:          order.Price,
		CancelAt:       order.CancelAt,
		Memo:           order.Memo,
		CreatedAt:      created,
		Status:         posted,
	}, nil
}



func (p *PGStorage) FindOrder(id string) (*domain.Order, error) {
	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdAt, status	FROM orders WHERE orderId = $1 LIMIT 1")
	if err != nil {
	
		return nil, err
	}
	raw := &RawOrder{}
	var status int32;
	err = s.QueryRow(id).Scan(
		&raw.ID,
		&raw.UserPubkey,
		&raw.MakerAsset,
		&raw.TakerAsset,
		&raw.Quantity,
		&raw.Price,
		&raw.CancelAt,
		&raw.Memo,
		&raw.CreatedAt,
		&status,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}


	if err != nil {
		return nil, err
	}
	pubkey, err := conv.PubkeyFromStr(raw.UserPubkey)
	if err != nil {

		return nil, err
	}

	quantity := raw.Quantity

	price := raw.Quantity

	order := &domain.Order{
		ID:             raw.ID,
		UserPubkey:     pubkey,
		MakerAsset:     domain.AssetID(raw.MakerAsset),
		TakerAsset:		domain.AssetID(raw.TakerAsset),
		Quantity:       quantity,
		Price:          price,
		CancelAt:       raw.CancelAt,
		Memo:           raw.Memo,
		CreatedAt:		raw.CreatedAt,
		Status:         status,
			}
	return order, nil

}

func (p *PGStorage) ChainIDForAsset(assetID domain.AssetID) (string, error) {
	s, err := p.db.Prepare("SELECT chainId FROM assets WHERE assetid = $1")
	if err != nil {
		return "", err
	}
	var chainID string
	err = s.QueryRow(assetID).Scan(&chainID)
	if err != nil {
		return "", err
	}

	return chainID, nil
}

func (p *PGStorage) FindOrderbyUser(req *pb.GetOrderRequestV1)(*pb.GetOrderResponseV1, error) {
	var SingleResponse *pb.GetOrderResponse
    var FinalResponse *pb.GetOrderResponseV1
    FinalResponse = new(pb.GetOrderResponseV1)
    var status int
    var userpubkey string

 	if req.MakerAsset != "" && req.TakerAsset !=""  {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE makerAsset = $1 AND takerAsset =$2" )
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.MakerAsset,req.TakerAsset)
    	if !rows.Next(){
    		return FinalResponse, nil
    	}
    } 

    if req.OrderId != "" {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE orderId = $1")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.OrderId)
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())

			SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    	return FinalResponse, nil
    }

    if req.MakerAsset != "" && req.TakerAsset == "" && req.UserPubkey == "" {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE makerAsset = $1")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.MakerAsset)
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			

    		rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())


        	SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.MakerAsset == "" && req.TakerAsset != "" && req.UserPubkey == "" {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE takerAsset = $1")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.TakerAsset)
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	


        	SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.MakerAsset == "" && req.TakerAsset == "" && req.UserPubkey != "" {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE userPubkey = $1")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.UserPubkey)
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.MakerAsset != "" && req.TakerAsset != "" && req.UserPubkey == "" {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE takerAsset = $1 AND makerAsset = $2")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.TakerAsset,req.MakerAsset)
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.MakerAsset != "" && req.TakerAsset == "" && req.UserPubkey != "" {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE userPubkey = $1 AND makerAsset = $2")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.UserPubkey,req.MakerAsset)
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.MakerAsset == "" && req.TakerAsset != "" && req.UserPubkey != "" {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE userPubkey = $1 AND takerAsset = $2")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.UserPubkey,req.TakerAsset)
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
        }
    } else if req.MakerAsset != "" && req.TakerAsset != "" && req.UserPubkey != "" {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders WHERE userPubkey = $1 AND takerAsset = $2 AND makerAsset =$3")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query(req.UserPubkey,req.TakerAsset,req.MakerAsset)
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    }  else {
    	s, err := p.db.Prepare("SELECT orderId, userPubKey, makerAsset, takerAsset, quantity, price, cancelAt, memo, createdat, signature, status FROM orders ")
    	if err != nil {
			return FinalResponse, err
		}
    	rows, _ := s.Query()
    	for rows.Next(){
    		SingleResponse = new(pb.GetOrderResponse)
			rows.Scan(&SingleResponse.OrderId, &userpubkey, &SingleResponse.MakerAsset, &SingleResponse.TakerAsset, &SingleResponse.Quantity,&SingleResponse.Price, &SingleResponse.CancelAt, &SingleResponse.Memo,&SingleResponse.CreatedAt, &SingleResponse.Signature,&status) 
        	
			userPubkey,_ := conv.PubkeyFromStr(userpubkey)

			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	SingleResponse.Status = strconv.Itoa(status)
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } 

	return FinalResponse,nil
}





func (p *PGStorage) RecordFillIntentOrder(fillIntent *domain.FillIntent) (error) {
	var idStr string

	if fillIntent.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return  err
		}
		idStr = id.String()
	} else {
		idStr = fillIntent.ID
	}

	err := WithTransaction(p.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO intent (intentID, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, signature, makerAsset, takerAsset, status,price) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,$13)`,
			idStr,
			fillIntent.UserPubkey,
			fillIntent.CpPubKey,
			fillIntent.OrderId,
			fillIntent.Quantity,
			fillIntent.CreatedAt,
			fillIntent.CancelAt,
			fillIntent.Memo,
	        "nil",
	        fillIntent.MakerAsset,
	        fillIntent.TakerAsset,
	        fillIntent.Status,
	        fillIntent.Price,
		)
		return err
	})

	if err != nil {
		return  err
	}

	return nil
}




func (p *PGStorage) GetFillIntentOrder(req *pb.GetIntentRequestV1) (*pb.GetIntentResponseV1,error) {
 
    var SingleResponse *pb.GetIntentResponse	
	var FinalResponse *pb.GetIntentResponseV1
    FinalResponse = new(pb.GetIntentResponseV1)
   
    var userpubkey string
    var cppubkey   string
   	
    if req.OrderId == "" && req.UserPubKey == "" && req.MakerAsset == "" && req.TakerAsset == ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature, status,price FROM intent")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query()
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())

        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId != "" && req.UserPubKey == "" && req.MakerAsset == "" && req.TakerAsset == ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE orderId = $1")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.OrderId)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())

        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId != "" && req.UserPubKey != "" && req.MakerAsset == "" && req.TakerAsset == ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE orderId = $1 AND userPubKey=$2")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.OrderId, req.UserPubKey)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			

			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())

        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId != "" && req.UserPubKey == "" && req.MakerAsset != "" && req.TakerAsset == ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE orderId = $1 AND makerasset=$2")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.OrderId, req.MakerAsset)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())

        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId != "" && req.UserPubKey == "" && req.MakerAsset == "" && req.TakerAsset != ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE orderId = $1 AND takerasset=$2")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.OrderId, req.TakerAsset)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())

        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId == "" && req.UserPubKey != "" && req.MakerAsset == "" && req.TakerAsset != ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE userPubKey = $1 AND takerasset=$2")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.UserPubKey, req.TakerAsset)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId == "" && req.UserPubKey != "" && req.MakerAsset != "" && req.TakerAsset == ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE userPubKey = $1 AND makerasset=$2")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.UserPubKey, req.MakerAsset)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId == "" && req.UserPubKey == "" && req.MakerAsset != "" && req.TakerAsset != ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE takerasset = $1 AND makerasset=$2")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.TakerAsset, req.MakerAsset)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId != "" && req.UserPubKey == "" && req.MakerAsset != "" && req.TakerAsset != ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE takerasset = $1 AND makerasset=$2 AND orderId = $3")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.TakerAsset, req.MakerAsset, req.OrderId)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId == "" && req.UserPubKey != "" && req.MakerAsset != "" && req.TakerAsset != ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE takerasset = $1 AND makerasset=$2 AND userPubKey = $3")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.TakerAsset, req.MakerAsset, req.UserPubKey)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId != "" && req.UserPubKey != "" && req.MakerAsset == "" && req.TakerAsset != ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE takerasset = $1 AND orderId=$2 AND userPubKey = $3")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.TakerAsset, req.OrderId, req.UserPubKey)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else if req.OrderId != "" && req.UserPubKey != "" && req.MakerAsset != "" && req.TakerAsset == ""{
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE makerasset = $1 AND orderId=$2 AND userPubKey = $3")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.MakerAsset, req.OrderId, req.UserPubKey)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    } else {
    	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, quantity, createdAt, cancelAt, memo, makerasset, takerasset, signature,status,price FROM intent WHERE makerasset = $1 AND orderId=$2 AND userPubKey = $3 AND takerasset=$4")
		if err != nil {
			return nil, err
		}
		rows, _ := s.Query(req.MakerAsset, req.OrderId, req.UserPubKey, req.TakerAsset)
		for rows.Next() {
			SingleResponse = new(pb.GetIntentResponse)
			rows.Scan(&SingleResponse.IntentId, &userpubkey, &cppubkey, &SingleResponse.OrderId,&SingleResponse.Quantity, &SingleResponse.CreatedAt,&SingleResponse.CancelAt, &SingleResponse.Memo, &SingleResponse.MakerAsset,&SingleResponse.TakerAsset, &SingleResponse.Signature, &SingleResponse.Status, &SingleResponse.Price) 
        	/*此处需要返回没有解压缩的key*/
        	userPubkey,_ := conv.PubkeyFromStr(userpubkey)
			SingleResponse.UserPubKey = hexutil.Encode(userPubkey.SerializeCompressed())
        	
			cpPubkey,_ := conv.PubkeyFromStr(cppubkey)
			SingleResponse.CpPubKey = hexutil.Encode(cpPubkey.SerializeCompressed())
        	FinalResponse.Responses = append(FinalResponse.Responses, SingleResponse)
    	}
    }

	return FinalResponse, nil
}


func (p *PGStorage) UpdateOrderStatus(id string, status int) (error) {	
	err := WithTransaction(p.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`UPDATE orders set status=$1 WHERE orderId  = $2`,
			strconv.Itoa(status),
			id,
		)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}


func (p *PGStorage) UpdateIntentStatus(id string, status int) (error) {	
   err := WithTransaction(p.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`UPDATE intent set status=$1 WHERE intentid  = $2`,
			strconv.Itoa(status),
			id,
		)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}




func (p *PGStorage) GetAssetId(idarray *[]string, chainid *[]string) (int,error) {
	assetNum := 0
	var assetname string
	var chainId string
	s, err := p.db.Prepare("SELECT assetid, chainid from assets" )
    if err != nil {
		return 0,err
	}
    rows, err2:= s.Query()
     if err2 != nil {
		return 0,err2
	}
    for rows.Next() {
    	rows.Scan(&assetname,&chainId) 
    	(*idarray)[assetNum] = assetname
    	(*chainid)[assetNum] = chainId
    	assetNum++
    }
    return assetNum,nil
}

	

func (p *PGStorage) RecordFillInfo(para *HTLCHashInfo,cpppara *HTLCHashInfo, orderid string, intentid string) (error) {
	
	fillconfirminfo := &domain.FillConfirmInfo{} 
    fillconfirminfo.Settled = true
    fillconfirminfo.HTLCHash = para.HTLCHash
    fillconfirminfo.PHash    = para.PHash
    fillconfirminfo.Timeout  = para.Timeout


    fillconfirminfo2 := &domain.FillConfirmInfo{} 
    fillconfirminfo2.Settled = true
    fillconfirminfo2.HTLCHash = cpppara.HTLCHash
    fillconfirminfo2.PHash    = cpppara.PHash
    fillconfirminfo2.Timeout  = cpppara.Timeout

	s, err := p.db.Prepare("SELECT intentId, userPubKey, cpPubKey, orderId, createdAt,  memo, takerasset, signature, quantity FROM intent WHERE intentid =$1")
    if err != nil {
		return err
	} 

    err = s.QueryRow(intentid).Scan(
		&fillconfirminfo2.IntentId,
		&fillconfirminfo2.UserPubKey,
		&fillconfirminfo2.CpPubKey,
		&fillconfirminfo2.OrderId,
		&fillconfirminfo2.CreatedAt,
		&fillconfirminfo2.Memo,
		&fillconfirminfo2.AssetId,
		&fillconfirminfo2.Signature,	
		&fillconfirminfo2.Quantity,			
	)

    
    if err != nil {
    	return err
    }

    s, err = p.db.Prepare("SELECT chainid FROM assets WHERE assetId =$1")
    if err != nil {
		return err
	}

	err = s.QueryRow(fillconfirminfo2.AssetId).Scan(
		&fillconfirminfo2.TakerChain,			
	)



	/////////////
	s, err = p.db.Prepare("SELECT orderid, createdAt,  memo, makerasset, signature, quantity FROM orders WHERE orderid =$1")
    if err != nil {
		return err
	} 

    err = s.QueryRow(orderid).Scan(
		&fillconfirminfo.IntentId,
		&fillconfirminfo.CreatedAt,
		&fillconfirminfo.Memo,
		&fillconfirminfo.AssetId,
		&fillconfirminfo.Signature,	
		&fillconfirminfo.Quantity,			
	)
	
 	fillconfirminfo.UserPubKey = fillconfirminfo2.UserPubKey
	fillconfirminfo.CpPubKey   = fillconfirminfo2.CpPubKey
    
    if err != nil {
    	return err
    }

    s, err = p.db.Prepare("SELECT chainid FROM assets WHERE assetId =$1")
    if err != nil {
		return err
	}

	err = s.QueryRow(fillconfirminfo.IntentId).Scan(
		&fillconfirminfo.MakerChain,			
	)
	//////////////

	err = WithTransaction(p.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO fills (userPubKey, cpPubKey, orderId, intentId,makerchain,takerchain, HTLCHash, assetid, chainid, memo, createdAt, signature, settled, preimage, timeout, quantity, cppHTLCHash, cppassetid, cppchainid,cppmemo, cppcreatedad，cppsignature,cpppreimage,cpptimeout,cppquantity) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25)`,
			fillconfirminfo.UserPubKey,
			fillconfirminfo.CpPubKey,
			fillconfirminfo.OrderId,
			fillconfirminfo.IntentId,
			"MakerChain",
			"TakerChain",
			fillconfirminfo.HTLCHash,
			fillconfirminfo.AssetId,
			fillconfirminfo.ChainId,
			fillconfirminfo.Memo,
			fillconfirminfo.CreatedAt,
			"nil",
			fillconfirminfo.Settled,
			fillconfirminfo.PHash,
			fillconfirminfo.Timeout,
			fillconfirminfo.Quantity,
			fillconfirminfo2.HTLCHash,
			fillconfirminfo2.AssetId,
			fillconfirminfo2.ChainId,
			fillconfirminfo2.Memo,
			fillconfirminfo2.CreatedAt,				
			"nil",
			fillconfirminfo2.PHash,
			fillconfirminfo2.Timeout,
			fillconfirminfo2.Quantity,
			)
			return err	
		}) 

	if err != nil {
    	return err
    }
	return nil
}



func (p *PGStorage) GetFillInfo(id string) (*pb.GetFillResponseV1,error) {	
  
  var finalrepsonse *pb.GetFillResponseV1
  finalrepsonse = new(pb.GetFillResponseV1)
  response := &pb.GetFillResponse{}
  response2 := &pb.GetFillResponse{}
  s, err := p.db.Prepare("SELECT userPubKey, cpPubKey, orderId, intentId, MakerChain, TakerChain, chainid,HTLCHash,Assetid,quantity,timeout, preimage, createdAt, memo,signature, cppHTLCHash,cppassetid,cppquantity,cpptimeout, cpppreimage, cppcreatedAt, cppmemo,cppsignature,cppchainid FROM fills WHERE orderid =$1")
  if err != nil {
		return nil,err
  }

  err = s.QueryRow(id).Scan(
		&finalrepsonse.UserPubKey,
		&finalrepsonse.CpPubKey,
		&finalrepsonse.OrderId,
		&finalrepsonse.IntentId,
		&finalrepsonse.Makerchain,
		&finalrepsonse.Takerchain,
		&response.ChainId,
		&response.HTLCHash,
		&response.AssetId,		
		&response.Quantity,
		&response.Timeout,
		&response.PHash,
		&response.CreatedAt,
		&response.Memo,
		&response.Signature,
		&response2.HTLCHash,
		&response2.AssetId,
		&response2.Quantity,
		&response2.Timeout,
		&response2.PHash,
		&response2.CreatedAt,
		&response2.Memo,
		&response2.Signature,
		&response2.ChainId,
	)
    if err != nil {
    	return nil,err
    }

    response2.MakerAddress = ""
    response2.TakerAddress = ""
    response.MakerAddress = ""
    response.TakerAddress = ""
    
    finalrepsonse.Responses = append(finalrepsonse.Responses, response)
    finalrepsonse.Responses = append(finalrepsonse.Responses, response2)
    return finalrepsonse, nil
}




func (p *PGStorage) GetFillconfirmInfo(id string) (*pb.FillConfirmResponseV1,error) {	
  
  var response *pb.FillConfirmResponseV1
  response = new(pb.FillConfirmResponseV1)
 
  s, err := p.db.Prepare("SELECT userPubKey, cpPubKey, orderId, intentId, HTLCHash, MakerChain, TakerChain, quantity,timeout, preimage, createdAt, memo,signature, settled,assetid,chainid FROM fills WHERE orderid =$1")
  if err != nil {
		return nil,err
  }

  err = s.QueryRow(id).Scan(
		&response.UserPubKey,
		&response.CpPubKey,
		&response.OrderId,
		&response.IntentId,
		&response.HTLCHash,
		&response.MakerChain,
		&response.TakerChain,
		&response.Quantity,
		&response.Timeout,
		&response.PHash,
		&response.CreatedAt,
		&response.Memo,
		&response.Signature,
		&response.Settled,
		&response.AssetId,
		&response.ChainId,
	)
    if err != nil {
    	return nil,err
    }

    response.MakerAddress = ""
    response.TakerAddress = ""

    return response, nil
}
