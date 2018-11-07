package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/briandowns/spinner"
)

var (
	filepath *string
)

const azureDownloadLink = "https://www.microsoft.com/en-gb/download/confirmation.aspx?id=41653"

func main() {
	filepath = flag.String("writeto", "./azure_tfvars.tf", "The terraform file to write to. (default: ./azure_tfvars.tf)")
	flag.Parse()
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("red", "bold")
	s.Start()
	if err := downloadAndParseAzureXML(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.Stop()
	fmt.Printf("Terraform file %q successfully created!\n", *filepath)
}

func getDownloadURL() (string, error) {
	resp, err := http.Get(azureDownloadLink)
	if err != nil {
		return "", fmt.Errorf("Error: failed fetching %s: %+v", azureDownloadLink, err)
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error: failed reading rendered html from %s: %+v", azureDownloadLink, err)
	}
	regexpString := `url=https:\/\/download\.microsoft\.com\/download.+xml`
	r, err := regexp.Compile(regexpString)
	if err != nil {
		return "", fmt.Errorf("Error: failed compiling regex: %+v", err)
	}
	if ok := r.MatchString(string(bs)); !ok {
		return "", fmt.Errorf("Error: regexp string %q has no match", regexpString)
	}
	return strings.Split(r.FindString(string(bs)), "=")[1], nil
}

func downloadAndParseAzureXML() error {
	downloadURL, err := getDownloadURL()
	if err != nil {
		return err
	}
	fmt.Println("Download url is:", downloadURL)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("Error: failed downloading file from %s: %+v", downloadURL, err)
	}
	// Read the downloaded XML file
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error: failed reading response body after downloading: %+v", err)
	}
	// Parse to AzureIP struct
	ipList := &AzureIP{}
	if err := xml.Unmarshal(bs, ipList); err != nil {
		return fmt.Errorf("Error: unmarshalling XML to AzureIP struct: %+v", err)
	}
	if err := ipList.createTfFile(); err != nil {
		return err
	}
	return nil
}

func (azip *AzureIP) createTfFile() error {
	t := template.Must(template.New("azip").Parse(tftempl))
	f, err := os.Create(*filepath)
	if err != nil {
		return fmt.Errorf("Error: failed creating file: %+v", err)
	}
	if err := t.Execute(f, azip.Regions); err != nil {
		return fmt.Errorf("Error: failed templating azip: %+v", err)
	}
	return nil
}

const tftempl = `
{{- range . }}
variable "azure_{{.Name}}_subnets" {
    type = "list"
    default = [
    {{- range .IPRanges }}
    "{{ .Subnet }}",
    {{- end}}
    ]
}
{{end}}
`

// AzureIP holds the struct for parsing the following XML format
//
//  <?xml version="1.0" encoding="utf-8"?>
//  <AzurePublicIpAddresses xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
//    <Region Name="australiac2">
//      <IpRange Subnet="20.36.64.0/19" />
//      <IpRange Subnet="20.36.112.0/20" />
//      <IpRange Subnet="20.39.72.0/21" />
//      <IpRange Subnet="20.39.96.0/19" />
//      <IpRange Subnet="40.82.12.0/22" />
//      <IpRange Subnet="40.82.244.0/22" />
//      <IpRange Subnet="40.90.130.32/28" />
//      <IpRange Subnet="40.90.142.64/27" />
//      <IpRange Subnet="40.90.149.32/27" />
//      <IpRange Subnet="40.126.128.0/18" />
//      <IpRange Subnet="52.143.218.0/24" />
//      <IpRange Subnet="52.239.218.0/23" />
//    </Region>
//  </AzurePublicIpAddresses>
//
type AzureIP struct {
	XMLName xml.Name `xml:"AzurePublicIpAddresses"`
	Regions []Region `xml:"Region"`
}

// Region is the Region field inside AzurePublicIpAddresess
type Region struct {
	Name     string    `xml:"Name,attr"`
	IPRanges []IPRange `xml:"IpRange"`
}

// IPRange is the IpRange field inside Region
type IPRange struct {
	Subnet string `xml:"Subnet,attr"`
}
