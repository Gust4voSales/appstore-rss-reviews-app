# App Store RSS Reviews Mobile App

A React Native mobile application built with Expo that displays App Store RSS Reviews fetched from the backend service. The app shows reviews from the last 96 hours (4 days) with filtering capabilities and a modern, responsive UI.

## Tech Stack

- **Framework**: React Native with Expo SDK 53
- **UI Library**: NativeWind (Tailwind CSS)
- **Testing**: Jest + React Native Testing Library

## Project Structure

```
mobile/
├── components/
│   ├── Container.tsx          # Layout container component
│   ├── Header.tsx             # App header with refresh functionality
│   ├── ReviewItem.tsx         # Individual review display component
│   └── ReviewsFilter.tsx      # Rating filter component
├── config/
│   └── env.ts                 # Environment configuration
├── hooks/
│   └── useQuery.ts            # Custom hook for data fetching
├── screens/
│   └── ReviewsScreen.tsx      # Main reviews screen
├── services/
│   └── reviewsService.ts      # API service for backend communication
├── types/
│   └── reviews.ts             # TypeScript type definitions
├── App.tsx                    # Root application component
├── app.json                   # Expo configuration
├── package.json               # Dependencies and scripts
└── tailwind.config.js         # Tailwind CSS configuration
```

## Prerequisites

- Node.js 18 or later
- npm or yarn
- Expo CLI (installed globally): `npm install -g @expo/cli`
- iOS Simulator (for iOS development) or Android Studio (for Android development)
- Running backend server (see `../server/README.md`)

## Installation

1. **Navigate to the mobile directory**:

   ```bash
   cd mobile
   ```

2. **Install dependencies**:

   ```bash
   npm install
   # or
   yarn install
   ```

3. **Configure the backend URL**:

   The app connects to the backend server. Remember to use the IP of the machine running the server; If running the app on a different machine or in an emulator, use the IP of the machine running the server

Create a `.env` file like the `.env.sample` and set the `EXPO_PUBLIC_API_URL` environment variable like the following example:

```bash
EXPO_PUBLIC_API_URL=http://192.168.2.102:8080
```

## Running the App

### Development Mode

**Start the development server**:

```bash
npm start
# or
yarn start
# or
expo start
```

This will open the Expo development tools in your browser with options to:

- Press `i` to run on iOS Simulator
- Press `a` to run on Android Emulator
- Scan QR code with Expo Go app on physical device

### Platform-Specific Commands

**iOS Simulator**:

```bash
npm run ios
# or
expo start --ios
```

**Android Emulator**:

```bash
npm run android
# or
expo start --android
```

### Physical Device Testing

1. **Install Expo Go** on your device:
   - iOS: [App Store](https://apps.apple.com/app/expo-go/id982107779)
   - Android: [Google Play](https://play.google.com/store/apps/details?id=host.exp.exponent)

2. **Remember to update API URL** to your machine's IP address like explained before

3. **Start the server** and scan the QR code with Expo Go

## Testing

### Run Tests

```bash
npm test
```

## Troubleshooting

### Network Configuration for Device Testing

When testing on physical devices, you need to ensure:

1. **Backend server is accessible**: The backend must be accessible from your device's network
2. **Correct IP address**: Use your machine's local IP address (not localhost/127.0.0.1)

**Find your machine's IP**:

```bash
# Windows
ipconfig | findstr IPv4

# macOS/Linux
ifconfig | grep inet
```
