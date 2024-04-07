# BannerGetRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TagIds** | Pointer to **[]int32** | Идентификаторы тэгов | [optional] 
**FeatureId** | Pointer to **int32** | Идентификатор фичи | [optional] 
**Content** | Pointer to **map[string]interface{}** | Содержимое баннера | [optional] 
**IsActive** | Pointer to **bool** | Флаг активности баннера | [optional] 

## Methods

### NewBannerGetRequest

`func NewBannerGetRequest() *BannerGetRequest`

NewBannerGetRequest instantiates a new BannerGetRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBannerGetRequestWithDefaults

`func NewBannerGetRequestWithDefaults() *BannerGetRequest`

NewBannerGetRequestWithDefaults instantiates a new BannerGetRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTagIds

`func (o *BannerGetRequest) GetTagIds() []int32`

GetTagIds returns the TagIds field if non-nil, zero value otherwise.

### GetTagIdsOk

`func (o *BannerGetRequest) GetTagIdsOk() (*[]int32, bool)`

GetTagIdsOk returns a tuple with the TagIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTagIds

`func (o *BannerGetRequest) SetTagIds(v []int32)`

SetTagIds sets TagIds field to given value.

### HasTagIds

`func (o *BannerGetRequest) HasTagIds() bool`

HasTagIds returns a boolean if a field has been set.

### GetFeatureId

`func (o *BannerGetRequest) GetFeatureId() int32`

GetFeatureId returns the FeatureId field if non-nil, zero value otherwise.

### GetFeatureIdOk

`func (o *BannerGetRequest) GetFeatureIdOk() (*int32, bool)`

GetFeatureIdOk returns a tuple with the FeatureId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeatureId

`func (o *BannerGetRequest) SetFeatureId(v int32)`

SetFeatureId sets FeatureId field to given value.

### HasFeatureId

`func (o *BannerGetRequest) HasFeatureId() bool`

HasFeatureId returns a boolean if a field has been set.

### GetContent

`func (o *BannerGetRequest) GetContent() map[string]interface{}`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *BannerGetRequest) GetContentOk() (*map[string]interface{}, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *BannerGetRequest) SetContent(v map[string]interface{})`

SetContent sets Content field to given value.

### HasContent

`func (o *BannerGetRequest) HasContent() bool`

HasContent returns a boolean if a field has been set.

### GetIsActive

`func (o *BannerGetRequest) GetIsActive() bool`

GetIsActive returns the IsActive field if non-nil, zero value otherwise.

### GetIsActiveOk

`func (o *BannerGetRequest) GetIsActiveOk() (*bool, bool)`

GetIsActiveOk returns a tuple with the IsActive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsActive

`func (o *BannerGetRequest) SetIsActive(v bool)`

SetIsActive sets IsActive field to given value.

### HasIsActive

`func (o *BannerGetRequest) HasIsActive() bool`

HasIsActive returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


