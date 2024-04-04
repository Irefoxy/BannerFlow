# BannerIdDeleteRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TagIds** | Pointer to **[]int32** | Идентификаторы тэгов | [optional] 
**FeatureId** | Pointer to **NullableInt32** | Идентификатор фичи | [optional] 
**Content** | Pointer to **map[string]interface{}** | Содержимое баннера | [optional] 
**IsActive** | Pointer to **NullableBool** | Флаг активности баннера | [optional] 

## Methods

### NewBannerIdDeleteRequest

`func NewBannerIdDeleteRequest() *BannerIdDeleteRequest`

NewBannerIdDeleteRequest instantiates a new BannerIdDeleteRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBannerIdDeleteRequestWithDefaults

`func NewBannerIdDeleteRequestWithDefaults() *BannerIdDeleteRequest`

NewBannerIdDeleteRequestWithDefaults instantiates a new BannerIdDeleteRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTagIds

`func (o *BannerIdDeleteRequest) GetTagIds() []int32`

GetTagIds returns the TagIds field if non-nil, zero value otherwise.

### GetTagIdsOk

`func (o *BannerIdDeleteRequest) GetTagIdsOk() (*[]int32, bool)`

GetTagIdsOk returns a tuple with the TagIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTagIds

`func (o *BannerIdDeleteRequest) SetTagIds(v []int32)`

SetTagIds sets TagIds field to given value.

### HasTagIds

`func (o *BannerIdDeleteRequest) HasTagIds() bool`

HasTagIds returns a boolean if a field has been set.

### SetTagIdsNil

`func (o *BannerIdDeleteRequest) SetTagIdsNil(b bool)`

 SetTagIdsNil sets the value for TagIds to be an explicit nil

### UnsetTagIds
`func (o *BannerIdDeleteRequest) UnsetTagIds()`

UnsetTagIds ensures that no value is present for TagIds, not even an explicit nil
### GetFeatureId

`func (o *BannerIdDeleteRequest) GetFeatureId() int32`

GetFeatureId returns the FeatureId field if non-nil, zero value otherwise.

### GetFeatureIdOk

`func (o *BannerIdDeleteRequest) GetFeatureIdOk() (*int32, bool)`

GetFeatureIdOk returns a tuple with the FeatureId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeatureId

`func (o *BannerIdDeleteRequest) SetFeatureId(v int32)`

SetFeatureId sets FeatureId field to given value.

### HasFeatureId

`func (o *BannerIdDeleteRequest) HasFeatureId() bool`

HasFeatureId returns a boolean if a field has been set.

### SetFeatureIdNil

`func (o *BannerIdDeleteRequest) SetFeatureIdNil(b bool)`

 SetFeatureIdNil sets the value for FeatureId to be an explicit nil

### UnsetFeatureId
`func (o *BannerIdDeleteRequest) UnsetFeatureId()`

UnsetFeatureId ensures that no value is present for FeatureId, not even an explicit nil
### GetContent

`func (o *BannerIdDeleteRequest) GetContent() map[string]interface{}`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *BannerIdDeleteRequest) GetContentOk() (*map[string]interface{}, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *BannerIdDeleteRequest) SetContent(v map[string]interface{})`

SetContent sets Content field to given value.

### HasContent

`func (o *BannerIdDeleteRequest) HasContent() bool`

HasContent returns a boolean if a field has been set.

### SetContentNil

`func (o *BannerIdDeleteRequest) SetContentNil(b bool)`

 SetContentNil sets the value for Content to be an explicit nil

### UnsetContent
`func (o *BannerIdDeleteRequest) UnsetContent()`

UnsetContent ensures that no value is present for Content, not even an explicit nil
### GetIsActive

`func (o *BannerIdDeleteRequest) GetIsActive() bool`

GetIsActive returns the IsActive field if non-nil, zero value otherwise.

### GetIsActiveOk

`func (o *BannerIdDeleteRequest) GetIsActiveOk() (*bool, bool)`

GetIsActiveOk returns a tuple with the IsActive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsActive

`func (o *BannerIdDeleteRequest) SetIsActive(v bool)`

SetIsActive sets IsActive field to given value.

### HasIsActive

`func (o *BannerIdDeleteRequest) HasIsActive() bool`

HasIsActive returns a boolean if a field has been set.

### SetIsActiveNil

`func (o *BannerIdDeleteRequest) SetIsActiveNil(b bool)`

 SetIsActiveNil sets the value for IsActive to be an explicit nil

### UnsetIsActive
`func (o *BannerIdDeleteRequest) UnsetIsActive()`

UnsetIsActive ensures that no value is present for IsActive, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


