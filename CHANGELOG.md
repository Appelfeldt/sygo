# Changelog

## 0.1 - 2024-05-26
First release of this tool, expect bugs-a-plenty.

- Sygo v0.1 has only been tested on .png files, do not expect the tool to function properly for other formats.
- Command: encode for encoding data into an image
- Command: decode for decoding data from an image
- Command: capacity for calculating storage capacity in bits of an image
- Command: size for calculating the bit size of a data strings resulting payload (data string bit size + 32bit header)
- Flag: output to set the filepath of the encoded image
- Flag: channels to set which channels are used when encoding/decoding
- Flag: bits-per-channel for how many least-significant-bits to use per channel when encoding/decoding
