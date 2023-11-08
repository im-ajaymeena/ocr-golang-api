graph TD;
A[Start] --> B[Parse JSON request];
B --> |Single Image| C[Decode Image Data];
C --> D[Run OCR];
D --> E[Generate Text Response];
E --> F[Send Text Response];

B --> |Image List| G[Loop through Images];
G --> C;
G --> D;
G --> |Last Image| H[Generate Text Response List];
H --> | encode list as string | F;

B --> |No Image Data| I[No Image Error Response];
I --> J;

F --> J[End];
