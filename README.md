Service 1 - Task Producer
-------------------------

This document describes the Task Producer service, responsible for generating and sending tasks to the Task Consumer service.

### Functionality

*   Generates tasks with the following properties:

    *   **Task Type:** Random integer between 0 and 9 (inclusive)

    *   **Value:** Random integer between 0 and 99 (inclusive)

*   Persists tasks in a database with an initial state of "received".

*   Tracks the number of produced tasks using Prometheus metrics.


### Configuration

*   **Prometheus:**

    *   port: Port to expose Prometheus metrics.

    *   endpoint: Fixed endpoint path for Prometheus metrics (/metrics).

*   **Communication:** Configuration for communication between Task Producer and Task Consumer (details depend on chosen communication method).

*   **Max Backlog:** Maximum number of unprocessed messages the producer can generate before stopping.

*   **Logging:**

    *   level: Logging level (e.g., info, debug, error).

    *   format: Log format (console or JSON structured).

*   **Profiling:** Port to expose for profiling.

*   **Message Production Rate:** Number of messages produced per second.

*   **Database:** Connection configuration for the database used to persist tasks.


### Version Information

Calling the service with the -version command line argument returns the version with which the service was built.

### Running the Service

(Specific instructions will depend on your development environment and deployment choices.)

1.  Build the service.

2.  Start the service with the desired configuration options.


### Monitoring

*   Use Prometheus to monitor the number of produced tasks and other relevant metrics.


### Logging

*   Logs are written based on the configured level and format.


### Development

*   Code for the Task Producer service should be located in the service1 directory.

*   Unit tests should be written to ensure proper functionality.


### Additional Notes

*   This document provides a high-level overview. Refer to the code for specific implementation details.

*   Consider implementing error handling and retries for communication failures.

Service 2 - Task Consumer
-------------------------

### Functionality

*   Receives tasks from the Task Producer.

*   Sets the task state to "processing" in the database upon receipt.

*   Handles tasks with a simulated delay based on their value.

*   Sets the task state to "done" in the database after processing.

*   Enforces a rate limit on incoming tasks.

*   Tracks various metrics using Prometheus:

    *   Number of tasks being processed and done.

    *   Number of tasks per task type.

    *   Total sum of the "value" field for each task type.

*   Logs task content and calculated total sum for each incoming task.


### Configuration

*   **Prometheus:**

    *   port: Port to expose Prometheus metrics.

    *   endpoint: Fixed endpoint path for Prometheus metrics (/metrics).

*   **Communication:** Configuration for communication between Task Producer and Task Consumer (details depend on chosen communication method).

*   **Logging:**

    *   level: Logging level (e.g., info, debug, error).

    *   format: Log format (console or JSON structured).

*   **Profiling:** Port to expose for profiling.

*   **Message Consumption Rate:** Rate limit for message consumption.

*   **Database:** Connection configuration for the database used to persist tasks.


### Version Information

Calling the service with the -version command line argument returns the version with which the service was built.

### Running the Service

(Specific instructions will depend on your development environment and deployment choices.)

1.  Build the service.

2.  Start the service with the desired configuration options.


### Monitoring

*   Use Prometheus to monitor the number of tasks being processed, and done, and other relevant metrics.


### Logging

*   Logs are written based on the configured level and format.


### Development

*   The code for the Task Consumer service should be located in the service2 directory.

*   Unit tests should be written to ensure proper functionality.


### Additional Notes

*   This document provides a high-level overview. Refer to the code for specific implementation details.

*   Consider implementing error handling and retries for communication failures.

*   The simulated delay in task handling can be replaced with actual processing logic.