package go3scale
import "fmt"
import "strconv"
// import "net/http"
import "github.com/parnurzeal/gorequest"

type Client struct {
	Provider_key, Host string
}

func New(provider_key string, ) (client *Client) {
    return &Client{
        Provider_key: provider_key,
				Host: "su1.3scale.net",
    }
}

func (client *Client) Print() {
    fmt.Println("Client provider_key:", client.Provider_key)
}

type Usage struct {
    Name string
    Value int
}

func (client *Client) Authrep_with_user_key(user_key string, usage Usage ) Response {
	r:= new(Response)

	url := "/transactions/authorize.xml?"
  query := "&usage[" + usage.Name + "]="+strconv.Itoa(usage.Value)
	query += "&user_key=" + user_key
	query += "&provider_key="+ client.Provider_key

	request := gorequest.New()
	res, body, errs := request.Get("https://"+client.Host+url+query).
	  Set("X-3scale-User-Agent", "plugin-golang-v#test").
	  End()

	if errs != nil {
		// handle error
	}
	if res.StatusCode == 200 ||  res.StatusCode == 409{
			fmt.Printf("\n"+"All good")
			r.Succeed()
	}else if res.StatusCode >= 400 && res.StatusCode < 409{
		fmt.Printf("\n"+"Error")
	}else{
		fmt.Printf("\n"+"Mega error")
	}
		fmt.Printf("\n"+body)
	return *r
}

/****** Response object ******/

type Response struct {
  error_code int
  error_message string
}

func (r Response) Succeed(){
  r.error_code = 0
  r.error_message = ""
}

func (r Response) Error(code int, message string){
  r.error_code = code
  r.error_message = message
}

func (r Response) IsSuccess() bool{
	return (r.error_code == 0 && r.error_message == "")
}
