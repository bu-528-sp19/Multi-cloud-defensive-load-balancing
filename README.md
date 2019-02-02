** **

## Project Description Template

The purpose of this Project Description is to present the ideas proposed and decisions made during the preliminary envisioning and inception phase of the project. The goal is to analyze an initial concept proposal at a strategic level of detail and attain/compose an agreement between the project team members and the project customer (mentors and instructors) on the desired solution and overall project direction.

This template proposal contains a number of sections, which you can edit/modify/add/delete/organize as you like.  Some key sections we’d like to have in the proposal are:

- Vision: An executive summary of the vision, goals, users, and general scope of the intended project.

- Solution Concept: the approach the project team will take to meet the business needs. This section also provides an overview of the architectural and technical designs made for implementing the project.

- Scope: the boundary of the solution defined by itemizing the intended features and functions in detail, determining what is out of scope, a release strategy and possibly the criteria by which the solution will be accepted by users and operations.

Project Proposal can be used during the follow-up analysis and design meetings to give context to efforts of more detailed technical specifications and plans. It provides a clear direction for the project team; outlines project goals, priorities, and constraints; and sets expectations.

** **

## 1.   Vision and Goals Of The Project:

Hydra will be a framework for applications in the cloud to mitigate DDoS attacks and provider outages by providing resiliency at multiple levels both intra- and cross-clouds. If parts of AWS or GCP go down, the application itself should be alive and kicking.

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

Global Architectural Structure Of the Project:

Hydra will provide resiliency at multiple levels.

### Compute

Hydra will feature a request server layer in each cloud 

This section provides a high-level architecture or a conceptual diagram showing the scope of the solution. If wireframes or visuals have already been done, this section could also be used to show how the intended solution will look. This section also provides a walkthrough explanation of the architectural structure.

 

Design Implications and Discussion:

This section discusses the implications and reasons of the design decisions made during the global architecture design.

## 5. Acceptance criteria

MVP:

1) Host heartbeating for uptime and load monitoring
2) Request load balancing based on heartbeat data
3) Turning off host instances should not bork the application
4) DNS failover in the event that the main server goes down
5) Turning off all compute instances of a Cloud should not kill the application
6) 1 Request server per Cloud

Stretch:

1) Distributed DB writes
2) Distributed reads from any healthy database
3) Eventual consistency of new databases via replication and write queuing
4) Turning off databse instances should not bork the application
5) Multiple Request servers per Cloud

## 6.  Release Planning:

Release planning section describes how the project will deliver incremental sets of features and functions in a series of releases to completion. Identification of user stories associated with iterations that will ease/guide sprint planning sessions is encouraged. Higher level details for the first iteration is expected.

** **

For more help on markdown, see
https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet

In particular, you can add images like this (clone the repository to see details):

![alt text](https://github.com/BU-NU-CLOUD-SP18/sample-project/raw/master/cloud.png "Hover text")


