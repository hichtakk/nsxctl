nsxctl
======

`nsxctl` is simple command-line client for NSX-T and NSX ALB. It helps operating both of NSX solutions with single and intuitive command from their local environment.

It supports not only execute implemented operational commands but also just call REST API easily. nsxctl helps authentication to target NSX endpoint and you can focus on only REST API method, URI and params.

## Usage
### Configure endpoints and credentials

Nsxctl can manage multiple NSX endpoints. To register endpoint, use `config` subcommand with name argument and parameters bellow.

| Parameter | Description |
| --- | --- |
| name | any name for managing the combination of endpoint and user credential|
| -e / --endpoint | NSX Manager or NSX ALB Controller endpoint |
| -u / --user | user name to use |
| -p / --password | password for specified user |
| --alb | (optional) Use this flag to configure NSX ALB site. NSX-T is default. |
| --init | (optional) Use this flag to create a new configuration file |

```
# example
nsxctl config set-site ${SITE_NAME} --endpoint https://${MANAGER-IP} --user ${USER} --password ${PASSWORD} --init
```

This configuration is stored `~/.config/nsxctl.json` by default. You can change the path and file name with `-c/--config` option. Sensitive data such as user credential will be base64 encoded.

Now you are ready to run `nsxctl`! see examples[]

## Examples
### show resources
show tier 0 gateways
```
nsxctl show gateway --tier 0
ID              Name            HA Mode         Failover Mode
tier0-01        tier0-01        ACTIVE_ACTIVE   NON_PREEMPTIVE
```
show routing table of specified tier 0 gateway on each edge node
```
nsxctl show routes tier0-01
/edge-cluster/ec01/node/edge01
B> 0.0.0.0/0 [20] via 10.111.41.50
C> 100.64.56.0/31 is directly connected
C> 100.64.96.0/31 is directly connected
C> 100.64.200.0/31 is directly connected
C> 100.64.208.0/31 is directly connected
C> 169.254.0.0/25 is directly connected
i> 169.254.0.128/25 [0] blackhole
C> 172.17.0.0/16 is directly connected
c> 172.21.0.0/28 [3] via 100.64.208.1
c> 172.21.0.16/28 [3] via 100.64.96.1
l> 172.21.3.1/32 [3] via 100.64.208.1
l> 172.21.3.2/32 [3] via 100.64.208.1
C> 172.21.15.0/24 is directly connected
```
show BGP advertised routes of specified tier 0 gateway
```
nsxctl show adv tier0-01
BGP neighbor: 10.111.41.50, Remote ASN: 100
Edge node: edge01, Source IP: 10.111.41.47

Network          Next Hop       Metric          Local Pref      Path 
10.111.41.0/26   0.0.0.0             0                 100      
172.17.0.0/16    0.0.0.0             0                 100      
172.21.0.0/28    0.0.0.0             0                 100      
172.21.0.16/28   0.0.0.0             0                 100      
172.21.3.1/32    0.0.0.0             0                 100      
172.21.3.2/32    0.0.0.0             0                 100      
172.21.15.0/24   0.0.0.0             0                 100      
```
show NSX ALB virtual services
```
nsxctl show alb-virtualservice
ID            Name         VIP             Port          Cloud           SEGroup         Status
9c77032320d5  http01       172.21.15.200   22            Default-Cloud   Default-Group   UP
2c8417fbc486  dns01        172.21.15.201   53            Default-Cloud   Default-Group   UP
541c092afd1e  system-mc01  172.21.15.10    6443          tkg             tkg             UP
e3447bf45756  system-mc02  172.21.15.12    6443          tkg             tkg             UP
79d377182366  Shared-L7-0  172.21.15.14    80,443(SSL)   tkg             tkg             UP
```

### monitor gateway traffic
monitor uplink traffic of specified tier 0 gateway with 5 seconds refresh interval
```
nsxctl top gateway --tier 0 tier0-01 --interval 5
[Press ESC or Ctrl-C to exit]

ID: tier0-01, Name: tier0-01
HA: ACTIVE_ACTIVE, Preempt: NON_PREEMPTIVE

IfName      TX [bps]  TX[pps]  RX [bps]  RX[pps]
────────────────────────────────────────────────
ext-edge01  0.00      0.00     3933.32   46.00
ext-edge02  86.47     1.00     3668.52   44.00
```

### call API
get NSX-T segments
```
nsxctl exec get /policy/api/v1/infra/segments
```
get NSX ALB virtual services
```
nsxctl exec get --alb /api/virtualservice
```
update NSX-T pool with json parameter (please refer NSX-T REST API specification regarding JSON parameters)
```
nsxctl exec patch /policy/api/v1/infra/ip-pools/pool01 -f ./pool.json
```

## Contributing
Thank you very much for taking the time to give feedback and improvement suggestion.
If you want to submit pull requests to fix bugs or any enhancements, please open an issue and link it to your pull request.
If you have any questions, feel free to open an issue.