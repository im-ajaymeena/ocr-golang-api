```mermaid
graph TD
  subgraph Client
    A[User] -->|Request| B(CreateOCRTask)
    A[User] -->|Request| C(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    B -->|Starts Task| D[Task Queue]
    D -->|Worker 1| E((...))
    D -->|Worker 2| F((...))
    D -->|Worker 3| G((...))
    D -->|Worker ...| H((...))graph TDgraph TD
  subgraph Client
    A[User] -->|Request| B(CreateOCRTask)
    A[User] -->|Request| C(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    B -->|Starts Task| D[Task Queue]
    D -->|Worker 1| E((...))
    D -->|Worker 2| F((...))
    D -->|Worker 3| G((...))
    D -->|Worker ...| H((...))
  end

  subgraph Workers (Goroutines)
    E -->|Process| I(redisWriter)
    F -->|Process| I(redisWriter)
    G -->|Process| I(redisWriter)
    H -->|Process| I(redisWriter)
  end

  subgraph Redis
    I -->|Save Result| J[TaskID: Result]
  end

  subgraph Client
    C -->|Request| K(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    K -->|Retrieves Result| L((...))
    K -->|Retrieves Result| M((...))
    K -->|Retrieves Result| N((...))
    K -->|Retrieves Result| O((...))
  end

  subgraph Client
    A[User] -->|Request| B(CreateOCRTask)
    A[User] -->|Request| C(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    B -->|Starts Task| D[Task Queue]
    D -->|Worker 1| E((...))
    D -->|Worker 2| F((...))
    D -->|Worker 3| G((...))
    D -->|Worker ...| H((...))
  end

  subgraph Workers (Goroutines)
    E -->|Process| I(redisWriter)
    F -->|Process| I(redisWriter)
    G -->|Process| I(redisWriter)
    H -->|Process| I(redisWriter)
  end

  subgraph Redis
    I -->|Save Result| J[TaskID: Result]
  end

  subgraph Client
    C -->|Request| K(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    K -->|Retrieves Result| L((...))
    K -->|Retrieves Result| M((...))
    K -->|Retrieves Result| N((...))
    K -->|Retrieves Result| O((...))
  end
graph TD
  subgraph Client
    A[User] -->|Request| B(CreateOCRTask)
    A[User] -->|Request| C(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    B -->|Starts Task| D[Task Queue]
    D -->|Worker 1| E((...))
    D -->|Worker 2| F((...))
    D -->|Worker 3| G((...))
    D -->|Worker ...| H((...))
  end

  subgraph Workers (Goroutines)
    E -->|Process| I(redisWriter)
    F -->|Process| I(redisWriter)
    G -->|Process| I(redisWriter)
    H -->|Process| I(redisWriter)
  end

  subgraph Redis
    I -->|Save Result| J[TaskID: Result]
  end

  subgraph Client
    C -->|Request| K(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    K -->|Retrieves Result| L((...))
    K -->|Retrieves Result| M((...))
    K -->|Retrieves Result| N((...))
    K -->|Retrieves Result| O((...))
  end

  end

  subgraph Workers (Goroutines)
    E -->|Process| I(redisWriter)graph TD
  subgraph Client
    A[User] -->|Request| B(CreateOCRTask)
    A[User] -->|Request| C(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    B -->|Starts Task| D[Task Queue]
    D -->|Worker 1| E((...))
    D -->|Worker 2| F((...))
    D -->|Worker 3| G((...))
    D -->|Worker ...| H((...))
  end

  subgraph Workers (Goroutines)
    E -->|Process| I(redisWriter)
    F -->|Process| I(redisWriter)
    G -->|Process| I(redisWriter)
    H -->|Process| I(redisWriter)
  end

  subgraph Redis
    I -->|Save Result| J[TaskID: Result]
  end

  subgraph Client
    C -->|Request| K(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    K -->|Retrieves Result| L((...))
    K -->|Retrieves Result| M((...))
    K -->|Retrieves Result| N((...))
    K -->|Retrieves Result| O((...))
  end

    F -->|Process| I(redisWriter)
    G -->|Process| I(redisWriter)
    H -->|Process| I(redisWriter)
  end

  subgraph Redis
    I -->|Save Result| J[TaskID: Result]
  end

  subgraph Client
    C -->|Request| K(GetOCRTaskResult)
  end

  subgraph API Server (Gin)
    K -->|Retrieves Result| L((...))
    K -->|Retrieves Result| M((...))
    K -->|Retrieves Result| N((...))
    K -->|Retrieves Result| O((...))
  end
```