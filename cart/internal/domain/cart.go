package domain

type CartItem struct {
	SKU   uint32
	Count uint16
}

type FullCartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}
