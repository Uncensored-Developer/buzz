# BUZZ

This is a mini API for a simple dating app that fetches potential user profiles based on location
using h3 geospatial index



All that's needed to setup and run this app is Golang and Docker installed locally.

---
#### Run Tests
> go test -v ./...

#### Run API
> docker compose up

#### Shutdown
> docker compose down

---

### APIs
On start up the app would run on [http://localhost:8003](http://localhost:8003)

### Create user: 

This creates a random user with a default password of `password123`

**URL:** /user/create

**Method:** POST

**Body:** NONE

**RESPONSE**
```json

```

===================================

### Login:

**URL:** /login

**Method:** POST

**Body:** 
```json
{
    "email": "dorthyspinka@hintz.io",
    "password": "password123"
}
```

**RESPONSE**
```json
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjh9.pJgdc2DQjtWvHA2rP5_WMzD2Z0Ldd4PPtWWw4ocfGkM"
}
```

===================================

### Update Location:

**URL:** /user/update

**HEADERS:** `Content-Type: application/json`, `Authorization: Bearer <Token>`

**Method:** PUT

**Body:**
```json
{
  "Latitude": 41.089097,
  "Longitude": -81.349288
}
```

**RESPONSE**
```json
{
  "result": {
    "id": 8,
    "email": "dorthyspinka@hintz.io",
    "password": "77546b6a714655717144705547696e414a71485acbfdac6008f9cab4083784cbd1874f76618d2a97",
    "name": "Olga Goodwin",
    "gender": "F",
    "age": 0,
    "longitude": -81.349288,
    "latitude": 41.089097,
    "h3_index": 617744365936443391
  }
}
```

===================================

### Swipe:

Use to respond to a user's profile (either YES or NO)

**URL:** /swipe

**HEADERS:** `Content-Type: application/json`, `Authorization: Bearer <Token>`

**Method:** POST

**Body:**
```json
{
  "userId": 1,
  "action": "YES" // or NO
}
```

**RESPONSE**
```json
{
  "results": {
    "matched": true,
    "matchedID": 1
  }
}
```

===================================

### Discover:

Use to fetch potential matches based on age range, gender and radius

**URL:** /discover

**HEADERS:** `Content-Type: application/json`, `Authorization: Bearer <Token>`

**Method:** GET

**QUERY PARAMS:** `age_range`, `gender`, `distance_from`
```text
age_range should follow the `{minimum age}-{maximum age}` e.g 18-30.

gender should be either M, F or 0 for male, female and others respectively.

distance_from should be an integer representing the distance from user in Kilometers.
```

**RESPONSE**
```json
{
  "results": [
    {
      "id": 12,
      "name": "Lottie Ledner",
      "gender": "M",
      "age": 33,
      "distanceFromMe": 60
    },
    {
      "id": 13,
      "name": "Vern King",
      "gender": "M",
      "age": 23,
      "distanceFromMe": 20
    },
    {
      "id": 16,
      "name": "Clifton Hodkiewicz",
      "gender": "M",
      "age": 38,
      "distanceFromMe": 40
    }
  ]
}
```
---


### ASSUMPTIONS

* Minimum age for creating a profile is 18 years.
* Maximum age is 60 years.
* User location is stored as longitude and latitude.
* For a dating app we need to find users within a few kilometers radius, hence I used a resolution of 9 for H3 indexing which is a good balance between performance and precision.
* Attractiveness is calculated by number of YES swipes on a user's profile.


### GOING FORWARD
* More refactors to be done (probably introduce a dependency injection framework to properly handle dependencies).
* More extensive tests.
* Add test coverage.
* Handle NO swipes on user profiles, Probably store them to exclude the involved profiles from being shown to the user in future.
* Build separate service that regularly aggregates potential matches for active users based on more profile info, activities and other statistics (might include some ML).
* Add Super Like Feature that ranks profile to the top of the list of potential matches.
* What happens after match????
