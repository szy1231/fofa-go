package fofa_test

import (
	"github.com/buger/jsonparser"
	"github.com/xiaoyu-0814/fofa-go/fofa"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

var (
	clt = fofa.NewFofaClient([]byte(os.Getenv("FOFA_EMAIL")), []byte(os.Getenv("FOFA_KEY")))
)

func EqualBytes(a, b []byte) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNewFofaClient(t *testing.T) {
	rand.Seed(time.Now().Unix() ^ 0x1a2b3c4d)
	for i := 0; i < 100; i++ {
		email := strconv.Itoa(rand.Int())
		key := strconv.Itoa(rand.Int())
		clt := fofa.NewFofaClient([]byte(email), []byte(key))
		if !EqualBytes([]byte(email), clt.Email) || !EqualBytes([]byte(key), clt.Key) {
			t.Errorf("expect email = %s  key = %s , but email = %s  key = %s\n", email, key, clt.Email, clt.Key)
		}
	}
}

func TestNewFofaClientError(t *testing.T) {
	email := os.Getenv("FOFA_EMAIL")
	key := os.Getenv("FOFA_KEY")
	clt := fofa.NewFofaClient([]byte(email+"0000"), []byte(key))
	userinfo, err := clt.UserInfo()
	if err == nil {
		t.Errorf("%v\n", err.Error())
	} else if userinfo != nil {
		t.Errorf("expect userinfo is empty, but %v\n", userinfo)
	}
	clt = fofa.NewFofaClient([]byte(email), []byte(key+"0000"))
	userinfo, err = clt.UserInfo()
	if err == nil {
		t.Errorf("%v\n", err.Error())
	} else if userinfo != nil {
		t.Errorf("expect userinfo is empty, but %v\n", userinfo)
	}
}

func TestFofa_UserInfo(t *testing.T) {
	email := os.Getenv("FOFA_EMAIL")
	key := os.Getenv("FOFA_KEY")

	clt := fofa.NewFofaClient([]byte(email), []byte(key))
	info, err := clt.UserInfo()
	if err != nil {
		t.Fatalf("Failed get userInfo, got error: %v", err)
	}
	t.Log(info)
}

func TestQueryAsJSON(t *testing.T) {
	var (
		arr          = []byte(nil)
		err          = error(nil)
		modeNormal   = []byte("normal")
		modeExtended = []byte("extended")
		query        = []byte(nil)
		fields       = []byte(nil)
		page         = uint(0)
	)
	// -------------------------------------------
	clt := fofa.NewFofaClient([]byte(os.Getenv("FOFA_EMAIL")), []byte(os.Getenv("FOFA_KEY")))
	if clt == nil {
		t.Errorf("create fofa client failed!")
	}
	{
		{

			{ // ------extended-----domain-----------------------------

				query = []byte(`host="nosec.org"`)
				fields = []byte(`fields=domain`)
				page = 1
				arr, err = clt.QueryAsJSON(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					jsonExpectEqual(t, arr, modeExtended, query, page)
				}
			} // -------------------------------------------

			{ // -------extended-----host-----------------------------
				query = []byte(`host="nosec.org"`)
				fields = []byte(`fields=host`)
				page = 1
				arr, err = clt.QueryAsJSON(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					jsonExpectEqual(t, arr, modeExtended, query, page)
				}
			} // -------------------------------------------

			{ // -------normal------domain------------------------------
				query = []byte(`"nosec.org"`)
				fields = []byte(`fields=domain`)
				page = 1
				arr, err = clt.QueryAsJSON(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					jsonExpectEqual(t, arr, modeNormal, []byte(`"nosec.org"`), page)
				}
			} // -------------------------------------------

			{ // --------normal------host-----------------------------
				query = []byte(`"nosec.org"`)
				fields = []byte(`fields=host`)
				page = 1
				arr, err = clt.QueryAsJSON(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					jsonExpectEqual(t, arr, modeNormal, query, page)
				}
			} // -------------------------------------------
		} // -------------------------------------------

		{ // -------------------------------------------
			query = []byte(`host="nosec.org"`)
			fields = nil
			page = 1
			arr, err = clt.QueryAsJSON(page, query)
			if err != nil {
				t.Errorf("%v\n", err.Error())
			} else {
				jsonExpectEqual(t, arr, modeExtended, query, page)
			}
		} // -------------------------------------------
	} // -------------------------------------------
}

func TestQueryAsStruct(t *testing.T) {
	var (
		data         = fofa.Data{}
		err          = error(nil)
		modeNormal   = "normal"
		modeExtended = "extended"
		query        = []byte(nil)
		fields       = []byte(nil)
		page         = uint(0)
	)
	// -------------------------------------------
	clt := fofa.NewFofaClient([]byte(os.Getenv("FOFA_EMAIL")), []byte(os.Getenv("FOFA_KEY")))
	if clt == nil {
		t.Errorf("create fofa client failed!")
	}
	{
		{

			{ // ------extended-----domain-----------------------------

				query = []byte(`host="nosec.org"`)
				fields = []byte(`fields=domain`)
				page = 1
				data, err = clt.QueryAsObject(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					dataExpectEqual(t, data, modeExtended, string(query), page)
				}
				t.Log(data)
			} // -------------------------------------------

			{ // -------extended-----host-----------------------------
				query = []byte(`host="nosec.org"`)
				fields = []byte(`fields=host`)
				page = 1
				data, err = clt.QueryAsObject(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					dataExpectEqual(t, data, modeExtended, string(query), page)
				}
				t.Log(data)
			} // -------------------------------------------

			{ // -------normal------domain------------------------------
				query = []byte(`"nosec.org"`)
				fields = []byte(`fields=domain`)
				page = 1
				data, err = clt.QueryAsObject(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					dataExpectEqual(t, data, modeNormal, `"nosec.org"`, page)
				}
				t.Log(data)
			} // -------------------------------------------

			{ // --------normal------host-----------------------------
				query = []byte(`"nosec.org"`)
				fields = []byte(`fields=host`)
				page = 1
				data, err = clt.QueryAsObject(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					dataExpectEqual(t, data, modeNormal, string(query), page)
				}
				t.Log(data)
			} // -------------------------------------------
		} // -------------------------------------------

		{ // -------------------------------------------
			query = []byte(`host="nosec.org"`)
			fields = nil
			page = 1
			data, err = clt.QueryAsObject(page, query)
			if err != nil {
				t.Errorf("%v\n", err.Error())
			} else {
				dataExpectEqual(t, data, modeExtended, string(query), page)
			}
			t.Log(data)
		} // -------------------------------------------
	} // -------------------------------------------
}

func TestIP(t *testing.T) {
	if clt == nil {
		t.Fatalf("create fofa client failed!")
	}

	var (
		query, fields = []byte(nil), []byte(nil)
		data          = fofa.Data{}
		err           = error(nil)
	)

	query = []byte(`ip="106.75.75.203"`)
	fields = []byte(`host,domain,ip,port,title,city,country`)
	data, err = clt.QueryAsObject(1, query, fields)
	switch {
	case err != nil:
		t.Errorf("%v\n", err.Error())
	case len(data.Results) != 4:
		t.Errorf("expect 4 records, but get %d\n", len(data.Results))
	}
	for _, v := range data.Results {
		switch {
		case v[3] == "6379":
			if v[1] != "" || v[4] != "" {
				t.Errorf("%s\n", v)
			}
		case v[3] == "443":
			switch v[4] {
			case "FOFA Pro - 网络空间安全搜索引擎":
				if v[1] != "fofa.so" && v[1] != "106.75.75.204" {
					t.Errorf("%s\n", v)
				}
			}
		}
	}
}

func jsonExpectEqual(t *testing.T, data, mode, query []byte, page uint) {
	m := getMode(data)
	q := getQuery(data)
	p := getPage(data)
	if !EqualBytes(m, mode) || !EqualBytes(q, query) || page != p {
		_, f, r, _ := runtime.Caller(1)
		t.Errorf("%s %d: Expect\tmode=%s  query=%s  page=%d\nBut\tmode=%s  query=%s  page=%d ", f, r, mode, query, page, m, q, p)
	}
}

func dataExpectEqual(t *testing.T, data fofa.Data, mode, query string, page uint) {
	if data.Mode != mode || data.Query != query || data.Page != int(page) {
		_, f, r, _ := runtime.Caller(1)
		t.Errorf("%s %d: Expect\tmode=%s  query=%s  page=%d\nBut\tmode=%s  query=%s  page=%d ", f, r, mode, query, page, data.Mode, data.Query, data.Page)
	}
}
func getDomain(data []byte) []byte {
	domain, _ := jsonparser.GetString(data, "domain")
	return []byte(domain)
}
func getHost(data []byte) []byte {
	host, _ := jsonparser.GetString(data, "host")
	return []byte(host)
}
func getIP(data []byte) []byte {
	ip, _ := jsonparser.GetString(data, "ip")
	return []byte(ip)
}
func getPort(data []byte) []byte {
	port, _ := jsonparser.GetString(data, "port")
	return []byte(port)
}
func getCountry(data []byte) []byte {
	country, _ := jsonparser.GetString(data, "country")
	return []byte(country)
}
func getCity(data []byte) []byte {
	city, _ := jsonparser.GetString(data, "city")
	return []byte(city)
}
func getPage(data []byte) uint {
	page, _ := jsonparser.GetInt(data, "page")
	return uint(page)
}
func getMode(data []byte) []byte {
	mode, _ := jsonparser.GetString(data, "mode")
	return []byte(mode)
}
func getQuery(data []byte) []byte {
	query, _ := jsonparser.GetString(data, "query")
	return []byte(query)
}
func getSize(data []byte) uint {
	size, _ := jsonparser.GetInt(data, "size")
	return uint(size)
}
