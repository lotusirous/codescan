# System Design

Codescan consists of the following 2 main components

- The web http server renders the UI/API for user to submit the public repo.
- The task scheduler manages the queue and worker. The worker that downloads public git repositories, scans them for vulnerabilities based on rules, and stores the scan data in the database.

![Architecture](architecture.png 'Codescan Architecture')

## Improvements

**Rate limiter**: The resources available for a user depend on the cost they pay. It is important to limit the number of tasks for the reliability of our service.
