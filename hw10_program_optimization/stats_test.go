// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("cyrillic domains", func(t *testing.T) {
		data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"вася@москва.рф","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
		{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"маша@москва.рф","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
		{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
		{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
		{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"контакт@кремль.рф","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

		result, err := GetDomainStat(bytes.NewBufferString(data), "рф")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"москва.рф": 2,
			"кремль.рф": 1,
		}, result)
	})

	t.Run("only email field added to stat", func(t *testing.T) {
		data := `{"Id":1,"Name":"aliquid_qui_ea@test.com","Username":"aliquid_qui_ea@test.com","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"aliquid_qui_ea@test.com","Password":"aliquid_qui_ea@test.com","Address":"aliquid_qui_ea@test.com"}
		{"Id":2,"Name":"mLynch@test.com","Username":"mLynch@test.com","Email":"mLynch@broWsecat.com","Phone":"mLynch@test.com","Password":"mLynch@test.com","Address":"mLynch@test.com"}
		{"Id":3,"Name":"RoseSmith@test.com","Username":"RoseSmith@test.com","Email":"RoseSmith@Browsecat.com","Phone":"RoseSmith@test.com","Password":"RoseSmith@test.com","Address":"RoseSmith@test.com"}
		{"Id":4,"Name":"5Moore@test.com","Username":"5Moore@test.com","Email":"5Moore@Teklist.net","Phone":"5Moore@test.com","Password":"5Moore@test.com","Address":"5Moore@test.com"}
		{"Id":5,"Name":"nulla@test.com","Username":"nulla@test.com","Email":"nulla@Linktype.com","Phone":"nulla@test.com","Password":"nulla@test.com","Address":"nulla@test.com"}`

		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})
}
