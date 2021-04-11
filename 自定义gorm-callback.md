# Customlize your gorm-Callbacks

GORM 本身是由回调驱动的，所以我们可以根据需要完全定制 GORM，以此达到我们的目的，如下：

- 注册一个新的回调
- 删除现有的回调
- 替换现有的回调
- 注册回调的顺序

在 GORM 中包含以上四类 Callbacks，我们结合项目选用 “替换现有的回调” 来解决一个小痛点。

## 问题

在 models 目录下，我们包含 tag.go 和 article.go 两个文件，他们有一个问题，就是 BeforeCreate、BeforeUpdate 重复出现了，那难道 100 个文件，就要写一百次吗？

1、tag.go

![图片](https://mmbiz.qpic.cn/mmbiz_jpg/KVl0giak5ib4jMEmHXYhh4gXuNfjAmsN4peIibMcjiaTicEXRTlQZIySz35icH8r3bJWib7PcibvnuaOVm9rDJIhM0XMqw/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

2、article.go

![图片](https://mmbiz.qpic.cn/mmbiz_jpg/KVl0giak5ib4jMEmHXYhh4gXuNfjAmsN4pruVEb5pVRcodveoqlQOcS8cDF9aqldYhacMA8y5Qicfy2N8f0WvBxdQ/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)



显然这是不可能的，如果先前你已经意识到这个问题，那挺 OK，但没有的话，现在开始就要改

## 解决

### 实现Callbacks

1、updateTimeStampForCreateCallback

打开models目录下的models.go文件,添加如下代码：

```go
func updateTimeStampForCreateCallback(scope *gorm.Scope){
	if !scope.HasError(){
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok{
			if createTimeField.IsBlank{
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok{
			if modifyTimeField.IsBlank{
				modifyTimeField.Set(nowTime)
			}
		}
	}
}
```

在这段代码中，会完成以下功能：

- 检查是否有含有错误（db.Error）
- `scope.FIeldByName`通过`scope.FIelds()`获取所有字段，判断当前是否包含所需字段

```go
for _, field := range scope.FIelds(){
    if field.Name == name || fields.DBName == name{
        return field, true
    }
    if field.DBName == dbName{
        mostMatchedField = field
    }
}
```

- `field.IsBlank`可利用反射判断该字段的值是否为空

```go
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
```

2、updateTimeStampForUpdateCallback

接着写入如下代码：

```go
func updateTimeStampForUpdateCallback(scope *gorm.Scope){
	if !scope.HasError(){
		nowTime := time.Now().Unix()
		if _, ok := scope.Get("gorm:update_column"); ok{
			scope.SetColumn("ModifiedOn", nowTime)
		}
	}
}
```

- `scope.Get(...)`根据入参获取设置了字面值的参数，例如本文是`gorm:update_column`，它会去查找含这个字面值的字段属性
- `scope.SetColumn(...)`假设没有指定`udpate_column`的字段，我们默认在更新回调设置`ModifiedOn`的值

### 验证：

addTags：

```go
http://127.0.0.1:8000/api/v1/tags?name=math&state=0&created_by=songpeng&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QiLCJwYXNzd29yZCI6InRlc3QxMjM0NTYiLCJleHAiOjE2MTgxNjg1MjcsImlzcyI6Imdpbi1ibG9nIn0.8wwXMidpJ-XKAtDK3YougTHX47YHIwGQW9IX_UbWyD4
```

结果如下：

![image-20210412002324159](自定义gorm-callback.assets/image-20210412002324159.png)

![image-20210412002609046](自定义gorm-callback.assets/image-20210412002609046.png)