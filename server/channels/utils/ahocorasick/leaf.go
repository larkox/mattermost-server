package ahocorasick

type Leaf struct {
	Term   string
	Values []LeafObject
}

type LeafObject struct {
	UserID      string
	MentionType int
}

func (l *Leaf) AddValue(v LeafObject) {
	for _, existingValue := range l.Values {
		if v.Equals(existingValue) {
			return
		}
	}
	l.Values = append(l.Values, v)
}

func (l *Leaf) RemoveValue(v LeafObject) {
	for i, existingValue := range l.Values {
		if v.Equals(existingValue) {
			l.Values = append(l.Values[:i], l.Values[i+1:]...)
			return
		}
	}
}

func (l Leaf) HasValues() bool {
	return len(l.Values) != 0
}

func (l LeafObject) Equals(other LeafObject) bool {
	return l.UserID == other.UserID && l.MentionType == other.MentionType
}
