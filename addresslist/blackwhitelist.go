package addresslist

import(
	"net"
	"encoding/xml"
	"log"
)

func Inc(address net.IP) net.IP {
	addressCopy := make(net.IP, len(address))
	copy(addressCopy, address)
	return inc(addressCopy)
}

//Destroys the array passed in
func inc(address net.IP) net.IP {
	address[len(address) - 1]++
	if address[len(address) - 1] == byte(0) && len(address) > 1 {
		return append(inc(address[:len(address) - 1]), address[len(address) - 1])
	}
	return address
}

type IPRange struct {
	Start net.IP
	End net.IP
}

func (pair *IPRange) String() string {
	return pair.Start.String() + " to " + pair.End.String()
}

func (pair *IPRange) Equal(newRange *IPRange) bool {
	return pair.Start.Equal(newRange.Start) && pair.End.Equal(newRange.End)
}

func (pair *IPRange) Has(address net.IP) bool {
	return !IPLess(address, pair.Start) && !IPLess(pair.End, address)
}

func (pair *IPRange) shouldCombine(newRange *IPRange) bool {
	return pair.Has(newRange.Start) || newRange.Start.Equal(Inc(pair.End))
}

type BlackWhiteList struct {
	Whitelist, Blacklist []*IPRange
}

//Destroys the array passed in
func remove(list []*IPRange, index int) []*IPRange {
	for i := index + 1; i < len(list); i++ {
		list[i - 1] = list[i]
	}
	return list[:len(list) - 1]
}

func UnmarshalBlackWhite(data []byte) (*BlackWhiteList, error) {
	retVal := new(BlackWhiteList)
	err := xml.Unmarshal(data, retVal)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

//This is very useful for testing
func (list *BlackWhiteList) Equal(newList *BlackWhiteList) bool {
	if len(list.Whitelist) != len(newList.Whitelist) || len(list.Blacklist) != len(newList.Blacklist) {
		return false
	}
	for i, iprange := range list.Whitelist {
		if !iprange.Equal(newList.Whitelist[i]) {
			return false
		}
	}
	for i, iprange := range list.Blacklist {
		if !iprange.Equal(newList.Blacklist[i]) {
			return false
		}
	}
	return true
}

func (list *BlackWhiteList) String() string {
	retVal := "Whitelist:\n"
	for _, elem := range list.Whitelist {
		retVal += elem.String() + "\n"
	}
	retVal += "BlackList:\n"
	for _, elem := range list.Blacklist {
		retVal += elem.String() + "\n"
	}
	return retVal
}

func (list *BlackWhiteList) Marshal() []byte {
	data, err := xml.Marshal(list)
	if err != nil {
		log.Println(err)
	}
	return data
}

func (list *BlackWhiteList) FilterList(peerlist PeerList) PeerList {
	resultlist := make(PeerList, 0, len(peerlist))
	for _, peeritem := range peerlist {
		keep := false
		for _, elem := range list.Whitelist {
			if elem.Has(peeritem.IP) {
				keep = true
				break
			}
		}
		if !keep {
			continue
		}
		for _, elem := range list.Blacklist {
			if elem.Has(peeritem.IP) {
				keep = false
				break
			}
		}
		if !keep {
			continue
		}
		resultlist = append(resultlist, peeritem)
	}
	return resultlist
}

func (list *BlackWhiteList) BlacklistRange(newRange *IPRange) {
	for i, iprange := range list.Blacklist {
		if iprange.shouldCombine(newRange) {
			if IPGreater(newRange.End, iprange.End) {
				iprange.End = newRange.End
				for i++; i < len(list.Blacklist); {
					if iprange.shouldCombine(list.Blacklist[i]) {
						if IPLess(iprange.End, list.Blacklist[i].End) {
							iprange.End = list.Blacklist[i].End
							list.Blacklist = remove(list.Blacklist, i)
							break
						}
						list.Blacklist = remove(list.Blacklist, i)
					} else {
						break
					}
				}
			}
			return
		}
		if newRange.shouldCombine(iprange) {
			iprange.Start = newRange.Start
			if IPGreater(newRange.End, iprange.End) {
				iprange.End = newRange.End
				for i++; i < len(list.Blacklist); {
					if iprange.shouldCombine(list.Blacklist[i]) {
						if IPLess(iprange.End, list.Blacklist[i].End) {
							iprange.End = list.Blacklist[i].End
							list.Blacklist = remove(list.Blacklist, i)
							break
						}
						list.Blacklist = remove(list.Blacklist, i)
					} else {
						break
					}
				}
			}
			return
		}
	}
	//If we get this far we know that range being inserted is disjoint from every other range
	for i, iprange := range list.Blacklist {
		if IPLess(newRange.End, iprange.Start) {
			temp := make([]*IPRange, len(list.Blacklist[i:]))
			copy(temp, list.Blacklist[i:])
			list.Blacklist = append(append(list.Blacklist[:i], newRange), temp...)
			return
		}
	}
	list.Blacklist = append(list.Blacklist, newRange)
}

