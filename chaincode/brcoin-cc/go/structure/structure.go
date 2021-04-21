package structure

//바락코인 지갑 구조체
type BarakWallet struct {
	Regdate  int64  `json:"regdate"`  // 지갑 등록 일자
	Password string `json:"password"` // 지갑 패스워드 (인증서로 사용할 예정)
	//Addinfo  string         `json:"addinfo"`
	JobType string        `json:"jobType"` //job 내용 (실행 내용)
	JobArgs string        `json:"jobArgs"` //job params
	JobDate int64         `json:"jobDate"` //job 진행 시간
	Balance []BalanceInfo `json:"balance"` // 잔고 리스트
	//Pending  map[int]string `json:"pending"`
	Nonce string `json:"nonce"` //논스값
}

// 잔고 구조체
type BalanceInfo struct {
	Balance    string `json:"balance"`    // 잔고
	TokenId    int    `json:"tokenId"`    // 토큰 id
	UnlockDate int64  `json:"unlockdate"` //토큰 거래 정지 날짜
}

// Token BRC000 - 토큰
type Token struct {
	Owner          string `json:"owner"`          // 생성자 지갑 주소
	Symbol         string `json:"symbol"`         // 토큰 심볼
	CreateDate     int64  `json:"createdate"`     // 토큰 생성 시간 (체인코드에서 생성)
	TotalSupply    string `json:"totalsupply"`    //토큰 총 공급량
	ReservedAmount string `json:"reservedamount"` // 소유자의 토큰 보유량?
	Token          int    `json:"token"`          // 토큰 아이디 (체인코드에서 생성)
	Name           string `json:"name"`           //토큰 이름
	Information    string `json:"information"`    // 토큰 정보
	URL            string `json:"url"`            // 토큰 관련 url
	// Image          string           `json:"image"`       // 토큰 이미지 우선 사용 보류
	Decimal int            `json:"decimal"` //토큰 상태라는데 (소수점 단위 같음)
	Reserve []TokenReserve `json:"reserve"` //토큰 예약 리스트
	Type    string         `json:"type"`    //xhzms 타입
	JobType string         `json:"jobType"` //job 내용 (실행 내용)
	JobArgs string         `json:"jobArgs"` //job params
	JobDate int64          `json:"jobDate"` //job 진행 시간 (체인코드에서 생성)
}

// 토큰 예약 리스트 구조체
type TokenReserve struct {
	Address    string `json:"address"`    //지갑 주소
	Value      string `json:"value"`      //토큰 량?
	UnlockDate int64  `json:"unlockdate"` //거래 제한 날짜
}

//결과 구조체
type jsonResponse struct {
	Key           string `json:"key"`
	ResultFlag    bool   `json:"resultFlag"`
	ResultCode    string `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}
