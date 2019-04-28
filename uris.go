package charityhonor

import (
	"math/rand"
	"strings"
	"time"
)

var adjs = []string{
	"Gold",
	"Silver",
	"Bronze",
	"Black",
	"White",
	"Red",
	"Pink",
	"Dark",
	"Light",
	"Radiant",
	"Majestic",
	"Noble",
	"Solar",
	"Lunar",
	"Big",
	"Small",
	"Honest",
	"Pretty",
	"Honorable",
	"Studious",
}

var nouns = []string{
	"Legend",
	"Myth",
	"Tale",
	"Quest",
	"Epic",

	"Crow",
	"Bear",
	"Lion",
	"Elephant",
	"Cougar",
	"Tiger",
	"Eagle",
	"Hawk",
	"Dog",
	"Husky",
	"Warthog",
	"Parakeet",
	"Cat",
	"Dove",
	"Snail",

	"Shield",
	"Helm",
	"Armor",
	"Boots",
	"Doublet",
	"Crown",
	"Chainmail",
	"Breastplate",

	"Mage",
	"Bard",
	"Pally",
	"Sorcerer",
	"Warlock",
	"Wizard",
	"Knight",
	"Paladin",

	"Sun",
	"Star",
	"Moon",
	"Galaxy",
	"Comet",

	"Castle",
	"Fort",
	"Tower",

	"King",
	"Queen",
	"Prince",
	"Jester",
	"Princess",
	"Pauper",
	"Duke",
	"Duchess",
	"Earl",
	"Viscount",
	"Count",
	"Lord",
	"Lady",
}

func init() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
}

func GenerateUri() string {
	numParts := 3
	parts := make([]string, numParts)
	parts[0] = adjs[rand.Intn(len(adjs))]
	for parts[1] == "" || parts[1] == parts[0] {
		parts[1] = adjs[rand.Intn(len(adjs))]
	}
	parts[2] = nouns[rand.Intn(len(nouns))]

	return strings.Join(parts, "")
}
