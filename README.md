# iRecruit GraphQL API

This is a graphql API for the iRecruit web app, setup to run with
a MongoDb backend. And a "demo-able" internal MongoDb mock, for the
basic functionality, i.e. no complex, `$gt` and `$lt` style MongoDb
queries, or `$where` or anything like that.

Still gonna put docs for the API usage.

The API routes are functionally tested, and all the systems are unit tested..
Well, there's really only one system, i.e. the CRUD abstraction.

## Some Excuses
Uh, yea, so I guess first things first, setup.

Steps are still being defined, because there are still some dependancies
like platforms and needing the "make" tool already installed, which is a
convenience. Otherwise you can just run the defacto `go build`. 


There are also some 1 or 2 shell commands in the `Makefile` which might make it difficult to run on Windows without tweaking.
But if Windows is what you live for, this shouldn't be a problem.

I'll look into setting up a docker image a little bit later on, though there is
a `Dockerfile` for building and running the code, if you wanna use that. I haven't
read up enough on that, and just know the bare minimum to setup this Dockerfile.

## Setup

Now there are a couple of ways to do this.


If you're gonna build or run from source first you need to "go get" the dependancies:

    go get -u ./...	


### Mongo-Mode from Source

- Start your mongodb server

- Copy example.env to .env.development for development or .env.production for production

- Configure the .env file if necessary... probably is.. maybe

- Run `make`
#
    make run



### Mock-Mode from Source

No setup at all really, just run:

    make mock


### Binary-Mode
Somehow get a copy of the executable binary, and bam, just set your env, usually either development, test or production.

    export ENV=development

The run your binary:

    ./binary

### Mock-Mode from Binary

Just run your binary with  the `-mock` flag.

	./binary -mock


#### Mongo-Mode from Binary

- startup your mongodb server

- configure your .env file

- run the binary

#

	./binary -mock


## Testing

... Uhm, yea, I guess just run `go test` on everything or whatever. I'm still figuring things out here. Still need to setup some things. But for now, if you're like me, and just wanna know if everything works according to the tests just run:

    make test

That ought to see you through for the most part.