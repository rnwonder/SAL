# ShopAnythingLagos API

This is the API for the ShopAnythingLagos project. It is a RESTful API built with Go.

## Getting Started

- The base URL for the API is `https://localhost:4500/`
- It uses Bear Token Authentication

## Prerequisites

- Postman or any other API testing tool
- Download the postman collection [here](https://somelink.com)

## Endpoints

- ### Auth
  - Register as a merchant
    - POST `/auth`
    - It requires a JSON body with the following fields:
      - `email` - The email of the merchant - `string`
      - `password` - The password of the merchant - `string`
      - `name` - The name of the merchant - `string`
      - `skuId` - The SKU ID of the merchant - `string`
    
  - Login as a merchant
    - POST `/auth/login`
    - It requires a JSON body with the following fields:
      - `email` - The email of the merchant - `string`
      - `password` - The password of the merchant - `string`

- ### Products
  - Get all products
    - GET `/product`
    - It is a public route, hence it does not require authentication
    - Use the `page` and `limit` query parameters to paginate the results
    - Use search query parameter to search for products
    - use sortKey and sortOrder query parameters to sort the results
    - The default value for `page` is 1 and `limit` is 10
    - The default value for `sortKey` is `createdAt` and `sortOrder` is `desc`
    - The default value for `search` is an empty string
    - It returns 
      - `products` - An array of products 
      - `message` - A message indicating the success of the request - `string`
      - `meta` - An object containing pagination information"
        - `currentPage` - The current page - `number`
        - `totalPages` - The total number of pages - `number`
        - `limit` - The limit of products per page - `number`
        - `totalProducts` - The total number of products - `number`
        - `nextPage` - The query parameters for the next page - `string`
        - `prevPage` - The query parameters for the previous page - `string`
    
  - Get a single product
    - GET `/product/:id`
    - It is a public route, hence it does not require authentication
    - It requires the `id` of the product as a URL parameter
    - It returns a single product
      - `id` - The ID of the product - `string`
      - `name` - The name of the product - `string`
      - `description` - The description of the product - `string`
      - `price` - The price of the product - `number`
      - `createdAt` - The date the product was created - `string`
      - `updatedAt` - The date the product was last updated - `string`

  - Create a product
    - POST `/product`
    - Its an authenticated route, hence it requires a bearer token
    - It requires a JSON body with the following fields:
      - `name` - The name of the product - `string`
      - `description` - The description of the product - `string`
      - `price` - The price of the product - `number`
    - It returns the created product
  
    - Update a product
      - PUT `/product/:id`
      - Its an authenticated route, hence it requires a bearer token
      - It requires the `id` of the product as a URL parameter
      - It requires a JSON body with any of the following fields:
        - `name` - The name of the product - `string`
        - `description` - The description of the product - `string`
        - `price` - The price of the product - `number`
        - It returns the updated product
        
    - Delete a product
      - DELETE `/product/:id`
      - Its an authenticated route, hence it requires a bearer token
      - It requires the `id` of the product as a URL parameter
      - It returns a message indicating the success of the request