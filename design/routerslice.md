### routeslice
中间件的设计思想: 本项目中使用Slice来顺序存放需要执行的中间件和业务逻辑函数, for循环遍历该Slice中的方法，顺序地执行函数。主要实现类:

- RouterHandler       

- RouterSlice
  - Use()
  - AddHandle()
- GroupRouter
  - Use()
  - AddHandle()

需要注意的是，对于每个msgID, AddHandle()方法只能执行一次，否则会panic（详见源码）。如果需要添加多个函数，在执行AddHandle()时，传入多个RouterHandler