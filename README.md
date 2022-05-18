# GRPC-Gateway? Why not?

I struggled for some time with this decision.
I tried GRPC-Gateway, and I wouldn't say I liked it.
Here is why.

Let me start by saying that I like the idea of code generation: I did work on this topic during my master thesis.
Using GRPC-Gateway, I could specify the application's input format, output format, and interaction logic only once.
Consequently, every time I change the proto specification, I only have to generate the code.

Nice, but how much does it cost?

Code generation is amazing, but only until you are satisfied with the behaviors of the generated code.

Let's analyze the redirection scenario.

The GRPC server answers with a GetRedirectionLocationResponse, but I want the Gateway to redirect the user to the speficied location instead of returning such a location a JSON response.

To the bost of my knowledge, I cannot express a redirect using `google.api.annotations` in the proto file, thus I have to implement some custom beheaviour on the gateway side, intercepting a specific type of response to change the behaviour, like in the code below:

```go
func responseMessageMatcher(ctx context.Context, w http.ResponseWriter, resp protoreflect.ProtoMessage) error {
    // Check if the response is of type GetRedirectionLocationResponse
    t, ok := resp.(*pb.GetRedirectionLocationResponse)
    if ok {
        // Set the Location header and return a StatusFound http code
        w.Header().Set("Location", t.Location)
        w.WriteHeader(http.StatusFound)
    }
    return nil
}
```

That isn't lovely, in my honest opinion.
One could say *"but it is just for this time"*, but I think that the effort of building a custom REST API is too low to justify such a choice.