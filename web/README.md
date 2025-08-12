# App Store RSS Reviews - Web

A React web application for viewing and filtering App Store reviews, built with React, TypeScript, and Tailwind CSS. This is the web counterpart to the mobile React Native application.

## Tech Stack

- **Framework**: React
- **UI Library**: Tailwind CSS
- **Testing**: Vitest + React Native Testing Library

## Project Structure

```
src/
├── public/              # Public resources
├── components/          # Reusable UI components
│   ├── DefaultButton.tsx  # Default button component
│   ├── Header.tsx       # App header with refresh
│   ├── ReviewItem.tsx   # Individual review card
│   └── ReviewsFilter.tsx # Rating and time filters
├── config/              # Configuration files
│   └── env.ts          # Environment configuration
├── hooks/               # Custom React hooks
│   └── useQuery.ts     # Data fetching hook
├── pages/             # Main application pages
│   └── Reviews.tsx # Main reviews page
├── services/            # API service layer
│   └── reviewsService.ts # Reviews API client
└── types/               # TypeScript type definitions
    └── reviews.ts       # Review-related types
```

## Getting Started

### Prerequisites

- Node.js (v16 or higher)
- npm or yarn
- Backend server running

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

## Testing

### Run Tests

```bash
npm test
```
