import base64
import sys

"""
About:
  This script decodes a given base64-encoded image content
  and decodes it to a valid image file
"""

if len(sys.argv) != 3:
  print("Usage: ./decode-image.py <EncodedDataFilename> <OutputFilename>")
  exit()

with open(sys.argv[1], "r") as f1:
  data = f1.read()
  with open(sys.argv[2], "wb+") as f2:
    f2.write(base64.b64decode(data))
    f2.close()
  f1.close()
