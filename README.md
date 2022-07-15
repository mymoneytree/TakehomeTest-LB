## Moneytree Takehome Exercise

### Summary

Send, Receive and Order Messages from Multiple Producers over HTTP

Your task is to write 1) text-processing script 2) a "producer" and 3) an http server. The 2) producer will send http requests to a server (described below) while the 3) server will handle responses from that Encoding/Decoding server and output the results to stdout. 

Please feel free to use any programming language to complete these tasks. Each component can be written in a different language if desired. 

### Decoding Server 

Moneytree will provide you with a binary named `decode-server` and a newline delimited file of encoded words.

The binary is an HTTP server that is listening on port 5555. This server exposes one route - `http://localhost:5555/event`.

This route accepts `POST` requests and is expecting a payload in the format of:

```
{
  "encoded": <string>,
  "meta": {...}
}
```

Where "encoded" is a string and "meta" is an arbitrary JSON blob. 

This server will then decode the string and send a subsequent http POST request to `http://localhost:5556/event`. It will also send the `meta` object with no modifications. 

The decode server will process and respond to incoming messages in an arbitrary order. 

### Pre-processing Script

Write a script that can read the initial file and add any required metadata necessary to complete the task. This script should take an `<int> N` as a command line argument and output `N` separate files that evenly divide the initial set of words in the master file.

The contents of each string should not be modified, but additional data CAN be added to this list of encoded strings.

### Producer

Write a program that reads in a list of encoded words from the output of the pre-processing script. This can be done via command line args, reading from file, etc.

This program should then send each string (and whatever additional metadata you desire) to the Decoding Server. 

You should be able to run N instances of this program in parallel, where N is the number of output files created by the pre-processing script.

### Server

Write a server that listens on port `5556` with a route named `/event` that can receive POST requests over http from the `decode-server`.

This server should expect a JSON payload of:

```
{
  "decoded": <string>,
  "meta": {...}
}
```

Where the decoded message is the decoded version of the initial string and the meta object is an unmodified copy of the `meta` object from the initial message. 

This server should reorder the messages received from the `decode-server` to print them in the original order as written in the master CSV (before splitting). The server should be robust to dropped messages from the `decode-server` and always produce an output. 

### Test/Success Criteria

The initial list of words should be arbitrarily subdividable and an arbitrary number of Producers should be able to send the full set of words (and metadata) to the Decode-Server. 

Your server should listen at `http://localhost:5556/event` and be able to receive the messages in an arbitrary order and print them in the correct order (the original order in the master CSV list).

### Deliverables

1) Producer that can read input and send messages over http to the Decode Server

2) Server that can receive messages on port `5556` and re-order messages received to print them in the original order

3) A script to process the original data file and transform it with any necessary metadata to complete the task.

4) Instructions on how to build, run and orchestrate 1,2 & 3 to acheive the desired results

### Evaluation Criteria
- Code Clarity, Naming
- Architecture, Structure, Separation, Reusability
- Implementation, Functionality
- Consideration of Edge Cases
- Documentation