threepio
==

![threepio](https://s-media-cache-ak0.pinimg.com/736x/c5/35/c9/c535c913ca0135bd19010f013a7e65f6.jpg)

threepio is a little golang app. It responds to a URI string dictating where an asset for editting is stored, syncs this and opens either prelude or premiere for editting.

Grammar
--

```
uri             = "threepio+", application, "://", path, "?", params;
application     = "prelude" | "premiere";

path            = "/", alphanumeric, {path};

params          = param, {"&", params};

param           = key, "=", value;
key             = "accessKey" | "secretKey" | "sessionToken" | "uuid";
value           = alphanumeric;


alphanumeric    = letter | digit, {alphanumeric | punctuation};
```

Build
--

Golang
--

OSX
---
```brew update```
```brew install go --with-cc-all```
```brew reinstall go --cross-compile-all```

Check the go version
----
```go version```
```go version go1.7.3 darwin/amd64```


Dependencies
---

Get the dependencies
```
go get
```

Compile and build
Build binary
```make dist```

Running
---
```
go run threepio.go -f=config/threepio.ini -u=threepio+prelude:///test?uuid=123456

-f configuration file
-u passed uri
```

or from binary

```./threepio -f=config/threepio.ini -u=threepio+prelude:///test?uuid=123456```

Packaging
---

Install dmg builder
----
```npm install -g appdmg```

Build the dmg
```appdmg dmg/threepio-dmg.json dmg/target/threepio.dmg```


Registering
---


