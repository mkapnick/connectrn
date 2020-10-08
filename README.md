# API Endpoint REST design

### Account
- `/api/v1/account/signup/`
- `/api/v1/account/login/`
- `/api/v1/profile/`

### Restaurants
- `/api/v1/restaurants/`
- `/api/v1/restaurants/:restaurant_id/`
- `/api/v1/restaurants/:restaurant_id/tables/`
- `/api/v1/restaurants/:restaurant_id/tables/:table_id`
- `/api/v1/restaurants/:restaurant_id/tables/?date=2020-10-08`

### Reservations
- `/api/v1/restaurants/:restaurant_id/tables/:table_id/reserve/`
- `/api/v1/restaurants/:restaurant_id/tables/reserve/`
- `/api/v1/restaurants/:restaurant_id/tables/:table_id/reservations/:user_reservation_id/cancel/`

# Coding challenge constraints
**>>>> A restaurant has 4 tables of capacity 1,2,3 and 4 seats**
- Use the API and create 4 tables with the proper capacity `[POST Create restaurant table]`

**>>>> Two reservations can be made at the same time of multiple seating capacities**
- Use the API to make multiple reservations at the same time `[POST Reserve tables]`

**>>>> No single reservation can be greater than 10 people**
- When you create 4 tables with seats 4, 3, 2, and 1 (max of 10 seats), then the
`[POST Reserve tables]` request would never be able to exceed 10 spots. You will
encounter the error `not enough seats available`

**>>>> No waitlist - successful reservation or failed reservation are the only 2 states possible**
- Reservations are fully atomic

**>>>> Customer and Owner are the only 2 people who will use this API**
- You can create a `user` account
- You can create an `owner` account
- `owner` accounts are associated with a restaurant. The implementation is
crude and naive but it's simple for the purpose of the assignment

**>>>> Both customer and owner can make/cancel a reservation**
- The following endpoints will accomplish this:
- `[POST Reserve table]`
- `[POST Cancel reservation]`

**>>>> Create 3 REST APIs**
- I must have missed this part! I created a lot more endpoints :)

**>>>> As a User (customer or owner), I want to make a reservation**
- The following endpoint will accomplish this:
- `[POST Reserve a table]`

**>>>> As a User, I want to cancel a reservation**
- The following endpoint will accomplish this:
- `[POST Cancel reservation]`

**>>>> As a User, at any given point, I want to check which tables are free and of what
capacity**
- The following endpoint will accomplish this:
- `[GET Restaurant tables]`

# Video Demo
[![Demo](https://img.youtube.com/vi/KNeQbMrvGZU/0.jpg)](https://www.youtube.com/watch?v=KNeQbMrvGZU "Demo")

# Run in insomnia [as seen in the video]
[![Run in Insomnia}](https://insomnia.rest/images/run.svg)](https://insomnia.rest/run/?label=Restaurant%20API&uri=https%3A%2F%2Fgithub.com%2Fmkapnick%2Fconnectrn%2Fblob%2Fmaster%2Finsomnia_2020-10-08.json)

# Prod Deployment
- https://connectrn-api.herokuapp.com/

