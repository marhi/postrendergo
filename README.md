Postrender is a 2 component pre-rendering system.
It's called postrender just to make it a bit less boring ;)

The first part of the system is the js file (generated from TypeScript), which should be included into the web page you wish to postrender.
The script waits for the DOM to be available or mutated then after a fixed time delay it sends it back to the backend.

The second part is the backend server written in GO, which accepts the renered HTML and saves it into a directory.
It also generates the .js file for inclusion based on given parameters.

Sample usage [output the .js file and start listening]:

mkdir output
go run main.go -delay 2000 -host 'http://localhost:1337/renderer' -listen :1337 -auth 'lalal' -dir output > ~/angular5/shop/src/assets/postrender.js

Licensed under Apache v2.0