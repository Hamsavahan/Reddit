# Reddit-Like Simulation Engine

This repository contains a Reddit-like simulation engine built using the Actor Model with Go and Proto.Actor. The simulation replicates key functionalities of Reddit, such as subreddit creation, post management, commenting, and a karma-based voting system. It also includes a user simulator that models real-world user behavior using Zipf distribution.

## Features
- Simulates Reddit-like functionalities including subreddits, posts, and comments.
- Implements user interactions with Zipf-distributed activity patterns.
- Highly scalable, supporting up to 20,000 concurrent users.
- Uses the Proto.Actor framework for fault tolerance and distributed processing.

## Prerequisites
- Go (version 1.18 or higher)

## Install dependencies:
- go mod tidy

## Running the Program
- To start the simulation:
- go run main.go


## Configuration Notes

- Number of Clients: By default, the number of clients is set to 150 in the main.go file. If you wish to simulate a different number of clients, update this value in the following line of code in main.go:

- sim := simulation.NewSimulation(system, 150) // Change the number of users here
- Replace 150 with the desired number of clients.
- Startup Delay: The program includes a 5-second delay for the engine to initialize before client activities begin. This ensures all actors are ready to process requests. Note that this 5-second delay is included in the reported simulation time.






