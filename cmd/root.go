package cmd

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

//ASN ...
type ASN struct {
	Name    string
	Numbers []string
}

//CnASN ...
type CnASN struct {
	Asn []*ASN
}

//Appendumbers ...
func (a *ASN) Appendumbers(numbers []string) {
	length := len(numbers)
	for index, key := range numbers {
		if index != (length - 1) {
			a.Numbers = append(a.Numbers, key+",")
		} else {
			a.Numbers = append(a.Numbers, key)
		}

	}

}

//CreateAsnFile ...
func CreateAsnFile(output string, asn *CnASN) error {
	filename := path.Join(output, "cn.conf")
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	t := template.New("asn")
	t = template.Must(t.Parse(asnTempalte))
	return t.Execute(file, asn)
}

const asnTempalte = `
{{ range .Asn }}
define {{ .Name}} = [
	{{range .Numbers }}
		{{ . }}
	{{end}}
];
{{ end }}
`

var (
	output  string
	rootCmd = &cobra.Command{
		Use:   "china-asn",
		Short: "用于获取中国运营商asn号",
		Run: func(cmd *cobra.Command, args []string) {
			ctKeyWord := []string{"CHINANET", "CHINATELCOM", "ChinaTelecom"}
			cuKeyWord := []string{"UNICOM"}
			cmKeyWord := []string{"CMNET", "CTTNET", "TieTong"}
			creKeyWord := []string{"CERNET"}
			ctASN := make([]string, 0)
			cuASN := make([]string, 0)
			cmASN := make([]string, 0)
			creASN := make([]string, 0)
			otherASN := make([]string, 0)
			var row []string
			var rows [][]string
			res, err := http.Get("https://whois.ipip.net/countries/CN")
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
				tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
					rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
						row = append(row, tablecell.Text())
					})
					if len(row) == 4 {
						found := false
						str := strings.Replace(row[0], " ", "", -1)
						number := strings.Split(str, "AS")
						if len(number) != 2 {
							return
						}
						rows = append(rows, row)
						for _, k := range ctKeyWord {
							if strings.Contains(row[1], k) {
								ctASN = append(ctASN, number[1])
								found = true
								break
							}
						}
						for _, k := range cuKeyWord {
							if strings.Contains(row[1], k) {
								cuASN = append(cuASN, number[1])
								found = true
								break
							}
						}
						for _, k := range cmKeyWord {
							if strings.Contains(row[1], k) {
								cmASN = append(cmASN, number[1])
								found = true
								break
							}
						}
						for _, k := range creKeyWord {
							if strings.Contains(row[1], k) {
								creASN = append(creASN, number[1])
								found = true
								break
							}
						}
						if !found {
							otherASN = append(otherASN, number[1])

						}
					}
					row = nil
				})
			})
			cnAsn := &CnASN{}
			cnAsn.Asn = make([]*ASN, 0)
			cuAsn := &ASN{}
			cuAsn.Name = "cu_asn"
			cuAsn.Numbers = make([]string, 0)
			cuAsn.Appendumbers(cuASN)
			if err != nil {
				log.Fatal(err)
			}
			cnAsn.Asn = append(cnAsn.Asn, cuAsn)
			ctAsn := &ASN{}
			ctAsn.Name = "ct_asn"
			ctAsn.Numbers = make([]string, 0)
			ctAsn.Appendumbers(ctASN)
			cnAsn.Asn = append(cnAsn.Asn, ctAsn)
			cmAsn := &ASN{}
			cmAsn.Name = "cm_asn"
			cmAsn.Numbers = make([]string, 0)
			cmAsn.Appendumbers(cmASN)
			cnAsn.Asn = append(cnAsn.Asn, cmAsn)
			oAsn := &ASN{}
			oAsn.Name = "other_asn"
			oAsn.Numbers = make([]string, 0)
			oAsn.Appendumbers(otherASN)
			cnAsn.Asn = append(cnAsn.Asn, oAsn)
			err = CreateAsnFile(output, cnAsn)
			if err != nil {
				log.Fatal(err)

			}

		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "/etc/", "存储asn文件路径")
}
