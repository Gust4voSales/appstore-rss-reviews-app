# App Store RSS Reviews - Web

A React web application for viewing and filtering App Store reviews, built with React, TypeScript, and Tailwind CSS. This is the web counterpart to the mobile React Native application.

## Features

- **Real-time Reviews**: Fetches and displays App Store reviews from the backend API
- **Rating Filter**: Filter reviews by star rating (1-5 stars)
- **Time Range Filter**: Filter reviews by time range (24h, 48h, 72h, 96h)
- **Responsive Design**: Mobile-friendly interface using Tailwind CSS
- **Auto-refresh**: Manual refresh capability with loading states
- **Back to Top**: Floating button for easy navigation
- **Error Handling**: Graceful error handling with user feedback

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Container.tsx    # Layout container
│   ├── Header.tsx       # App header with refresh
│   ├── ReviewItem.tsx   # Individual review card
│   └── ReviewsFilter.tsx # Rating and time filters
├── config/              # Configuration files
│   └── env.ts          # Environment configuration
├── hooks/               # Custom React hooks
│   └── useQuery.ts     # Data fetching hook
├── screens/             # Main application screens
│   └── ReviewsScreen.tsx # Main reviews screen
├── services/            # API service layer
│   └── reviewsService.ts # Reviews API client
└── types/               # TypeScript type definitions
    └── reviews.ts       # Review-related types
```

## Getting Started

### Prerequisites

- Node.js (v16 or higher)
- npm or yarn
- Backend server running on port 8080

### Installation

1. Install dependencies:

```bash
npm install
```

2. Set up environment variables:

   - Copy `env.example` to `.env`
   - Update `VITE_API_URL` if your backend runs on a different URL

3. Start the development server:

```bash
npm run dev
```

4. Open your browser to `http://localhost:5173`

### Environment Configuration

The application uses Vite environment variables:

- `VITE_API_URL`: Backend API URL (default: http://127.0.0.1:8080)

## Available Scripts

- `npm run dev`: Start development server
- `npm run build`: Build for production
- `npm run preview`: Preview production build
- `npm run lint`: Run ESLint

## Architecture

This web application mirrors the architecture of the mobile React Native version:

- **Component-based**: Modular, reusable components
- **Custom Hooks**: `useQuery` hook for data fetching (avoiding external dependencies)
- **Service Layer**: Centralized API communication
- **Type Safety**: Full TypeScript support
- **Responsive Design**: Tailwind CSS for styling

## API Integration

The app communicates with the backend server to fetch App Store reviews:

- **GET /reviews**: Fetch reviews with optional rating and time range filters
- **Query Parameters**:
  - `rating`: Filter by star rating (1-5)
  - `hours`: Time range in hours (24, 48, 72, 96)

## Styling

Uses Tailwind CSS for styling with:

- Responsive design patterns
- Consistent color scheme
- Hover and focus states
- Loading animations
- Mobile-first approach
