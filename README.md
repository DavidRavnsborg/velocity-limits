# Overview

# How-to-run

Run with `go run .`

Test with `go test`

# My Development Notes

## April 14, 2021
line 687 of the output-test.txt has an extra line with the value: {"id":"6928","customer_id":"562","accepted":true}
It is the only discrepancy between output.txt and output-test.txt. There appears to be a duplicate id on the same customer, so see why it wasn't filtered out.
FIXED: I needed another table, the Responses table, since non-approved responses aren't logged in the Success table (which is used for checking velocity limits on amounts/transactions).

I finished adding tests for the ResponsesTable and FundSuccessTable. I also fixed the length of the history for weekly limits (it should only go back to Monday of the current week at 00:00:00).

## April 13, 2021
Create only one version given the time constraints, but design it using an interface that can be used by other versions for the persistence layer (I.e. if the persistence layer was a DB, the Go wrapper for it could implement the Limits interface). That will accomplish the extensible design I'm going for, while keeping complexity low, as this project is already taking quite a bit of my afterwork time. I'd rather polish the core design and testing of this app than try to cram more features into it.

## April 11, 2021
Completed Go crash course. Start by creating the simplest version of the application.

Make it runnable with different command line args (simplest version first - easiest to implement):

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
