package charityhonor

import (
	"math/rand"
	"strings"
	"time"
)

var adjs = []string{
	"Big",
	"Black",
	"Bronze",
	"Calm",
	"Careful",
	"Charming",
	"Dark",
	"Fearless",
	"Friendly",
	"Gentle",
	"Gold",
	"Good",
	"Honest",
	"Honorable",
	"Jolly",
	"Light",
	"Lunar",
	"Majestic",
	"Noble",
	"Pink",
	"Pretty",
	"Radiant",
	"Red",
	"Silver",
	"Small",
	"Solar",
	"Studious",
	"White",
}

var nouns = []string{
	"Epic",
	"Legend",
	"Myth",
	"Quest",
	"Tale",

	//Animals
	"Bear",
	"Cat",
	"Crow",
	"Cougar",
	"Dog",
	"Dove",
	"Eagle",
	"Elephant",
	"Hawk",
	"Husky",
	"Lion",
	"Parakeet",
	"Snail",
	"Tiger",
	"Warthog",

	//Armor
	"Armor",
	"Boots",
	"Breastplate",
	"Chainmail",
	"Crown",
	"Doublet",
	"Helm",
	"Shield",

	//RPG classes
	"Bard",
	"Knight",
	"Mage",
	"Pally",
	"Paladin",
	"Sorcerer",
	"Warlock",
	"Wizard",

	//Defenders
	"Defender",
	"Guard",
	"Guardian",
	"Protector",
	"Sentinel",
	"Sentry",
	"Warden",


	//Jobs
	"Doctor",
	"Nurse",
	"Sailor",

	//Celestial things
	"Sun",
	"Star",
	"Moon",
	"Galaxy",
	"Comet",

	//Fortresses
	"Castle",
	"Fort",
	"Fortress",
	"Tower",

	//Titles
	"Captain",
	"Count",
	"Duke",
	"Duchess",
	"Earl",
	"Emperor",
	"Empress",
	"Jester",
	"King",
	"Lady",
	"Lord",
	"Pauper",
	"Prince",
	"Princess",
	"Queen",
	"Viscount",
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
