# Souuup Examples

This directory contains examples demonstrating various use cases for the Souuup validation library.

## Basic Validation

A simple example showing how to validate basic data structures.

```bash
cd basic
go run main.go
```

## HTTP Validation

An example showing how to use Souuup for validating HTTP requests in a web API.

```bash
cd http
go run main.go
```

Then send a POST request to `http://localhost:8080/register` with a JSON body like:

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "Password123",
  "confirmPassword": "Password123",
  "age": 25
}
```

You can use tools like curl, Postman, or httpie to send the request.

## Complex Validation

An example showing complex validation scenarios with nested schemas, custom rules, and varied data types.

```bash
cd complex
go run main.go
```

## Creating Your Own Examples

Feel free to modify these examples or create your own. The key patterns to follow are:

1. Define your data structures
2. Create validation rules specific to your use case
3. Build a schema mapping field names to validatable fields
4. Create a validator named `uuu` with your schema
5. Call `Validate()` to check for validation errors

For more details, see the main README and package documentation.
