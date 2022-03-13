
# Busha assessment

A simple crud rest api to get movies and character data from [api](https://swapi.dev)
movies are cached in redis and subsequent calls data is fetched from the cache.

characters of each movie is also fetched and chached also.

# Requirements
To run this application, you'll need docker installed.
clone the [repository](https://github.com/odamilola36/busha_)
* cd into busha
* run ```docker-compose up --build```
* the application runs on port ```8080```
* refer to [**API Reference**](Response.md) for usage

## API Reference

All URIs are relative to *http://localhost:8080*


### CharacterApi

Method | HTTP request | Description
------------- | ------------- | -------------
[**getCharacters**](CharacterApi.md#getCharacters) | **GET** /characters/{movieId} | characters endpoint

#### **getCharacters**
> Response getCharacters(movieId, sort, order, filter)

characters endpoint
### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**movieId** | **Long**| id of movie |
**sort** | **String**| key for sorting, valid values are name, height, gender | [optional]
**order** | **String**| order of results, asc (ascending) or desc[default](descending) | [optional]
**filter** | **String**| filter parameters based on gender values: male, female, n/a or blank(same as not sorting the responses at all) | [optional]

### Return type

[**Response**](Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | successful operation |  -  |
**400** | Invalid ID supplied |  -  |
**404** | Character not found |  -  |
**500** | Internal server error |  -  |

### CommentApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**create**](CommentApi.md#create) | **POST** /comment | create comment for a movie
[**deleteComment**](CommentApi.md#deleteComment) | **DELETE** /comment/{commentId} | Delete Comment by ID
[**getAll**](CommentApi.md#getAll) | **GET** /comments/{movieId} | get all comments for a movie
[**getCommentById**](CommentApi.md#getCommentById) | **GET** /comment/{commentId}/{movieId} | get comments by commentId and movieId
[**updateComment**](CommentApi.md#updateComment) | **PATCH** /comment/{commentId} | update a comment

# **create**
> Response create()

create comment for a movie

method to add comment to a movie with specified Id


### Parameters
This endpoint does not need any parameter.

### Return type

[**Response**](Response.md)

### Authorization

No authorization required

### HTTP request headers

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getMovies**](MoviesApi.md#getMovies) | **GET** /movies | Get all movies from the api

### **getMovies**
> getMovies()

Get all movies from the api

### Parameters
This endpoint does not need any parameter.

### Return type

[**Response**](Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Ok |  -  |
**405** | Method not supported |  -  |
**500** | Internal server error |  -  |

# Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**status** | **Long** |  |  [optional]
**error** | **Object** |  |  [optional]
**message** | **String** | status message |  [optional]
**data** | **Object** |  |  [optional]


