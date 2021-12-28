package main

import (
	sap_api_caller "sap-api-integrations-accounting-document-reads/SAP_API_Caller"
	"sap-api-integrations-accounting-document-reads/sap_api_input_reader"

	"github.com/latonaio/golang-logging-library/logger"
)

func main() {
	l := logger.NewLogger()
	fr := sap_api_input_reader.NewFileReader()
	inoutSDC := fr.ReadSDC("./Inputs/SDC_Accounting_Document_Item_sample.json")
	caller := sap_api_caller.NewSAPAPICaller(
		"https://sandbox.api.sap.com/s4hanacloud/sap/opu/odata/sap/", l,
	)

	accepter := inoutSDC.Accepter
	if len(accepter) == 0 || accepter[0] == "All" {
		accepter = []string{
			"Item",
		}
	}

	caller.AsyncGetAccountingDocument(
		inoutSDC.AccountingDocument.CompanyCode,
		inoutSDC.AccountingDocument.FiscalYear,
		inoutSDC.AccountingDocument.AccountingDocument,
		accepter,
	)
}
