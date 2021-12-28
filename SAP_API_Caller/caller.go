package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-accounting-document-reads/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library/logger"
	"golang.org/x/xerrors"
)

type SAPAPICaller struct {
	baseURL string
	apiKey  string
	log     *logger.Logger
}

func NewSAPAPICaller(baseUrl string, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL: baseUrl,
		apiKey:  GetApiKey(),
		log:     l,
	}
}

func (c *SAPAPICaller) AsyncGetAccountingDocument(companyCode, fiscalYear, accountingDocument string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "Item":
			func() {
				c.Item(companyCode, fiscalYear, accountingDocument)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) Item(companyCode, fiscalYear, accountingDocument string) {
	itemData, err := c.callAccountingDocumentSrvAPIRequirementItem("A_OperationalAcctgDocItemCube", companyCode, fiscalYear, accountingDocument)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(itemData)

}

func (c *SAPAPICaller) callAccountingDocumentSrvAPIRequirementItem(api, companyCode, fiscalYear, accountingDocument string) ([]sap_api_output_formatter.Item, error) {
	url := strings.Join([]string{c.baseURL, "API_OPLACCTGDOCITEMCUBE_SRV", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithItem(req, companyCode, fiscalYear, accountingDocument)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToItem(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithItem(req *http.Request, companyCode, fiscalYear, accountingDocument string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("CompanyCode eq '%s' and FiscalYear eq '%s' and AccountingDocument eq '%s'", companyCode, fiscalYear, accountingDocument))
	req.URL.RawQuery = params.Encode()
}
