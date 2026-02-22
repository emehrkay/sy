# SY Monitor App

## Design

I took a hexagonal approach when designing this service. It is split up into three parts: `storage`, `service`, and `implementation` layers. Each layer provides a public interface, concrete instances of that interface, and its own types relevant to that layer (they're all the same, but in a real application it would be differences).

## Storage Layer

This layer handles storing the stats and heartbeats for a given device. The implementation of this is an in-memory set slice of records. I built I in a way that would semi-mimic a real datastore where you have primary key (deviceID) and a collection of records (Stats, Heartbeats) associated with it. An actual database could be dropped in if it adheres to the `storage.Storage` interface.

## Service Layer

This is the main brains of the application. In a complete application, it would be the only way to interact with the storage. This creates a set of business rules that can be executed from many different contexts: http, grpc, cli, etc.

## Implementation Layers

There are two ways to interact with this application: http and cli. For applications like this, I like to provide a single executable with multiple entry points. This would allow for an easy way to interact with the app and its data without having to navigate the full http or rpc auth stack. I like to think of this layers as simply a mechanism to build requests for the service layer and translate responses to its callers.

### Dependencies

* `github.com/spf13/cobra` -- It makes build CLIs a bit easier

### Omitted for Brevity

* Observability -- In a real application I would create custom syslog handlers and log output of errors and traces for performance-related metrics
* Tests -- I structured this code in a way where every layer is independly testable, allow for any of its dependencies to be mocked and injected where necessary. A tool like GoMock or even a custom in-house mocking tool would do the job just fine.
* Custom errors -- I love creating both sentinel and custom error types. It makes it easier to normalize at the public interface level (convert an ErrNotFound to 404 for example)
* Edge cases -- Since I used the device simulator as my test bed, I wasn't quite able to harden the code against some corner cases like when a device is considered offline and see now that affects the numbers
* Context cancellations -- In a real application I would use the context to handle timeouts and cancellations etc

## Runninng the Application

Running the application can be done by building it `go build cmd/main.go` and running commands against the executable or by running the same commands during `go run cmd/main.go`.

```
go run cmd/main.go apistart --csv __meta/devices.csv
```

> this starts the api server and automatically loads the devices.csv file

There is also a Makefile with one target: `run_api` that will run the api server

## Wrap Up

Overall it took me maybe five hours to develop this code base. I wanted to focus on making sure it was a solid foundation that anyone can build upon and swap/upgrade components when needed. I think that the most time spent was the storage layer, I do some calculations in there -- `AverageUploadTime` `MinuteRange` -- that could otherwise be computed in the actual storage engine. 

The application isn’t too complicated so its time and space complexity should be easy to gauge. The memory storage keeps a map of device ids related to a slice of either `Stats` or `Heartbeat` so it would be the number of devices and then the number of each stat or heartbeat O(D + S + H). During the get stats process, both the Stats and Heartbeats are reordered, O(LogN), heartbeats take the first and last elements for its calculations O(N), Average UploadTime looks at every Stat O(N)

No AI was used to write this. 

I would add security on both the implementation and the service layers. The Implementation layer would use the existence cookies or jwts to restrict access and pass it to the Service layer to verify that the user making the request is valid and allowed to perform the requested action.
