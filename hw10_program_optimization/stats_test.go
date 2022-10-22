// +build !bench

package hw10programoptimization

import (
	"archive/zip"
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
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
}

//go test -bench=BenchmarkGetDomainStat -benchmem -count 1
func BenchmarkGetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, _ := zip.OpenReader("testdata/users.dat.zip")
		data, _ := r.File[0].Open()
		GetDomainStat(data, "biz")
	}
}

func TestGetDomainStatIncorrectData(t *testing.T) {
	t.Run("Empty data", func(t *testing.T) {
		data := ""
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
	t.Run("Plain text", func(t *testing.T) {
		data := "plain text"
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.Contains(t, err.Error(), "get users error:")
		require.Equal(t, "get users error: parse error: syntax error near offset 0 of 'plain text'", err.Error())
		require.Nil(t, result)
	})

	t.Run("Incorrect email", func(t *testing.T) {
		data := `{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
	t.Run("Incorrect json", func(t *testing.T) {
		data := `{"Id":5,"Name":"ddress":"Russell Trail 61"}`
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.Contains(t, err.Error(), "parse error: syntax error")
		require.Nil(t, result)
	})
}
