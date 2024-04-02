
# Likes Microservice

Management of likes and dislikes relation between users and media (movie, books, music).

---
<br />

## API Reference

### Instance Management

#### Delete Media

Delete media's node (instance) given a media id and media type.

```http
  DELETE /likes/media/${id}
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required**. media id |

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `media_type` | `enum('MOV', 'SON' , 'BOO')` | **Required**. media type |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `204` | `success` | No content|
| `400` | `error` | "Media id not provided" |
| `400` | `error` | "Media type not provided" |
| `500` | `error` | Any other error message|

#### Delete user

Delete user's node (instance) given a user id.

```http
  DELETE /likes/user/${id}
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required**. user id |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `204` | `success` | No content|
| `400` | `error` | "User id not provided" |
| `500` | `error` | Any other error message|



### Likes Management

#### Create Like

Create new like/dislike relation.

```http
  POST /likes
```

```typescript
// Body interface
interface Create_Like{
  user_id: number
  media_id: number
  media_type: 'MOV' | 'BOO' | 'SON'
  like_type: 'LK' | 'DLK'
  rating?: float
  wishlist?: boolean
}
```

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `201` | `success` | "Relation created"|
| `400` | `error` | "Guard failed" |
| `500` | `error` | Any other error message|

#### Delete Like

Delete like/dislike relation.

```http
  DELETE /likes
```

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `user_id` | `int` | **Required**. user id |
| `media_id` | `int` | **Required**. media id |
| `media_type` | `enum('MOV', 'SON' , 'BOO')` | **Required**. media type |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `204` | `success` | No content |
| `400` | `error` | "User id not provided" |
| `400` | `error` | "Media id not provided" |
| `400` | `error` | "Media type not provided" |
| `404` | `error` | "Relation not found" |
| `500` | `error` | Any other error message|

#### Get User Likes

Returns all the like/dislikes relations make by a given user.

```http
  GET /likes/user/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required**. user id |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `200` | `success` | "Returns all the user likes relations|
| `400` | `error` | "User id not provided" |
| `404` | `error` | "User not found" |
| `500` | `error` | Any other error message|

```typescript

interface Like_Relation{
  id: int // Media id
  user_id: int
  type: 'MOV' | 'BOO' | 'SON' // Media type
  rating?: float // given by the user searched
  like_type: 'LK' | 'DLK' | 'BLK' // Liked | Disliked | Blank (no info yet)
  wishlist: boolean // Inside user wishlist? Yes or No
}

// Body interface
interface Get_Likes{
  id: number // User id
  movies: Like_Relation[]
  books: Like_Relation[]
  songs: Like_Relation[]
}
```

#### Get Media Likes

Returns all the like/dislikes relations of a given media id (book, movie, song).

```http
  GET /likes/media/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required**. media id |

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `media_type` | `enum('MOV', 'SON' , 'BOO')` | **Required**. media type |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `200` | `success` | "Returns all the media likes |
| `400` | `error` | "Media id not provided" |
| `400` | `error` | "Media type not provided" |
| `404` | `error` | "Media not found" |
| `500` | `error` | Any other error message|

```typescript

interface Like_Relation{
  id: number // Media id
  user_id: number
  type: 'MOV' | 'BOO' | 'SON' // Media type
  rating?: float // given by the user searched
  like_type: 'LK' | 'DLK' | 'BLK' // Liked | Disliked | Blank (no info yet)
  wishlist: boolean // Inside user wishlist? Yes or No
}

// Body interface
interface Get_Likes_Media{
  likes: Like_Relation[]
  avg_rating: float
}
```

#### Get Rating

Get average rating of a media.

```http
  GET /likes/average/${id}
```

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `media_type` | `enum('MOV', 'SON' , 'BOO')` | **Required**. media type |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `200` | `success` | Returns media average |
| `400` | `error` | "Media id not provided" |
| `400` | `error` | "Media type not provided" |
| `500` | `error` | Any other error message|

```typescript
// Body interface
interface Get_Likes_Media{
  id: number // Media id
  type: 'MOV' | 'SON' | 'BOO' // Media Type
  avg_rating: float
}
```

### Wishlist

#### Get media on Wishlist

Get wishlist of a given user divided by media type.

```http
  GET /likes/wishlist/${id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required**. user id |


| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `200` | `success` | "Returns the user wishlist|
| `400` | `error` | "User id not provided" |
| `404` | `error` | "User not found" |
| `500` | `error` | Any other error message|


```typescript
// Body interface
interface Get_Likes{
  id: number // User id
  movies: number[] // Wishlist movie ids
  books: number[] // Wishlist book ids
  songs: number[] // Wishlist song ids
}
```

---
<br />
<br />
<br />

## Deployment

To deploy this project:

Deploy database

```bash
  ./run_DB.sh
```

Deploy API

```bash
  ./run_API.sh
```

Call API using (http://localhost:3000)[http://localhost:3000]

## Run Locally

Clone the project

```bash
  git clone https://github.com/QuickCrafts/PerfectPick_Likes_ms.git
```

Go to the project directory

```bash
  cd PerfectPick_Likes_ms
```

Open Database

```bash
  ./run_DB.sh
```

Start the server

```bash
  make run
```
Call API using (http://localhost:3000)[http://localhost:3000]
