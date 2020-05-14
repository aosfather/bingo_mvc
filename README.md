# bingo_mvc
简单高效的mvc框架，支持go自带的http库，和性能卓越的fasthttp库实现。  


# 特性
* 双http引擎
* 支持GET、POST多种方式访问
* 支持json参数
* 支持xml参数
* 支持序列化成json、xml格式
* 支持返回modeview格式，自动调用模板进行渲染
* 支持静态文件处理
* 提供Inject tag进行自动装载
* 提供Value tag进行属性从配置文件中自动赋值
* 提供sqltemplate实现，不用写sql也能增删改查
* 提供本地和分布式session实现方式
* 提供多request mapper的方式
* 提供自动将请求参数赋值给方法参数
* 提供Field tag，来指明对应的输入参数名
* 提供拦截器扩展，可以在服务执行前后进行响应处理
* 提供了Boot方式启动，通过加载配置文件，将引擎配制好。
* Boot中提供简单的IOC容器，自动装载对象和赋值


# 样例
### 简单的例子
```go

func main(){
  f:=fasthttp.FastHTTPDispatcher{}
  f.Port=8080
  f.Run()

}
```
### 来个有点意思的
```go
import (
	"fmt"
	"github.com/aosfather/bingo_mvc"
	"log"
)
//方法需要的输入参数
type MyRequest struct {
	Name string `Field:"name"`  //指定参数输入名称是name
}
//主要的服务提供者，提供了两个url，一个/test,一个/test1
type MyHandle struct {
	Test  string `mapper:"name(test);url(/test);method(GET);style(HTML)"`
	Test1 string `mapper:"name(test1);url(/test1);method(GET);style(JSON)"`
}
// 这个框架会调用的方法，通过这个方法返回了对应url响应的方法及对应的参数对象
func (this *MyHandle) GetHandles() bingo_mvc.HandleMap {
	result := bingo_mvc.NewHandleMap()
	r:=&MyRequest{}
	result.Add("test",this.DoTest,r)
	result.Add("test1",this.DoTest1,r)
	return result
}

func (this *MyHandle) DoTest(a interface{}) interface{} {
	t:=a.(*MyRequest)//框架会按指定的参数类型，进行赋值回调
	log.Println(t.Name)
	return "hello"
}

func (this *MyHandle) DoTest1(a interface{}) interface{} {
	t:=a.(*MyRequest)
	log.Println(t.Name)
	return fmt.Sprintf("hello %s",t.Name)
}
func main() {
	h := HttpDispatcher{}
	h.Port = 8090
    //向dispatch注册url mapping信息，简单明了
	h.AddRequestMapperBystruct(&MyHandle{})
	h.Run()
}

```
使用 curl localhost:8090/test?name=xxxx,试验一下吧。两个例子用的dispatcher类不一样， 
只是为了演示了下双引擎，两者的方法是一样的，外部行为没有什么不一样。
### 基于boot来启动应用
```go
import "github.com/aosfather/bingo_mvc/context"
//MyRequest 和MyHandle和上面一样，只是来了点小魔法，自动装载数据库链接
type MyHandle struct {
	Test  string `mapper:"name(test);url(/test);method(GET);style(HTML)"`
	Test1 string `mapper:"name(test1);url(/test1);method(GET);style(JSON)"`
    DB *sqltemplate.Datasource `Inject` //系统将自动装载
}
//演示了使用数据库，至于需不需要强制分 controller、service等的方式
//仁者见仁，智者见智了。根据系统复杂程度来定吧。
func (this *MyHandle) DoTest(a interface{}) interface{} {
	t:=a.(*MyRequest)//框架会按指定的参数类型，进行赋值回调
	log.Println(t.Name)
    dao:=this.DB.GetDao()
    dao.FindBySql(t,"select name from mytest1 where id=123")
    log.Println(t.Name) //看看是不是不一样了
	return "hello"
}
func main(){
    boot:=context.Boot{}
    boot.Init(&HttpDispatcher{},onLoad)
    boot.Start()
}

func onLoad()[]interface{}{
    return []interface{}{&MyHandle{}}
}
```
boot会自动加载bingo.yaml文件，该文件放在应用的根目录下
```yaml
# bingo.yaml
bingo：
  port: 8080   #监听的端口
  static: statics #静态资源的目录，可以没有
  template： templates #模板目录
  template_fix： tp #模板文件后缀
  useddB：true   # 启用数据库，可以不启用
  db：
     type: mysql      #数据库类型
     name: test       # 数据库名
     url: localhost:3306 #服务器地址及端口
     user: dev          #数据库用户名
     password: dev      #数据库用户对应的密码
  
```
使用boot后就拥有了一个ioc容器，使用Inject tag，会自动装载，可以指明对象的名称。
使用Value tag会将配置中对应的属性值赋值。



