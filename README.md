# Background-Removal API
> Removes the background from a given image file.

## Endpoints:
📔 Note: All endpoints accept ``application/json`` only.

---

**POST /api/remote**

Receives a valid json payload containing the fields:
* ``filename`` - was only tested for `JPEG`, `JPG`, `PNG`.
* ``uri`` - where the image is stored at the public internet.

Returns as json:
* ``filename`` - the exact same (given) filename.
* ``content`` - the result image(background removed) content in base64.

---

**POST /api/bs64**

Receives a valid json payload containing the fields:
* ``filename`` - was only tested for `JPEG`, `JPG`, `PNG`.
* ``content`` - the base64-encoded image content.

Returns as json:
* ``filename`` - the exact same (given) filename.
* ``content`` - the result image(background removed) content in base64.

---

**POST /api/file**

Receives a valid json payload containing the fields:
* ``filename`` - was only tested for `JPEG`, `JPG`, `PNG`.
* ``content`` - the raw binary image content.

Returns as json:
* ``filename`` - the exact same (given) filename.
* ``content`` - the result image(background removed) content in base64.

---
<br>

## Example:
Let us demonstrate the ``/api/remote`` endpoint (other endpoints work in a similar fashion).

We want to remove the background of this gorgeous looking duck image file(`jpg`):
![Duck with Background](https://upload.wikimedia.org/wikipedia/commons/b/bf/Bucephala-albeola-010.jpg)
So we make the following HTTP request:
```HTTP
POST /api/remote HTTP/1.1
Host: 127.0.0.1:3000
Content-Type: application/json
Content-Length: 173

{
	"uri":	"https://upload.wikimedia.org/wikipedia/commons/b/bf/Bucephala-albeola-010.jpg",
	"filename":	"duck.jpg"
}
```
Now, the server response is described as mentioned at the *Endpoints*.
The value of the `content` field (a base64-encoded string representing the image with it's background removed) is huge, thus I won't put it in here.
But after decoding the image and writing it as binary to a `jpg` file, we get (using `/tests/decode-image.py`):
![Duck without Background](/tmp/output_duck.jpg)
