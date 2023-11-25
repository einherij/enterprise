package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FuncName(t *testing.T) {
	a := assert.New(t)

	name := FuncName()
	a.Equal("test_func_name", name)
}

func Test_JoinQuoted(t *testing.T) {
	a := assert.New(t)

	quoted := JoinQuoted([]string{"three", "test", "lines"}, "'", ", ")
	a.Equal(`'three', 'test', 'lines'`, quoted)

	quoted = JoinQuoted([]string{"three", "test", "lines"}, "", "++")
	a.Equal(`three++test++lines`, quoted)
}

func Test_JoinParams(t *testing.T) {
	a := assert.New(t)

	combined := JoinParams("sspID", "countryCode")
	a.Equal("sspID#countryCode", combined)
	combined = JoinParams("", "countryCode")
	a.Equal("#countryCode", combined)
}

func TestGetFirst(t *testing.T) {
	a := assert.New(t)

	first := GetFirst("sspID#countryCode")
	a.Equal("sspID", first)
	first = GetFirst("")
	a.Equal("", first)
	first = GetFirst("#")
	a.Equal("", first)
}

func TestGetSecond(t *testing.T) {
	a := assert.New(t)

	second := GetSecond("sspID#countryCode")
	a.Equal("countryCode", second)
	second = GetSecond("")
	a.Equal("", second)
	second = GetSecond("#")
	a.Equal("", second)
}

func TestToUnderscore(t *testing.T) {
	a := assert.New(t)

	a.Equal(
		"reach_effect",
		ToUnderscoreACII("ReachEffect"))
	a.Equal(
		"tubecorp_bidmagnet_mgid_runative_ms_apac_tubecorporate_ms_evadav_adult_runative_inmobi_galaksion_topnews_vn_app",
		ToUnderscoreACII(`
tubecorp
Bidmagnet
MGID
Runative_ms_apac
tubecorporate_ms
Evadav_adult
runative
Inmobi
galaksion
topnews_vn_app
`))
}
