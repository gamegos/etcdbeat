# Etcdbeat

Etcdbeat is an [Elastic Beat](https://www.elastic.co/products/beats) that reads stats from the Etcd v2 API and indexes them into Elasticsearch or Logstash.   

## Description

> [Etcd](https://coreos.com/etcd) is an open-source distributed key value store that provides a reliable way to store data across a cluster of machines.  

An etcd cluster keeps track of a number of statistics including latency, bandwidth and uptime. Etcdbeat collects this stats with etcd API.  
  
### Usage
  
Check out [this blog post](http://www.berfinsari.com/2017/08/etcdbeat-elastic-beat-for-etcd_7.html) on how it works.
  
## Exported Fields

There are three types of documents exported:

- `type: leader` Contains etcd leader stats  
- `type: self` Contains etcd self stats  
- `type: store` Contains etcd store stats  

**Leader Stats**

The leader has a view of the entire cluster and keeps track of two interesting statistics: latency to each peer in the cluster, and the number of failed and successful Raft RPC requests.  

**Self Stats**

Each node keeps a number of internal statistics:  
  
* `id`: the unique identifier for the member  
* `leaderInfo.leader`: id of the current leader member  
* `leaderInfo.uptime`: amount of time the leader has been leader  
* `name`: this member's name  
* `recvAppendRequestCnt`: number of append requests this node has processed  
* `recvBandwidthRate`: number of bytes per second this node is receiving (follower only)  
* `recvPkgRate`: number of requests per second this node is receiving (follower only)  
* `sendAppendRequestCnt`: number of requests that this node has sent  
* `sendBandwidthRate`: number of bytes per second this node is sending (leader only). This value is undefined on single member clusters.  
* `sendPkgRate`: number of requests per second this node is sending (leader only). This value is undefined on single member clusters.  
* `state`: either leader or follower  
* `startTime`: the time when this node was started  

Document example:  
  
<pre>  
{  
    "self": {  
        "id": "8e9e05c52164694d",  
        "leaderInfo": {  
            "leader": "8e9e05c52164694d",  
            "startTime": "2017-08-07T12:24:14.306354646+03:00",  
            "uptime": "8m52.578056951s"  
        },  
        "name": "node2",  
        "recvAppendRequestCnt": 0,  
        "recvBandwidthRate": 6345,  
        "recvPkgRate": 824.1758351191694,  
        "sendAppendRequestCnt": 11.111234716807138,  
        "startTime": "2017-08-07T12:24:13.805809934+03:00",  
        "state": "StateLeader"  
    }  
}  
</pre>    
  
**Store Stats**

The store statistics include information about the operations that this node has handled.Operations that modify the store's state like create, delete, set and update are seen by the entire cluster and the number will increase on all nodes.   
  
Document example:    
  
<pre>  
{  
    "store": {  
        "compareAndDeleteFail": 0,  
        "compareAndDeleteSuccess": 0,  
        "compareAndSwapFail": 0,  
        "compareAndSwapSuccess": 0,  
        "createFail": 0,  
        "createSuccess": 2,  
        "deleteFail": 0,  
        "deleteSuccess": 0,  
        "expireCount": 0,  
        "getsFail": 4,  
        "getsSuccess": 42,  
        "setsFail": 3,  
        "setsSuccess": 6,  
        "updateFail": 0,  
        "updateSuccess": 0,  
        "watchers": 0  
    }  
}  
</pre>  

## Configuration  

Adjust the `etcdbeat.yml` configuration file to your needs.    

### `period`  
Defines how often to read statistics. Default to `30` s.	  

### `port`  
Defines the etcd port serviced. Default to `2379`  

### `host`  
Host name of ElasticSearch. Default to `localhost`    

### `statistics`  
You can decide which statistics to collect.    

```  
statistics:  
  leader: false  
  self: true  
  store: true  
```  

### `authentication`  
Authentication to be used for API connection. Default to `enable: false`    

```  
authentication:  
  enable: false  
  username: test  
  password: test1234  
```  

## Elasticsearch template  

The default template is provided, if you add any queries you should update the template accordingly.   
To apply default Etcdbeat template run:  

```
curl -XPUT 'http://<host>:9200/_template/etcdbeat' -d@etcdbeat.template.json  
```

## Getting Started with Etcdbeat

Ensure that this folder is at the following location:  
`${GOPATH}/github.com/gamegos`


### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Etcdbeat and also install the
dependencies, run the following command:  

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Etcdbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/gamegos/etcdbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Etcdbeat run the command below. This will generate a binary
in the same directory with the name etcdbeat.

```
make
```


### Run

To run Etcdbeat with debugging output enabled, run:

```
./etcdbeat -c etcdbeat.yml -e -d "*"
```


### Test

To test Etcdbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/etcdbeat.template.json and etc/etcdbeat.asciidoc

```
make update
```


### Cleanup

To clean  Etcdbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Etcdbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/gamegos
cd ${GOPATH}/github.com/gamegos
git clone https://github.com/gamegos/etcdbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.

## Author
Berfin Sarı <berfinsari21 'at' gmail.com>

## License
Covered under the Apache License, Version 2.0
Copyright (c) 2017 Berfin Sarı

