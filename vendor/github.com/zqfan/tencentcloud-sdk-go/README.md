# tencentcloud-sdk-go

## Dedicated API example

```
package main

import (
    "github.com/zqfan/tencentcloud-sdk-go/common"
    cvm "github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
    "fmt"
)

func main() {
    client, _ := cvm.NewClientWithSecretId("YOUR_SECRET_ID", "YOUR_SECRET_KEY", "REGION_NAME")
    request := cvm.NewDescribeAddressesRequest()
    request.Limit = common.IntPtr(10)
    // get response structure
    response, err := client.DescribeAddresses(request)
    // API errors
    if _, ok := err.(*common.APIError); ok {
        fmt.Printf("An API error has returned: %s", err)
        return
    }
    // unexpected errors
    if err != nil {
        panic(err)
    }
    eips := response.Response.AddressSet
    fmt.Println(common.StringValues(eips))
}
```

## Common API Example

```
package main

import (
    "github.com/zqfan/tencentcloud-sdk-go/client"
)

func main() {
    client := client.NewClient("YOUR_SECRET_ID", "YOUR_SECRET_KEY", "REGION_NAME")
    params := map[string]string {
        "Action": "DescribeInstances",
        "SignatureMethod": "HmacSHA256",
    }
    // get raw text response body from server
    response, err := client.SendRequest("cvm", params)
}
```

If you want to print debug info, set ``Client`` object's attribute ``Debug`` to true like: ``client.Debug = true``
