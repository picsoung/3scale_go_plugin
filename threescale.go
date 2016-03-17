package go3scale

import "fmt"
import "errors"
import "strconv"
import "github.com/parnurzeal/gorequest"

// Client is identified by its ProviderKey and Host
type Client struct {
	ProviderKey, Host string
}

func clientParams(args ...interface{}) (providerKey string, host string, err error){

  host = "su1.3scale.net"
  if 1 > len(args) {
      err = errors.New("Not enough parameters to init Client")
      return
  }

  for i,p := range args {
    switch i {
	    case 0: // name
        param, ok := p.(string)
        if !ok {
            err = errors.New("ProviderKey parameter not type string.")
            return
        }
        providerKey = param

	    case 1: // x
        param, ok := p.(string)
        if !ok {
            err = errors.New("Host parameter not type string.")
            return
        }
        host = param
	    default:
        err = errors.New("Too many parameters")
        return
    }
  }
  return
}

//New Creates a new Client
func New(args ...interface{}) (client *Client) {
	providerKey, host, err := clientParams(args...)
	if nil != err {
    panic(err.Error())
  }

	return &Client{
		ProviderKey: providerKey,
		Host:        host,
	}
}

// Usage is a name and a value
type Usage struct {
	Name  string
	Value int
}

func authrepUserKeyParams(args ...interface{}) (userKey string, usageArr []Usage, err error){
	var arr []Usage
	arr = append(arr,Usage{Name:"hits",Value:1})
	usageArr = arr

  if 1 > len(args) {
      err = errors.New("Not enough parameters for AuthrepUserKey")
      return
  }

  for i,p := range args {
    switch i {
	    case 0: // userKey
        param, ok := p.(string)
        if !ok {
            err = errors.New("userKey parameter not type string.")
            return
        }
        userKey = param

			case 1: // usage
        param, ok := p.([]Usage)
        if !ok {
            err = errors.New("usage parameter not type Usage.")
            return
        }
        usageArr = param
	    default:
        err = errors.New("Too many parameters")
        return
    }
  }
  return
}

//AuthrepUserKey authenticates a request with userKey
func (client *Client) AuthrepUserKey(args ...interface{}) Response {
	userKey, usageArr, err := authrepUserKeyParams(args...)
	if nil != err {
    panic(err.Error())
  }

	r := new(Response)
	url := "/transactions/authrep.xml?"
	query := "user_key=" + userKey
	query += "&provider_key=" + client.ProviderKey
	for _,element := range usageArr {
		query += "&usage[" + element.Name + "]=" + strconv.Itoa(element.Value)
	}

	// fmt.Printf("\n"+"https://"+client.Host+url+query+"\n")
	request := gorequest.New()
	res, body, errs := request.Get("https://"+client.Host+url+query).
		Set("X-3scale-User-Agent", "plugin-golang-v#test").
		End()

	if errs != nil {
		// handle error
	}
	if res.StatusCode == 200 || res.StatusCode == 409 {
		fmt.Printf("\n" + "All good")
		r.Succeed()
	} else if res.StatusCode >= 400 && res.StatusCode < 409 {
		fmt.Printf("\n" + "Error"+"\n",body)
		r.Error(res.StatusCode,"")
	} else {
		fmt.Printf("\n" + "Mega error")
	}
	fmt.Printf("\n" + body)
	return *r
}

//Response object
type Response struct {
	errorCode    int
	errorMessage string
}

//Succeed response
func (r Response) Succeed() {
	r.errorCode = 0
	r.errorMessage = ""
}

//Error response
func (r Response) Error(code int, message string) {
	r.errorCode = code
	r.errorMessage = message
}

//IsSuccess checks is response succeed
func (r Response) IsSuccess() bool {
	return (r.errorCode == 0 && r.errorMessage == "")
}
