package constants

type OrderStatus int
type OrderStatusString string

const (
	Pending        OrderStatus = 100
	PendingPayment OrderStatus = 200
	PaymentSuccess OrderStatus = 300
	Expired        OrderStatus = 400

	PendingString        OrderStatusString = "pending"
	PendingPaymentString OrderStatusString = "pending_payment"
	PendingSuccessString OrderStatusString = "payment_success"
	ExpiredString        OrderStatusString = "expired"
)

var mapStringToInt = map[OrderStatusString]OrderStatus{
	PendingString:        Pending,
	PendingPaymentString: PendingPayment,
	PendingSuccessString: PaymentSuccess,
	ExpiredString:        Expired,
}

var mapIntToString = map[OrderStatus]OrderStatusString{
	Pending:        PendingString,
	PendingPayment: PendingPaymentString,
	PaymentSuccess: PendingSuccessString,
	Expired:        ExpiredString,
}

func (p OrderStatusString) String() string {
	return string(p)
}

func (p OrderStatus) Int() int {
	return int(p)
}

func (p OrderStatus) GetStatusString() OrderStatusString {
	return mapIntToString[p]
}

func (ps OrderStatusString) GetStatusInt() OrderStatus {
	return mapStringToInt[ps]
}
