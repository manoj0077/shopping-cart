package store

//constant product names
const TROUSERS = "Trousers"
const BELTS = "Belts"
const SHIRTS = "Shirts"
const SUITS = "Suits"
const SHOES = "Shoes"
const TIES = "Ties"

type AmountCart struct {
	Quantity int
	Price    int
	Amount   float64
}
type CostCart map[string]AmountCart

type offer interface {
	calculateDiscountPrice(CostCart) CostCart
}

type TrouserOffer struct {
	Name string
}

func (o TrouserOffer) calculateDiscountPrice(cart CostCart) CostCart {
	tprod, ok := cart[TROUSERS]
	if !ok || tprod.Quantity < 2 {
		return cart
	}
	bprod, ok := cart[BELTS]
	if ok {
		amount := float64(bprod.Quantity) * 0.85 * float64(bprod.Price)
		if amount < bprod.Amount {
			bprod.Amount = amount
			cart[BELTS] = bprod
		}
	}
	sprod, ok := cart[SHOES]
	if ok {
		amount := float64(sprod.Quantity) * 0.85 * float64(sprod.Price)
		if amount < sprod.Amount {
			sprod.Amount = amount
			cart[SHOES] = sprod
		}
	}
	return cart
}

type TwoShirtOffer struct {
	Name string
}

func (o TwoShirtOffer) calculateDiscountPrice(cart CostCart) CostCart {
	sprod, ok := cart[SHIRTS]
	if !ok || sprod.Quantity < 3 {
		return cart
	}
	amount := float64(2*sprod.Price + (sprod.Quantity-2)*45)
	if amount < sprod.Amount {
		sprod.Amount = amount
		cart[SHIRTS] = sprod
	}
	return cart
}

type ThreeShirtOffer struct {
	Name string
}

func (o ThreeShirtOffer) calculateDiscountPrice(cart CostCart) CostCart {
	sprod, ok := cart[SHIRTS]
	if !ok || sprod.Quantity < 3 {
		return cart
	}
	tprod, ok := cart[TIES]
	if ok {
		tprod.Amount = float64(tprod.Price * tprod.Quantity / 2.0)
		cart[TIES] = tprod
	}
	return cart
}

func getOffers() []offer {
	var offers []offer
	offers = append(offers, TrouserOffer{Name: "TwoTrouserOffer"})
	offers = append(offers, TwoShirtOffer{Name: "TwoShirtOffer"})
	offers = append(offers, ThreeShirtOffer{Name: "ThreeShirtOffer"})
	return offers
}
