package components

type ItemCardProps struct {
	ImageURL string
	AltText  string
	Title    string
	SubTitle string
	Price    string
}
type SectionProps struct {
	Title string
	Items []ItemCardProps
}

var items = []SectionProps{
	{
		Title: "UraMaki",
		Items: []ItemCardProps{
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A chef presenting a plate of sushi",
				Title:    "Spicy Tuna Maki",
				SubTitle: "A tantalizing blend of spicy tuna, cucumber, and avocado, harmoniously rolled in nori and seasoned rice.",
				Price:    "$5",
			},
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A bowl of ramen with chopsticks",
				Title:    "Mango Maki",
				SubTitle: "Tempura-fried shrimp, cucumber, and cream cheese embrace a center of fresh avocado, delivering a satisfying contrast of textures.",
				Price:    "$12",
			},
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A plate of assorted tempura",
				Title:    "Salmon Maki",
				SubTitle: "Shiitake mushrooms, avocado, and pickled daikon radish nestle within a roll of seasoned rice, coated with nutty sesame seeds.",
				Price:    "$5",
			},
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A sushi roll with avocado and fish",
				Title:    "Tuna Maki",
				SubTitle: "A vibrant assortment of julienned carrots, bell peppers, and cucumber, tightly encased in a nori-wrapped rice roll.",
				Price:    "$5",
			},
		}},
	{
		Title: "Special Rolls",
		Items: []ItemCardProps{
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A chef presenting a plate of sushi",
				Title:    "Sunrise Bliss",
				SubTitle: "A delicate combination of fresh salmon, cream cheese, and asparagus, rolled in orange-hued tobiko for a burst of sunrise-inspired flavors.",
				Price:    "$16",
			},
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A bowl of ramen with chopsticks",
				Title:    "Mango Tango",
				SubTitle: "Tempura shrimp, cucumber, and avocado dance alongside sweet mango slices, drizzled with a tangy mango sauce.",
				Price:    "$16",
			},
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A plate of assorted tempura",
				Title:    "Truffle Indulgence",
				SubTitle: "Decadent slices of black truffle grace a roll of succulent wagyu beef, cucumber, and microgreens, culminating in an exquisite umami symphony.",
				Price:    "$16",
			},
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A sushi roll with avocado and fish",
				Title:    "Pacific Firecracker",
				SubTitle: "Spicy crab salad, tempura shrimp, and jalapeño peppers combine in a fiery ensemble, accented with a chili-infused aioli.",
				Price:    "$16",
			},
			{
				ImageURL: "/static/img/item1.avif",
				AltText:  "A sushi roll with avocado and fish",
				Title:    "Eternal Eel",
				SubTitle: "An enchanting blend of eel tempura, foie gras, and cucumber, elegantly layered with truffle oil and gold leaf for a touch of opulence.",
				Price:    "$16",
			},
		}},
}
