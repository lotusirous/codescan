# System Design

Codescan is a single program that comprises two primary components that run concurrently using Go routines.

- The web HTTP server provides the REST API for users.
- The task scheduler retrieves scan tasks from the web server, schedules, and updates the task to run.
- The manager configures and spawns a number of workers. These workers download public Git repositories, scan them for vulnerabilities using defined rules, and store the findings in a database.

![Architecture](architecture.png 'Codescan Architecture')

When user submits a new scan task, the manager will manages a transaction for creating/updating the scan task with the database.

## Discussion

Separating the worker into a separate process may seem like a logical solution, however, it would also come with additional complexities. This is because it would require implementing a distributed queue system, such as Redis or a cloud-based task queue. These types of queues allow for the distribution of tasks to multiple workers and can handle a large number of concurrent requests. However, implementing a distributed queue also introduces additional considerations and events that must be handled.

For example, the enqueue process may fail or the network connection may drop, which would require retrying a task for the worker. Additionally, when a worker completes or fails a task, it must send the result back to the manager, which could be done via Remote Procedure Calls (RPC) or a task queue. These extra events and processes add complexity and require additional development time and resources. For this reason, it's not considered a good choice.

While this approach may present challenges in terms of the size of the in-memory queue, these limitations can be mitigated through adjustments to the configuration or by utilizing multiple instances.

## Improvements

**Rate limiter**: The resources available for a user depend on the cost they pay. It is important to limit the number of tasks for the reliability of our service.

**Caching**: To improve database query performance, the web server can implement caching, which temporarily stores frequently accessed data to reduce load on the database and improve response time.
