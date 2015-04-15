# img-debug-proxy
Proxy in Go to help web-developers to deal with images 

Current version is a proof of concept. Proxy handles popular image formats and adds a label "\<width\>x\<height\>" at its top-left corner.

Planned features:

* Parse a special image's URL parameter like "&xxx800x600" - expected dimentions. Check an actual image dimention and mark image correspondingly. 
Generate empty image in case GET fails

* web-admin UI, statistics, etc.