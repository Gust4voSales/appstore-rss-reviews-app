# App Store RSS Reviews App

A full-stack application that monitors and displays iOS App Store RSS Reviews. The system consists of a Go backend service that polls App Store RSS feeds and a React Native mobile app that displays the reviews.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     App Store RSS Reviews App                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐    ┌──────────────────┐    ┌─────────────┐ │
│  │   Mobile App    │    │   Backend API    │    │ App Store   │ │
│  │  (React Native) │◄───┤      (Go)        │◄───┤ RSS Feed    │ │
│  └─────────────────┘    └──────────────────┘    └─────────────┘ │
│                                  │                              │
│                         ┌────────▼────────┐                     │
│                         │  JSON Storage   │                     │
│                         │ (reviews.json)  │                     │
│                         └─────────────────┘                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

```

## 📱 What is this project

This project is a technical challenge developed for a hiring process. The goal was to create a system composed of a backend and a frontend that consumes the App Store Connect RSS feed to fetch, store, and display recent reviews for a specific iOS app.
The backend in golang periodically fetches and persists the data to maintain state across restarts. The frontend, built with React Native, consumes an API to display reviews from the last 48 hours, sorted by newest first, including content, author, rating, and submission date.

Note: The current implementation is fixed to poll reviews for a specific app ID.

### 1. Start the Backend Server

Follow the instructions in [Server Documentation](./server/README.md)

The server will start on `http://localhost:8080` (or any other configured PORT) and begin polling for reviews immediately.

### 2. Start the Mobile App

Follow the instructions in [Mobile App Documentation](./mobile/README.md)
