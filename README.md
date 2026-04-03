# Golang_Chat_Advisor
Build a backend service to create an AI advisor where client can create a chat advisor. It is a microservice-style architecture because the application is designed as a small, independent service using clear separation between API, business logic, and data layers, and is deployed independently using serverless infrastructure.

*Pls note it is not using a real LLM and is only using mocked logic/ responses from advisors.
The purpose of this project is to demonstrate backend architecture, system design, session handling, API design and play around with AWS services for infra/deployment for learning purposes.
 
Users can:
- create a companion
- start a chat session
- send a message
- get a mocked AI reply
- fetch session history

Tech Stack Used: 
- Go
  - Backend code to handle requests, process logic and responses. Also talks to DynamoDB
- AWS Lambda
  - For quick and easy serverless deployment. Executes code only when requests come and stops after a time interval.  
- API Gateway
  - Comes before Lambda. Receives HTTP requests from user and sends them to Lambda. Also retrieves response.
- DynamoDB
  - Uses 2 tables. One for advisors created and one for chat sessions and messages
- AWS SAM
  - Builds cloud setup using template.yaml. Creates Lambda, API Gateway and DynamoDB tables
- GitHub Actions
  - CICD tool to build code and deploy to AWS. Runs on every push (for now)
