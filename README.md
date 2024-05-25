# Sygo
Sygo is an image steganography CLI tool. It can be used to encode data into an image and decode data it has encoded in an image.

When Sygo encodes a data string into an image, it also encodes a 32bit integer header into it, specifying the byte length of the data string. This is used to know how many bytes of data to read when decoding the image.

## Synopsis
```sygo [--output] [--channels] [--bits-per-channel] <command> <args>```

Example:  
```sygo encode ./input.png "the quick brown fox jumped for the lazy dog"```

## Commands
```encode <filepath> <datastring>```  
Encodes a data string in an image file and saves it as a copy.

```decode <filepath>```  
Decodes a data string from the given file.

```capacity <filepath>```  
Calculates the storage capacity in bits of the given file.

```size <datastring>```  
Calculates the payload size to encode, resulting from the given data string and header (data string bit size + 32bit header).

      
## Flags
```--output    (Default value: encoded.png)```  
   Sets the filepath for the encoded image.  
   Compatible commands: encode
   
```--channels    (Default value: rgb)```  
   Set which channels are used when encoding/decoding.  
   Compatible commands: encode, decode, capacity
   
```--bits-per-channel    (Default value: 1)```  
   Sets how many least-significant-bits are used per channel when encoding/decoding.  
   Compatible commands: encode, decode, capacity

## Limitations/Caveats

- Has only been tested on .png files
- Only supports 8bit color channels.
