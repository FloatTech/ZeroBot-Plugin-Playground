// Package github github工具类
package github

import "github.com/antchfx/htmlquery"

var (
	// Githubhost github域名
	Githubhost = "https://github.com"
)

// Githubfile github文件类
type Githubfile struct {
	Name string
	Type string
	Href string
}

// Files 取得某个目录的github文件
func Files(githubPath string) (fs []Githubfile, ds []Githubfile, err error) {
	doc, err := htmlquery.LoadURL(githubPath)
	if err != nil {
		return
	}
	list, err := htmlquery.QueryAll(doc, "//div[@aria-labelledby=\"files\"]/div[@role=\"row\"]")
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		k := htmlquery.FindOne(list[i], "//div[@role=\"rowheader\"]/span/a")
		if k != nil {
			f := Githubfile{}
			f.Name = htmlquery.InnerText(k)
			f.Href = k.Attr[len(k.Attr)-1].Val
			if k.Attr[1].Val == "This path skips through empty directories" {
				f.Type = "dir"
				ds = append(ds, f)
			} else {
				f.Type = "file"
				fs = append(fs, f)
			}
		}
	}
	return
}
