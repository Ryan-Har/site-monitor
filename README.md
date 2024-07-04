# Site monitor
### A web based app built with Golang, templ and htmx to monitor and notify of any disruption of various endpoints.

You can sign up and check out the hosted version at [site-monitor.rharris.dev](https://site-monitor.rharris.dev/signup)

## Getting Started

The application uses firebase for authentication you will need to create your own firebase app and supply your own api keys to run this locally.

All configuration is done with environment variables, an example .env file is included within the repository to show the information needed.

There is also a bundled docker-compose.yml included to get up and running quickly.

## Prerequisites
[golang](https://go.dev/doc/install)

[templ](https://github.com/a-h/templ)

[firebase account for authentication](https://firebase.google.com/)

one of the dependancies [pro-bing](https://github.com/prometheus-community/pro-bing) requires system permissions, depending on the user permissions of the account running the software, you may need to provide this manually, otherwise ICMP / ping requests may fail.
```sh
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
```

## Running from source
1. Clone the repository
    ```sh
    git clone https://github.com/Ryan-Har/site-monitor
    ```
2. Populate the .env file with the correct information

3. Generate the templ files and run
    ```sh
    source .env && cd src && templ generate && go run .
    ```

## Running as a docker container

1. Copy the docker-compose.yml and the example .env file to your device. This could be done manually or by cloning the repository  as above.

2. Create a db and conf folder in your working directory, these will be used to store the sqlite database and firebase service account json respectively.
    ```sh
    mkdir {db,conf}
    ```
3. Populate the .env file with the correct information

4. Run the docker container
    ```sh
    docker compose up
    ```
