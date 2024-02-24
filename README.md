# ShopAnythingLagos API

This is the API for the ShopAnythingLagos project. It is a RESTful API built with Go.

## Getting Started

- The base URL for the API is `https://localhost:4500/`
- It uses Bear Token Authentication

## Prerequisites

- Postman or any other API testing tool
- Download the postman collection [here](https://res.cloudinary.com/dfbebf7x0/raw/upload/v1708439359/SAL.postman_collection_o0m8ff.json)

## Endpoints

- ### Products
    - Get all products
        - **GET** `/product`
        - Use the `page` and `limit` query parameters to paginate the results
        - Use `search` query parameter to search for products
        - Use `sortKey` and `sortOrder` query parameters to sort the results
        - The default value for `page` is 1 and `limit` is 10
        - The default value for `sortKey` is `createdAt` and `sortOrder` is `desc`
        - **Response Body**
          ```json
          {
            "products": [],
            "message": "string",
            "meta": {
              "currentPage": "number",
              "totalPages": "number",
              "limit": "number",
              "totalProducts": "number",
              "nextPage": "string",
              "prevPage": "string"
            }
          }
          ```

    - Get a single product
        - **GET** `/product/:id`
        - It requires the `id` of the product as a URL parameter
        - **Response Body**
          ```json
          {
            "id": "string",
            "name": "string",
            "description": "string",
            "price": "number",
            "createdAt": "string",
            "updatedAt": "string"
          }
          ```

    - Create a product
        - **POST** `/product?skuId=skuId`
        - Its an authenticated route, hence it requires a bearer token
        - **Request Body**
          ```json
          {
            "name": "string",
            "description": "string",
            "price": "number"
          }
          ```
        - **Response Body**
          ```json
          {
            "id": "string",
            "name": "string",
            "description": "string",
            "price": "number",
            "createdAt": "string",
            "updatedAt": "string"
          }
          ```

    - Update a product
        - **PUT** `/product/:id?skuId=skuId`
        - Its an authenticated route, hence it requires a bearer token
        - It requires the `id` of the product as a URL parameter
        - **Request Body**
          ```json
          {
            "name": "string", 
            "description": "string", 
            "price": "number" 
          }
          ```
        - **Response Body**
          ```json
          {
            "id": "string",
            "name": "string",
            "description": "string",
            "price": "number",
            "createdAt": "string",
            "updatedAt": "string"
          }
          ```

    - Delete a product
        - **DELETE** `/product/:id?skuId=skuId`
        - Its an authenticated route, hence it requires a bearer token
        - It requires the `id` of the product as a URL parameter
        - **Response Body**
          ```json
          {
            "message": "string"
          }
          ```
