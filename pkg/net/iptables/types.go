package iptables

type DNATRule struct {
	// incoming dst port like cli --dport
	DPort string

	// ToHost:ToPort like cli --to-destination
	ToHost string
	ToPort string
}

type SNATRule struct {
	// output dst host like cli -d
	DHost string
}
