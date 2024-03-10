
# Likes Microservice

Management of likes and dislikes relation between users and media (movie, books, music).

---
<br />

## API Reference

### Instance Management

#### Create User

Create user's node (instance) given a user id.

```http
  POST /likes/user/${id}
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required**. user id |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `201` | `success` | "User instance created"|
| `400` | `error` | "User id not provided" |
| `400` | `error` | "User already has an instance" |
| `500` | `error` | Any other error message|

#### Create Media

Create media's node (instance) given a media id and media type.

```http
  POST /likes/media/${id}
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required**. media id |

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `media_type` | `enum('MOV', 'SON' , 'BOO')` | **Required**. media type |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `201` | `success` | "Media instance created"|
| `400` | `error` | "Media id not provided" |
| `400` | `error` | "Media type not provided" |
| `400` | `error` | "Media already has an instance" |
| `500` | `error` | Any other error message|

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
| `404` | `error` | "Media not found" |
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
| `404` | `error` | "User not found" |
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
  user_id: string
  media_id: string
  media_type: 'MOV' | 'BOO' | 'SON'
  like_type: 'LK' | 'DLK' | 'BLK'
  rating?: float
  wishlist: boolean
}
```

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `201` | `success` | "Relation created"|
| `400` | `error` | "Guard failed" |
| `400` | `error` | "User and Media already has a relation" |
| `500` | `error` | Any other error message|

#### Update Like

Update like/dislike relation.

```http
  PUT /likes
```

```typescript
// Body interface
interface Update_Like{
  like_type?: 'LK' | 'DLK' | 'BLK'
  rating?: float
  wishlist?: boolean
}
```

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `user_id` | `int` | **Required**. user id |
| `media_id` | `int` | **Required**. media id |
| `media_type` | `enum('MOV', 'SON' , 'BOO')` | **Required**. media type |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `201` | `success` | "Relation updated"|
| `400` | `error` | "Guard failed" |
| `400` | `error` | "User id not provided" |
| `400` | `error` | "Media id not provided" |
| `400` | `error` | "Media type not provided" |
| `404` | `error` | "Relation not found" |
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
  id: string // Media id
  user_id: string
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
  id: string // Media id
  user_id: string
  type: 'MOV' | 'BOO' | 'SON' // Media type
  rating?: float // given by the user searched
  like_type: 'LK' | 'DLK' | 'BLK' // Liked | Disliked | Blank (no info yet)
  wishlist: boolean // Inside user wishlist? Yes or No
}

// Body interface
interface Get_Likes_Media{
  likes: interface Like_Relation[]
  avg_rating: float
}
```

#### Get Rating

Get average rating of a relation.

```http
  GET /likes/average/
```

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `user_id` | `int` | **Required**. user id |
| `media_id` | `int` | **Required**. media id |
| `media_type` | `enum('MOV', 'SON' , 'BOO')` | **Required**. media type |

| Response Status | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `200` | `success` | Returns media average |
| `400` | `error` | "User id not provided" |
| `400` | `error` | "Media id not provided" |
| `400` | `error` | "Media type not provided" |
| `404` | `error` | "Relation not found" |
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

To deploy this project run

[//]: <> (@todo correct)

```bash
  npm run deploy
```

## Run Locally

Clone the project

[//]: <> (@todo correct all)

```bash
  git clone https://github.com/QuickCrafts/PerfectPick_Likes_ms.git
```

Go to the project directory

```bash
  cd PerfectPick_Likes_ms
```

Install dependencies

```bash
  todo
```

Start the server

```bash
  todo
```
