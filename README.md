# Profiles Microservice 👤

## Installation 🚀
Clone the repository: 📥
git clone git@github.com:meathub-com/profiles.git

Copy code

Install the required dependencies: 📦
go mod init profiles
go mod tidy

mathematica
Copy code

Run the service: 🏃‍♂️
go run main.go
docker-compose up --build
make docker-build

shell
Copy code

## Documentation: 📚
Access the API documentation at: 🌐
http://localhost:8080/swagger/index.html

markdown
Copy code

## Usage ℹ️
1. Register a new user by sending a POST request to `/profiles/register`.
2. Authenticate a user by sending a POST request to `/profiles/login`.
3. Retrieve user information by sending a GET request to `/profiles/{user_id}`.
4. Update user information by sending a PUT request to `/profiles/{user_id}`.
5. Delete a user by sending a DELETE request to `/profiles/{user_id}`.

## Contributing 🤝
We welcome contributions to improve the Profiles Microservice! To contribute, follow these steps:
1. Fork the repository.
2. Create a new branch.
3. Make your changes and commit them.
4. Push your changes to your fork.
5. Submit a pull request.

Please ensure that your code adheres to our coding guidelines and passes all tests before submitting a pull request.

## Bug Reports 🐞
If you encounter any issues or bugs, please report them on the GitHub issue tracker: 🐛
https://github.com/meathub-com/profiles/issues

csharp
Copy code

## License 📝
The Profiles Microservice is released under the MIT License. See the LICENSE file for more details.