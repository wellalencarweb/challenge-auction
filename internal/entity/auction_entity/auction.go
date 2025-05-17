package auction_entity

type AuctionStatus string

const (
	Open   AuctionStatus = "OPEN"
	Closed AuctionStatus = "CLOSED"
)

type ProductCondition string

const (
	New         ProductCondition = "NEW"
	Used        ProductCondition = "USED"
	Refurbished ProductCondition = "REFURBISHED"
)

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
}
