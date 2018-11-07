## Synopsis

If your team is using Azure services like Bots that calls some APIs from your server, 
you probably want to whitelist specific Azure IP addressess for security purposes.

In order to do that, you can go to https://www.microsoft.com/en-gb/download/details.aspx?id=41653 
and download the XML file from there. The problem is, the IP ranges may change on a regular basis 
and if you are using infra as code tools such as Terraform, you need to find an easy way to automate 
this task and update your Terraform variables.

**azureip2tfvars** is created to easily extract the latest Azure Data Centre IPs and create a 
Terraform variables file out of it. We use Terraform to create AWS security groups with ingress rules that whitelist these Azure data centre IP's. 
**azureip2tfvars** lets our team update our `tf` files on a regular basis easily to keep our ingress rules up to date.

## Quickstart

**WIP DOWNLOAD GUIDE AND RELEASES**

Download the binary for your platform and just run it. No dependencies required!

```bash
./azureip2tfvars -writeto /tmp/vars.tf

⡿ Download url is: https://download.microsoft.com/download/0/1/8/018E208D-54F8-44CD-AA26-CD7BC9524A8C/PublicIPs_20181107.xml
Terraform file "/tmp/vars.tf" successfully created!
```

Example output

```terraform
...

variable "azure_indiawest_subnets" {
    type = "list"
    default = [
    "20.40.8.0/21",
    "20.190.146.128/25",
    "40.79.219.0/24",
    "40.81.80.0/20",
    "40.87.220.0/22",
    "40.90.138.224/27",
    "40.126.18.128/25",
    "52.109.64.0/22",
    "52.114.28.0/22",
    "52.136.16.0/24",
    "52.136.32.0/19",
    "52.140.128.0/18",
    "52.183.128.0/18",
    "52.239.135.192/26",
    "52.239.187.128/25",
    "52.245.76.0/22",
    "52.249.64.0/19",
    "104.44.93.224/27",
    "104.44.95.112/28",
    "104.47.212.0/23",
    "104.211.128.0/18",
    ]
}

variable "azure_uswest2_subnets" {
...
}
```
