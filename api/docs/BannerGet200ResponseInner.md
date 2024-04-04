# BannerGet200ResponseInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BannerId** | Pointer to **int32** | Идентификатор баннера | [optional] 
**TagIds** | Pointer to **[]int32** | Идентификаторы тэгов | [optional] 
**FeatureId** | Pointer to **int32** | Идентификатор фичи | [optional] 
**Content** | Pointer to **map[string]interface{}** | Содержимое баннера | [optional] 
**IsActive** | Pointer to **bool** | Флаг активности баннера | [optional] 
**CreatedAt** | Pointer to **time.Time** | Дата создания баннера | [optional] 
**UpdatedAt** | Pointer to **time.Time** | Дата обновления баннера | [optional] 

## Methods

### NewBannerGet200ResponseInner

`func NewBannerGet200ResponseInner() *BannerGet200ResponseInner`

NewBannerGet200ResponseInner instantiates a new BannerGet200ResponseInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBannerGet200ResponseInnerWithDefaults

`func NewBannerGet200ResponseInnerWithDefaults() *BannerGet200ResponseInner`

NewBannerGet200ResponseInnerWithDefaults instantiates a new BannerGet200ResponseInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBannerId

`func (o *BannerGet200ResponseInner) GetBannerId() int32`

GetBannerId returns the BannerId field if non-nil, zero value otherwise.

### GetBannerIdOk

`func (o *BannerGet200ResponseInner) GetBannerIdOk() (*int32, bool)`

GetBannerIdOk returns a tuple with the BannerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBannerId

`func (o *BannerGet200ResponseInner) SetBannerId(v int32)`

SetBannerId sets BannerId field to given value.

### HasBannerId

`func (o *BannerGet200ResponseInner) HasBannerId() bool`

HasBannerId returns a boolean if a field has been set.

### GetTagIds

`func (o *BannerGet200ResponseInner) GetTagIds() []int32`

GetTagIds returns the TagIds field if non-nil, zero value otherwise.

### GetTagIdsOk

`func (o *BannerGet200ResponseInner) GetTagIdsOk() (*[]int32, bool)`

GetTagIdsOk returns a tuple with the TagIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTagIds

`func (o *BannerGet200ResponseInner) SetTagIds(v []int32)`

SetTagIds sets TagIds field to given value.

### HasTagIds

`func (o *BannerGet200ResponseInner) HasTagIds() bool`

HasTagIds returns a boolean if a field has been set.

### GetFeatureId

`func (o *BannerGet200ResponseInner) GetFeatureId() int32`

GetFeatureId returns the FeatureId field if non-nil, zero value otherwise.

### GetFeatureIdOk

`func (o *BannerGet200ResponseInner) GetFeatureIdOk() (*int32, bool)`

GetFeatureIdOk returns a tuple with the FeatureId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeatureId

`func (o *BannerGet200ResponseInner) SetFeatureId(v int32)`

SetFeatureId sets FeatureId field to given value.

### HasFeatureId

`func (o *BannerGet200ResponseInner) HasFeatureId() bool`

HasFeatureId returns a boolean if a field has been set.

### GetContent

`func (o *BannerGet200ResponseInner) GetContent() map[string]interface{}`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *BannerGet200ResponseInner) GetContentOk() (*map[string]interface{}, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *BannerGet200ResponseInner) SetContent(v map[string]interface{})`

SetContent sets Content field to given value.

### HasContent

`func (o *BannerGet200ResponseInner) HasContent() bool`

HasContent returns a boolean if a field has been set.

### GetIsActive

`func (o *BannerGet200ResponseInner) GetIsActive() bool`

GetIsActive returns the IsActive field if non-nil, zero value otherwise.

### GetIsActiveOk

`func (o *BannerGet200ResponseInner) GetIsActiveOk() (*bool, bool)`

GetIsActiveOk returns a tuple with the IsActive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsActive

`func (o *BannerGet200ResponseInner) SetIsActive(v bool)`

SetIsActive sets IsActive field to given value.

### HasIsActive

`func (o *BannerGet200ResponseInner) HasIsActive() bool`

HasIsActive returns a boolean if a field has been set.

### GetCreatedAt

`func (o *BannerGet200ResponseInner) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *BannerGet200ResponseInner) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *BannerGet200ResponseInner) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *BannerGet200ResponseInner) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *BannerGet200ResponseInner) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *BannerGet200ResponseInner) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *BannerGet200ResponseInner) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *BannerGet200ResponseInner) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


