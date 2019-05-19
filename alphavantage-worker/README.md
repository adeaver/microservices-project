# Alpha Vantage Worker

This worker is responsible for getting daily stock data for a number of symbols. It will use Alpha Vantage's API to collect this data.

## Design Decisions

I went back and forth on the design for this. The main consideration I had to make was whether I wanted a common Alpha Vantage client that could be shared between my different services or if I wanted to encapsulate the logic inside of only the services that needed it. The main options were:

#### Why I Didn't Choose a Shared Client

This made the most sense from a code reusability standpoint, which may or may not have been useful for a project this small, since I only really intended to use Alpha Vantage for this one worker. I actually went down the thought process of what I would want if I were working with thousands of microservices. I decided that a shared client would actually probably end up being somewhat untenable. The shared client would have to maintained and modified by a bunch of different teams, each team needed to expose the functionality that they care about, and ultimately ending up in a binary that's probably as large, if not larger than, the service itself. The only benefit to doing this would be if something had a particularly large and/or complicated API that would be hard to rewrite everytime you need it. Alpha Vantage is a simple web API, so I didn't think it was particularly necessary.

