### echo-logrus

Middleware echo-logrus is a [logrus](https://github.com/sirupsen/logrus) logger support for [echo](https://github.com/labstack/echo).

This fork supports echo V4.

#### Install

```sh
go get -u github.com/davrux/echo-logrus/v4
```

#### Usage

import package

```go
echologrus "github.com/davrux/echo-logrus/v4"
```

define new logrus

```go
echologrus.Logger = logrus.New()
e.Logger = echologrus.GetEchoLogger()
e.Use(echologrus.Middleware())
```
