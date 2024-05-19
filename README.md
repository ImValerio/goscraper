# GoScraper

Web app that allows you to scrape any (static) website in an easy way.
![image](https://github.com/ImValerio/goscraper/assets/48352092/7053ce42-fb1b-4a60-9cd4-5399d32fa8a4)

The stack used to build this app is the following:

- React
- Golang
- Redis
- Docker

## Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/ImValerio/goscraper.git
cd goscraper
```

if you want to run the app in development mode, run the following command:

```bash
docker compose up --build
```

Instead if you want to run the app in production mode, run the following command:

```bash
docker compose -f docker-compose.prod.yml up --build -d
```
