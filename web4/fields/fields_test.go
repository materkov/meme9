package fields

import "testing"

func TestParseFields(t *testing.T) {
	//ParseFields("foo,test(user,post(name))")
	ParseFields("foo,bar(user)")
}
