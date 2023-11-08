graph TD
  subgraph Client
    A[User] -->|Request| B(CreateOCRTask)
    A[User] -->|Request| C(GetOCRTaskResult)
  end

  subgraph API_Server_Gin
    B -->|Starts Task| D(Task Queue)
  end

  subgraph Workers_Goroutines
    D -->|Worker 1| O(OCR)
    O -->|Recognized Text, TaskID, <br> TaskStatus| E(redisWriter)
    
    D -->|Worker 2| P(OCR)
    P -->|Recognized Text, TaskID, <br> TaskStatus| F(redisWriter)
    D -->|Worker 3| Q(OCR)
    Q -->|Recognized Text, TaskID, <br> TaskStatus| G(redisWriter)
  end

  subgraph Redis
    E -->|Save Result| H[TaskID: Result]
    F -->|Save Result| H
    G -->|Save Result| H
  end

  subgraph API_Server_Gin
    C -->|Retrieves Result| J(GetOCRTaskResult)
  end

  subgraph API_Server_Gin
    J -->|Retrieves Result| K[Check Task Status]
    K -->|Type: String| L[Retrieve Result from Redis]
    K -->|Type: List| M[Retrieve Subtask IDs]
    M -->|Retrieve Subtask Results| N[Retrieve Subtask Results from Redis]
  end
