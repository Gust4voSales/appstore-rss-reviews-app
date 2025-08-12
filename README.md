# App Store RSS Reviews App

A full-stack application that monitors and displays iOS App Store RSS Reviews. The system consists of a Go backend service that polls App Store RSS feeds and frontend applications (React Native mobile app and React web app) that display the reviews.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     App Store RSS Reviews App                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐    ┌──────────────────┐    ┌─────────────┐ │
│  │   Mobile App    │    │                  │    │ App Store   │ │
│  │  (React Native) │◄───┤                  │◄───┤ RSS Feed    │ │
│  └─────────────────┘    │    Backend API   │    └─────────────┘ │
│                         │                  │                    │
│  ┌─────────────────┐    │       (Go)       │                    │
│  │    Web App      │◄───┤                  │                    │
│  │    (React)      │    └──────────────────┘                    │
│  └─────────────────┘             │                              │
│                         ┌────────▼────────┐                     │
│                         │  JSON Storage   │                     │
│                         │ (reviews.json)  │                     │
│                         └─────────────────┘                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

```

## 📱 What is this project

This project is a technical challenge developed for a hiring process. The goal was to create a system composed of a backend and multiple frontends that consume the App Store Connect RSS feed to fetch, store, and display recent reviews for a specific iOS app.
The backend in golang periodically fetches and persists the data to maintain state across restarts. The frontends, built with React Native (mobile) and React (web), consume an API to display reviews from the last 48 hours, sorted by newest first, including content, author, rating, and submission date.

Note: The current implementation is fixed to poll reviews for a specific app ID.

### 1. Start the Backend Server

Follow the instructions in [Server Documentation](./server/README.md)

The server will start on `http://localhost:8080` (or any other configured PORT) and begin polling for reviews immediately.

### 2. Start the Mobile App

Follow the instructions in [Mobile App Documentation](./mobile/README.md)

### 3. Start the Web App

Follow the instructions in [Web App Documentation](./web/README.md)

### 4. Notes

- The current implementation only polls reviews for a specific app ID.
- As a bonus, I added filters by hour (instead of being fixed to the last 48 hours) and by rating.
- I avoided using 3rd party libraries, as it was requested.

#### Improvements / Future

- Add pagination to the `/reviews` route
- Make it app-agnostic, allowing new apps to be subscribed and have their reviews fetched
- Add a country filter
