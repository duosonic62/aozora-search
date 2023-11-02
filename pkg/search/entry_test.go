package search

import (
	"fmt"
	"github.com/duosonic62/aozora-search/pkg/domain"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

func TestFindEntries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.String() == "/" {
			writer.Write([]byte(`
<table summary="作家データ">
<body>
<tr>
<td class="header">作家名：</td>
<td><font size="+2">テスト 太郎</font></td>
</tr>
<tr><td class="header">作家名読み：</td><td>てすと たろう</td></tr>
<tr><td class="header">ローマ字表記：</td><td>Test, Taro</td></tr>
<tr><td class="header">生年：</td><td>1892-03-01</td></tr>
<tr><td class="header">没年：</td><td>1927-07-24</td></tr>
</body>
</table>
<ol>
<li><a href="../cards/999999/card001.html">テスト書籍001</a></li>
<li><a href="../cards/999999/card002.html">テスト書籍002</a></li>
<li><a href="../cards/999999/card003.html">テスト書籍003</a></li>
</ol>
			`))
		} else {
			pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
			token := pat.FindStringSubmatch(request.URL.String())
			writer.Write([]byte(fmt.Sprintf(`
<table summary="作家データ">
<tbody><tr><td class="header">分類：</td><td>著者</td></tr>
<tr><td class="header">作家名：</td><td><a href="../../index_pages/person9999.html">テスト 太郎</a></td></tr>
<tr><td class="header">作家名読み：</td><td>テスト 太郎</td></tr>
<tr><td class="header">ローマ字表記：</td><td>Test, Taro</td></tr>
<tr><td class="header">生年：</td><td>1892-03-01</td></tr>
<tr><td class="header">没年：</td><td>1927-07-24</td></tr>
</tbody>
</table>
<table border="1" summary="ダウンロードデータ" class="download">
<tr>
<td><a href="./files/%[1]s_%[2]s.zip"></a>></td>
</tr>
</table>
`, token[1], token[2])))
		}
	}))
	defer ts.Close()

	tmp := pageUrlFormat
	pageUrlFormat = ts.URL + "/cards/%s/card%s.html"
	defer func() {
		pageUrlFormat = tmp
	}()

	got, err := FindEntries(ts.URL)
	if err != nil {
		t.Error(err)
	}
	want := []domain.Entry{
		{
			AuthorID: "999999",
			Author:   "テスト 太郎",
			TitleID:  "001",
			Title:    "テスト書籍001",
			InfoURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_001.zip",
		},
		{
			AuthorID: "999999",
			Author:   "テスト 太郎",
			TitleID:  "002",
			Title:    "テスト書籍002",
			InfoURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_002.zip",
		},
		{
			AuthorID: "999999",
			Author:   "テスト 太郎",
			TitleID:  "003",
			Title:    "テスト書籍003",
			InfoURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_003.zip",
		},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, but got %+v", want, got)
	}
}
