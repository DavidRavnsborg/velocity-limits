# Overview

This is my coding challenge for my KOHO Backend Developer application! This is the first app I have ever written in Go. I found Go to be relatively easy to learn, given its relatively small feature set. I also really appreciate working in a language with static type checking.

The application loads the requests from the input.txt file and outputs the resolved requests to output-test.txt. If you run `go test`, it also runs generates the file main_test.go tests that the final `output-test.txt` and `output.txt` files match.

I started off by planning how I would approach the problem. Initially, I was going to make a more complex version of this app that could have a input/state/output functions that would each accept an interface/interfaces. The idea was that you could, via dependency injection, provide the desired implementation of each part (e.g. in-memory database for state, vs. an actual persistent database vs. a filesystem based solution). However, I quickly realized how large a task this would be, especially in a new language, and instead opted to focus on the core design with good testing. Of input, state, and output, I only used interfaces on the state layer, with two interfaces: ResponseLimits and VelocityLimits. If you wanted to change the implementation of my ad-hoc in-memory database and table (FundSuccessDB and ResponsesTable) to real databases, you could create new types that implement ResponseLimits and VelocityLimits, that manage and fetch data from real databases. After creating the new interface implementations, you wouldn't need to modify anything else in the application except for the fields you pass to the checkConditions function for the interfaces ResponseLimits and VelocityLimits. 

The resulting app loads the data into a custom array of structs using JSON de-serialization of the individual request strings. Those form the input Data Transfer Objects (DTOs), or request objects, which then get resolved into response objects. I purposely set it up so that it would be relatively easy to change the input/output layers to an HTTP server that receives and acts on individual requests. Although on an HTTP server you would make use of the handleRequest receiver method of FundRequest, rather than handleBatchRequestsWriteToFile which was written specifically to operate on batches of requests (in ascending time order).


# How-to-run

Assuming you have Go already installed, to run the application, call `go run .`. You can see my Go version in go.mod, if you need to check for compatibility. It should generate an `output-test.txt` file, that is identical to `output.txt`.

To Test the application, call `go test`, which will run tests in main_test.go, state_memory_db_successes.go, and state_memory_table_responses.go. main_test.go should also generate a new `output-test.txt` file, that is identical to `output.txt` (the test will fail if it isn't).


# My Development Notes

## April 14, 2021
line 687 of the output-test.txt has an extra line with the value: {"id":"6928","customer_id":"562","accepted":true}
It is the only discrepancy between output.txt and output-test.txt. There appears to be a duplicate id on the same customer, so see why it wasn't filtered out.
FIXED: I needed another table, the Responses table, since non-approved responses aren't logged in the Success table (which is used for checking velocity limits on amounts/transactions).

I finished adding tests for the ResponsesTable and FundSuccessTable. I also fixed the length of the history for weekly limits (it should only go back to Monday of the current week at 00:00:00).

## April 13, 2021
Create only one version given the time constraints, but design it using an interface that can be used by other versions for the persistence layer (I.e. if the persistence layer was a DB, the Go wrapper for it could implement the Limits interface). That will accomplish the extensible design I'm going for, while keeping complexity low, as this project is already taking quite a bit of my afterwork time. I'd rather polish the core design and testing of this app than try to cram more features into it.

## April 11, 2021
Completed most of Stephen Girder's "Go: The Complete Developer's Guide (Golang)" on Udemy over the weekend. Start the implementation now by starting with the simplest version of the application.

Make it runnable with different command line args (simplest version first - easiest to implement):

### EDIT: following the schema below, I actually implemented "local_file_memory_file" (not listed), since my completed solution runs locally (not deployed anywhere), loads data from a local file (input.txt), persists state in memory, and outputs to a local file (output-test.txt).

  `x. "{where it runs} _ {data input} _ {state persistence} _ {data output}"`
  1. "local_file_memory_std"        Runs locally, loading data from local file, persisting state in memory, outputting in std io. 
  2. "local_file_file_file"         Runs locally, loading data from local file, persisting state in local files, outputting in file.
  3. "local_http_file_http"         Runs locally, receiving data in HTTP request, persisting state in local files, outputting in HTTP response.
  

These can be made more customizable with a config file. The config file can be written in Go, to allow for easy dependency injection, which would also make it easily configurable, particularly if I use seemingly common built-in Go interfaces like Reader and Writer. I will put the business logic in a separate package/module. 

Read: https://medium.com/rungo/anatomy-of-modules-in-go-c8274d215c16

These would rightfully be interpreted as scope creep, but otherwise good ideas:
- `4. "gae_http_cloudsql_http"       Runs on Google App Engine, receiving data in HTTP request, persisting state to Cloud SQL DB, outputting by HTTP response.`
  - Out-of-scope, because while I might be able to deploy a barebones GAE Go app in time, I wouldn't be able to spend the time necessary to secure it; and the assignment asked for no public hosting of the solution. I will assume that extends to the deployed solution, not just source code.
- Have a config to run the application as a web app, and create another web app to send the original web app the input (a simple version of this could make a good test case for the local_http_file_http scenario though).
- Purge db data older than a month (possibly with Cron job).
  - If I am getting into lifecycle management of data, I should be asking questions like "how long should my data persist for business, or regulatory reasons? Should there be backups? Etc."; but that level of detail doesn't exist for this problem. However, this would be a question to ask in a follow-up meeting with the client, product owner, lawyer, or stakeholder.


# Problem Statement

In finance, it's common for accounts to have so-called "velocity limits". In this task, you'll write a program that accepts or declines attempts to load funds into customers' accounts in real-time.

Each attempt to load funds will come as a single-line JSON payload, structured as follows:

```json
{
  "id": "1234",
  "customer_id": "1234",
  "load_amount": "$123.45",
  "time": "2018-01-01T00:00:00Z"
}
```

Each customer is subject to three limits:

- A maximum of $5,000 can be loaded per day
- A maximum of $20,000 can be loaded per week
- A maximum of 3 loads can be performed per day, regardless of amount

As such, a user attempting to load $3,000 twice in one day would be declined on the second attempt, as would a user attempting to load $400 four times in a day.

For each load attempt, you should return a JSON response indicating whether the fund load was accepted based on the user's activity, with the structure:

```json
{ "id": "1234", "customer_id": "1234", "accepted": true }
```

You can assume that the input arrives in ascending chronological order and that if a load ID is observed more than once for a particular user, all but the first instance can be ignored. Each day is considered to end at midnight UTC, and weeks start on Monday (i.e. one second after 23:59:59 on Sunday).

Your program should process lines from `input.txt` and return output in the format specified above, either to standard output or a file. Expected output given our input data can be found in `output.txt`.

You're welcome to write your program in a general-purpose language of your choosing, but as we use Go on the back-end and TypeScript on the front-end, we do have a preference towards solutions written in Go (back-end) and TypeScript (front-end).

We value well-structured, self-documenting code with sensible test coverage. Descriptive function and variable names are appreciated, as is isolating your business logic from the rest of your code. For example, consider decoupling input/output or data storage such that changing the underlying implementation wouldn't change how it is used from a client.

Thanks for your interest in KOHO - have fun!
