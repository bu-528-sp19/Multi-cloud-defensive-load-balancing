#                           Cloud Hydra: A Multi-Cloud Load Balancing and Failover framework

** **

## 1.   Vision and Goals Of The Project:

Hydra will be a framework for applications in the cloud to mitigate DDoS attacks and provider outages by providing resiliency at multiple levels both intra- and inter-clouds. If parts of AWS or GCP go down, for example, the application itself should be alive and kicking, as resources will be directed to the provider that is still up.

## 2. Users/Personas Of The Project:

1) Application Developers that want a resiliency framework

2) Users of the application that use Hydra that will experience no downtime

** **

## 3.   Scope and Features Of The Project:

- Request-level load balancing and queueing between hosts across clouds.
- If the compute layer of one cloud provider goes down, the application should still function.
- Distributed writes to all database solutions.
- Ability to read from any healthy database server.
- Eventual consistency of newly spawned databases via replication.
- If the database layer of one cloud provider goes down, the application should still function.

** **

## 4. Solution Concept

Hydra will provide resiliency at multiple levels.

### Compute

Hydra will feature a request server layer in each cloud that will receive incoming requests and round robin load balance them to hosts per Cloud. The Request Server layers in each Cloud will be heartbeating between eachother to ensure uptime and analyze load. In a two Request Server architecture, one will be marked as primary and will forward every other request to the secondary Request Server (e.g GCP forwards every other request to AWS). In the event that a secondary server goes down, the primary Request Server will remove it's eligibility in the round robin scheme until it can re-establish contact. If the primary Request Server goes down, the secondary servers will be notified via unresponsive heartbeat. They will elect a new primary Request Server via a simple leadership election algorithm any-cast. The new primary Request Server will then change DNS records to point to domain to itself. In the event that a priority is provided to the framework to prefer a certain Cloud over another (one may be cheaper, etc), the leadership election algorithm will introduce bias when generating the random IDs.

### Data 

To use multiple database services across many clouds, Hydra will also feature a distributed database access layer that will ensure consistency across all DBs. Hydra's aim is to ensure consistency and availability (and not necessarily paritioning), so writes will be distributed across all DBs. There will be a database access server (DAS) in all clouds that sits in between the webserver and the database service. Requests to write will pass through the DAS which connects to the DB service and execute the write. All DASs will be connected in this architecture, so write requests that a single DAS receives will be forwarded to every other DAS in the system, so that one write in one system fans out to a write in every system. Similarly, reads will be sourced from any system that is up, optimized for spacial locality. For example, if the AWS side of the system receives a request to read from the DB, and AWS RDS is up, it will simply read from RDS. If RDS is down, however, the AWS DAS will request a read from the GCP DAS and information will be retrieved from CloudSQL and forwarded to AWS.
 

Design Implications and Discussion:

This section discusses the implications and reasons of the design decisions made during the global architecture design.

## 5. Acceptance criteria

MVP (not in any particular order):

1) Request load balancing between clouds
2) Host level load balancing between clouds
3) Turning off host instances should not crash the application for an extended period of time
4) DNS failover in the event that the main server goes down
5) Turning off all compute instances of a Cloud should not kill the application
6) 1 Request server per Cloud

Stretch (also not in any particular order):

1) Multiple Request Servers and leadership election
2) Distributed DB writes
3) Distributed reads from any healthy database
4) Eventual consistency of new databases via replication and write queuing
5) Turning off databse instances should not bork the application
6) Multiple Request servers per Cloud

## 6.  Release Planning:

Release planning section describes how the project will deliver incremental sets of features and functions in a series of releases to completion. Identification of user stories associated with iterations that will ease/guide sprint planning sessions is encouraged. Higher level details for the first iteration is expected.

** **

For more help on markdown, see
https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet

In particular, you can add images like this (clone the repository to see details):

![alt text](https://github.com/BU-NU-CLOUD-SP18/sample-project/raw/master/cloud.png "Hover text")


