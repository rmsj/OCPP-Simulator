package simulator

import (
	"encoding/xml"
	"text/template"
	"log"
	"bytes"
	"fmt"
	"time"
)

type XMLStopTransactionResponse struct {
	XMLName xml.Name  `xml:"stopTransactionResponse"`
	IdTagInfo XMLIdTagInfo
}

type XMLStopTransactionBody struct {
	XMLName  xml.Name `xml:"Body"`
	StopTransactionResponse XMLStopTransactionResponse
}

type EnvelopeStopTransaction struct {
	XMLName  xml.Name    `xml:"Envelope"`
	Body   XMLStopTransactionBody
}

// Defines structure to render XML for Authorize request
// TODO: needs an array of TransactionData to render the template
type StopTransactionData struct {
	RequestData
	DateTimeStop string
}

type StopTransaction struct {
	Response *EnvelopeStopTransaction
}

// parses the XML - adding values to parameters, etc.
func (auth *StopTransaction) ParseRequestBody(data []string) string {

	// TODO: validate number of arguments

	var buffer bytes.Buffer
	tpl := template.Must(template.ParseFiles(auth.Template()))

	// date and time we are starting the transaction
	t1 := time.Now()

	// template data
	tplData := StopTransactionData{
		RequestData{
			data[1],
			data[2],
		},
		t1.Format(time.RFC3339),
	}

	fmt.Println("here");

	err := tpl.Execute(&buffer, tplData)
	if err != nil {
		log.Fatalln(err)
	}

	return buffer.String()
}

// Parse response
func (auth *StopTransaction) ParseResponseBody(responseData []byte) {
	err := xml.Unmarshal(responseData, &auth.Response)
	if err != nil {
		log.Fatalln(err)
	}

	//auth.Response = *response;
}

// Gets the XML to be used for this request
func (auth *StopTransaction) Template() string {
	return "xml/StopTransaction.xml"
}

// Gets the response status for the Authorize request
func (auth *StopTransaction) ResponseStatus() string {
	return auth.Response.Body.StopTransactionResponse.IdTagInfo.Status.Value
}

// Check if the authorize call to the central system has been accepted
func (auth *StopTransaction) Accepted() bool {
	return auth.ResponseStatus() == StatusAccepted
}

// Gets the response status for the Authorize request
func (auth *StopTransaction) ValidateArguments(data []string) {
	// TODO
}

func NewStopTransaction() *StopTransaction {
	return &StopTransaction{}
}